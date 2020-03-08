package controllers

import (
	"dksv/models"
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Post() {
	beego.Info("固定路由的get类型的方法 ")

	//cmd := exec.Command("ls", "-al")
	//cmd.Stdout =

	cliConf := new(ClientConfig)
	cliConf.createClient("10.211.55.15", 22, "root", "1Liujingliang")

	res := cliConf.RunShell("cat /var/run/my-docker/brid/config.json")
	fmt.Printf("type of res:%T, res:%v\n", res, res)

	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		//Data:   "存活",
	}
	c.Data["json"] = data
	c.ServeJSON()
}

type ClientConfig struct {
	Host       string
	Port       int64
	Username   string
	Password   string
	Client     *ssh.Client // ssh cli
	LastResult string      // 最近一次运行结果
}

func (cliConf *ClientConfig) createClient(host string, port int64, username, password string) {
	var (
		client *ssh.Client
		err    error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password
	cliConf.Port = port

	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	// 获取 cli
	if client, err = ssh.Dial("tcp", addr, &config); err != nil {
		log.Fatalln("error:", err)
	}
	cliConf.Client = client
}

func (cliConf *ClientConfig) RunShell(shell string) string {
	var (
		session *ssh.Session
		err     error
	)

	// 获取 session，这个session 是用来远程执行操作的
	if session, err = cliConf.Client.NewSession(); err != nil {
		log.Fatalln("error:", err)
	}

	// 执行 shell
	if output, err := session.CombinedOutput(shell); err != nil {
		log.Fatalln("error:", err)
	} else {
		cliConf.LastResult = string(output)
	}
	return cliConf.LastResult
}
