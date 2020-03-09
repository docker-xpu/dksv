package controllers

import (
	"bufio"
	"dksv/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ContainerController struct {
	beego.Controller
}

// 运行容器
func (this *ContainerController) Run() {
	type containerRunForm struct {
		ImageName     string   `json:"image_name"`
		ContainerName string   `json:"container_name"`
		Volume        string   `json:"volume"`
		PortMapping   []string `json:"port_mapping"`
		MemoryLimit   string   `json:"memory_limit"`
		CpuShare      string   `json:"cpu_share"`
		CpuSet        string   `json:"cpu_set"`
		Command       string   `json:"command"`
		Net           string   `json:"net"`
	}
	c := containerRunForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &c)

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	// 判断是否有相同的容器名存在
	containers := getAllContainerInfo()
	flag := false
	for index := range containers {
		if containers[index].Name == c.ContainerName {
			data.Status = -1
			data.Msg = "容器" + c.ContainerName + "已存在"
			flag = true
			break
		}
	}
	// todo 判断镜像是否存在
	//

	// 如果没有相同的容器
	// 执行 ./my-docker run -d -name brid2 -net mynginxbridge -p 8888:80 mynginx nginx
	if flag == false {
		cmd := exec.Command(
			models.MyDockerBinPath,
			"run", "-d",
			"-name", c.ContainerName,
			"-net", c.Net,
			"-v", c.Volume,
			"-p", c.PortMapping[0], // todo 多端口映射
			"-cpushare", c.CpuShare,
			"-cpuset", c.CpuSet,
			"-m", c.MemoryLimit,
			c.ImageName,
			c.Command,
		)
		err := cmd.Run()
		if err != nil {
			data.Status = -1
			data.Msg = fmt.Sprintf("创建容器失败:%v", err)
		} else {
			data.Data, _ = getContainerInfo(c.ContainerName)
		}
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 停止容器
func (this *ContainerController) Stop() {
	type containerStopForm struct {
		ContainerName string `json:"container_name"`
	}
	c := containerStopForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &c)

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	// 判断容器是否存在以及状态
	containerInfo, _ := getContainerInfo(c.ContainerName)
	if containerInfo.Status != models.RUNNING {
		data.Status = -1
		data.Msg = "容器已经停止或不存在"
	}
	if containerInfo.Status == models.RUNNING {
		cmd := exec.Command(models.MyDockerBinPath, "stop", c.ContainerName)
		err := cmd.Run()
		if err != nil {
			data.Status = -1
			data.Msg = fmt.Sprintf("容器停止失败:%v", err)
		}
		data.Data, _ = getContainerInfo(c.ContainerName)
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 删除容器
func (this *ContainerController) Remove() {
	fmt.Println("Remove()")
	type containerRemoveForm struct {
		ContainerName string `json:"container_name"`
	}
	c := containerRemoveForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &c)

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	// 判断容器是否存在以及状态
	containerInfo, _ := getContainerInfo(c.ContainerName)
	if containerInfo.Status != models.STOP {
		data.Status = -1
		data.Msg = "容器不存在或正在运行中"
	}
	if containerInfo.Status == models.STOP {
		cmd := exec.Command(models.MyDockerBinPath, "rm", c.ContainerName)
		err := cmd.Run()
		if err != nil {
			data.Status = -1
			data.Msg = fmt.Sprintf("容器删除失败:%v", err)
		}
		data.Data, _ = getContainerInfo(c.ContainerName)
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 列出容器
func (this *ContainerController) List() {
	containers := getAllContainerInfo()

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   containers,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 查看容器的日志
func (this *ContainerController) Logs() {
	containerName := this.GetString("container_name")

	containerLogs := getContainerLogs(containerName)

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   containerLogs,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 从本机获取所有容器的信息
func getAllContainerInfo() []*models.ContainerInfo {
	dirURL := fmt.Sprintf(models.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, _ := ioutil.ReadDir(dirURL)

	var containers []*models.ContainerInfo
	for _, file := range files {
		if file.Name() == "network" {
			continue
		}
		tmpContainer, _ := getContainerInfo(file.Name())
		containers = append(containers, tmpContainer)
	}

	return containers
}

// 根据 containerName 获取容器信息，并返回
func getContainerInfo(containerName string) (*models.ContainerInfo, error) {
	containerInfoPath := fmt.Sprintf(models.DefaultInfoLocation, containerName)
	containerConfigFile := containerInfoPath + models.ConfigName

	// 读取容器的 config.json 文件
	content, _ := ioutil.ReadFile(containerConfigFile)
	var containerInfo models.ContainerInfo
	json.Unmarshal(content, &containerInfo)

	resourceConfig := models.ContainerResourceConfig{}
	// 读取容器的 memory limit 相关的文件
	memCgPath := models.CgroupPath + "/memory/" + containerInfo.Id + "/"
	memLimitFile := memCgPath + "memory.limit_in_bytes"
	memLimit, _ := ioutil.ReadFile(memLimitFile)
	memUsedFile := memCgPath + "memory.usage_in_bytes"
	memUsed, _ := ioutil.ReadFile(memUsedFile)

	resourceConfig.MemoryLimit.Limit = strings.TrimSuffix(string(memLimit), "\n")
	resourceConfig.MemoryLimit.Used = strings.TrimSuffix(string(memUsed), "\n")

	// 读取 cpu.shares 限制
	cpuCgPath := models.CgroupPath + "/cpu/" + containerInfo.Id + "/"
	cpuShareFile := cpuCgPath + "cpu.shares"
	cpuShares, _ := ioutil.ReadFile(cpuShareFile)

	resourceConfig.CpuShare = strings.TrimSuffix(string(cpuShares), "\n")

	// 读取 cpu set
	cpuSetCgPath := models.CgroupPath + "/cpuset/" + containerInfo.Id + "/"
	cpuSetFile := cpuSetCgPath + "cpuset.cpus"
	cpuSet, _ := ioutil.ReadFile(cpuSetFile)

	resourceConfig.CpuSet = strings.TrimSuffix(string(cpuSet), "\n")
	containerInfo.Limits = resourceConfig

	// 使用 gopsutil 获取容器进程信息
	containerPid, err := strconv.Atoi(containerInfo.Pid)
	if err == nil { // 如果转换成功
		proc := process.Process{Pid: int32(containerPid)}
		extra := make(map[string]interface{})
		extra["Background"], _ = proc.Background()   // 程序是否后台运行
		extra["CPUAffinity"], _ = proc.CPUAffinity() // CPU 亲和力
		extra["CPUPercent"], _ = proc.CPUPercent()   // 进程使用多少cpu时间
		extra["Children"], _ = proc.Children()       // proc的子进程
		extra["Cmdline"], _ = proc.CmdlineSlice()    //进程的命令行参数
		extra["Connections"], _ = proc.Connections() // 网络连接数
		extra["CreateTime"], _ = proc.CreateTime()   // 进程创建时间
		extra["Cwd"], _ = proc.Cwd()                 // 进程工作目录
		extra["Exe"], _ = proc.Exe()                 // 进程的可执行路径。
		extra["IOCounters"], _ = proc.IOCounters()   // IO相关
		extra["IOnice"], _ = proc.IOnice()           // io nice 值（优先级）
		extra["IsRunning"], _ = proc.IsRunning()     // 是否在运行
		extra["MemoryInfo"], _ = proc.MemoryInfo()
		extra["MemoryInfoEx"], _ = proc.MemoryInfoEx()
		extra["MemoryPercent"], _ = proc.MemoryPercent() //此过程使用的总RAM的百分之多少
		extra["Name"], _ = proc.Name()                   // process name
		extra["Nice"], _ = proc.Nice()                   // 进程 nice 值
		extra["NumFDs"], _ = proc.NumFDs()               // 打开的 fd 数
		extra["NumThreads"], _ = proc.NumThreads()       // 线程数

		// extra["Rlimit"], _ = proc.Rlimit()
		containerInfo.Extra = extra
	}

	return &containerInfo, nil
}

// 获取容器日志信息
func getContainerLogs(containerName string) []string {
	containerInfoPath := fmt.Sprintf(models.DefaultInfoLocation, containerName)
	containerLogFile := containerInfoPath + models.ContainerLogFile

	containerLogs := make([]string, 1)

	f, _ := os.Open(containerLogFile)
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		buf, e := r.ReadBytes('\n')
		if e != nil && len(buf) == 0 {
			break
		}
		containerLogs = append(containerLogs, string(buf))
	}

	for index := range containerLogs {
		fmt.Println(containerLogs[index])
	}

	return containerLogs
}
