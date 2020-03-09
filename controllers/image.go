package controllers

import (
	"dksv/models"
	"github.com/astaxie/beego"
)

type ImageController struct {
	beego.Controller
}

// 从服务器拉取镜像
func (this *ImageController) Pull() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 列出本机的镜像
func (this *ImageController) List() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 删除本机的镜像
func (this *ImageController) Remove() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 将本机的容器打包成镜像
func (this *ImageController) Commit() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 将本机镜像推送到服务器
func (this *ImageController) Push() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	this.Data["json"] = data
	this.ServeJSON()
}
