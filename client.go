package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	//3.1 Client新增flag属性
	flag int //当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	//1.1 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	//1.2 连接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial Error:", err)
		return nil
	}

	client.conn = conn

	//1.3 返回对象
	return client
}

//3.2 新增menu方法，获取用户输入的模式
func (c *Client) menu() bool {
	var input int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&input)

	if input >= 0 && input <= 3 {
		c.flag = input
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围的数字<<<<<")
		return false
	}
}

//3.3 新增Run()主业务循环
func (c *Client) Run()  {
	for c.flag != 0 {
		for c.menu() != true {
		}

		//3.4 根据不同的模式处理不同的业务
		switch c.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式选择 ... ")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式选择 ... ")
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名选择 ... ")
			break
		}
	}
}

//2.1 定义全局变量参数，使用命令行传参
var serverIp string
var serverPort int

//2.2 在main函数前使用init加载命令行输入的参数
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8080, "设置服务器端口(默认是8080)")
}

func main() {
	//client := NewClient("127.0.0.1", 8080)

	//2.3 命令行解析
	flag.Parse()
	//2.4 使用命令行输入的ip和port请求服务器
	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println(">>>>>>连接服务器失败 ... ")
		return
	}

	fmt.Println(">>>>>>连接服务器成功 ... ")

	//启动客户端的业务
	//select {}
	client.Run()
}
