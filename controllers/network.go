package controllers

import (
	"dksv/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"os/exec"
)

type NetworkController struct {
	beego.Controller
}

// 创建网络
func (this *NetworkController) Create() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type networkCreateForm struct {
		Subnet string `json:"subnet"`
		Name   string `json:"name"`
	}
	req := networkCreateForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	// 判断 req.Name 是否存在
	f, err := os.Open(models.DefaultNetworkPath + req.Name)
	// 如果已经存在，直接返回
	if err == nil {
		fmt.Println(f)
		data.Status = -1
		data.Msg = fmt.Sprintf("网络名 %s 已存在", req.Name)
		this.Data["json"] = data
		this.ServeJSON()
	}

	// ./my-docker network create --driver bridge --subnet 192.168.10.1/24 mynginxbridge
	err = exec.Command(models.MyDockerBinPath,
		"network", "create",
		"--driver", "bridge",
		"--subnet", req.Subnet,
		req.Name,
	).Run()

	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("创建网络失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
	}

	data.Data = getNetworkByName(req.Name)
	this.Data["json"] = data
	this.ServeJSON()
}

// 查看网络
func (this *NetworkController) List() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	data.Data = getAllNetwork()
	this.Data["json"] = data
	this.ServeJSON()
}

// 删除网络
func (this *NetworkController) Remove() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type networkRemoveForm struct {
		Name   string `json:"name"`
	}
	req := networkRemoveForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	// 判断 req.Name 是否存在
	f, err := os.Open(models.DefaultNetworkPath + req.Name)
	// 如果不存在，直接返回
	if err != nil {
		fmt.Println(f)
		data.Status = -1
		data.Msg = fmt.Sprintf("网络名 %s 不存在", req.Name)
		this.Data["json"] = data
		this.ServeJSON()
	}

	// ./my-docker network remove mynginxbridge
	err = exec.Command(models.MyDockerBinPath, "network", "remove", req.Name).Run()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("删除网络失败:", err)
		this.Data["json"] = data
		this.ServeJSON()
	}

	this.Data["json"] = data
	this.ServeJSON()
}

func getAllNetwork() *[]models.ContainerNetWorkConfig {
	files, _ := ioutil.ReadDir(models.DefaultNetworkPath)
	networkInfo := make([]models.ContainerNetWorkConfig, 0)

	for index := range files {
		f := files[index]
		networkInfo = append(networkInfo, *getNetworkByName(f.Name()))
	}

	return &networkInfo
}

func getNetworkByName(name string) *models.ContainerNetWorkConfig {
	content, _ := ioutil.ReadFile(models.DefaultNetworkPath + name)
	var networkInfo models.ContainerNetWorkConfig
	json.Unmarshal(content, &networkInfo)
	return &networkInfo
}
