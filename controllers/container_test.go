package controllers

import (
	"bufio"
	"dksv/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test001(t *testing.T) {
	containerName := "brid"
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

	resourceConfig.MemoryLimit.Limit = string(memLimit)
	resourceConfig.MemoryLimit.Used = string(memUsed)

	// 读取 cpu.shares 限制
	cpuCgPath := models.CgroupPath + "/cpu/" + containerInfo.Id + "/"
	cpuShareFile := cpuCgPath + "cpu.shares"
	cpuShares, _ := ioutil.ReadFile(cpuShareFile)

	resourceConfig.CpuShare = string(cpuShares)

	// 读取 cpu set
	cpuSetCgPath := models.CgroupPath + "/cpuset/" + containerInfo.Id + "/"
	cpuSetFile := cpuSetCgPath + "cpuset.cpus"
	cpuSet, _ := ioutil.ReadFile(cpuSetFile)

	resourceConfig.CpuSet = string(cpuSet)

	containerInfo.Limits = resourceConfig
}

func Test002(t *testing.T) {
	containerName := "brid"
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
}
