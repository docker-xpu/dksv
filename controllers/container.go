package controllers

import (
	"bufio"
	"dksv/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"strings"
)

// /api/container/run/
type Container struct {
	beego.Controller
}

func (this *Container) Run() {
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
		HookURL       string   `json:"hook_url"`
	}
	c := containerRunForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &c)

	// todo 执行 my-docker run

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   "",
	}
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *Container) Stop() {
	type containerStopForm struct {
		ContainerName string `json:"container_name"`
	}
	c := containerStopForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &c)

	// todo 执行 my-docker stop

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   c,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *Container) Remove() {
	type containerRemoveForm struct {
		ContainerName string `json:"container_name"`
	}
	c := containerRemoveForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &c)

	// todo 删除容器

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   c,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *Container) List() {
	containers := getAllContainerInfo()

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   containers,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *Container) Logs() {
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

// 判断容器是否存在，如果存在，返回容器状态，否则返回 ""
func containerStatus(containerName string) string {
	return ""
}
