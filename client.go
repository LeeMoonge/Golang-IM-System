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
}

func NewClient(serverIp string, serverPort int) *Client {
	//1.1 创建客户端对象
	client := &Client{
		ServerIp:serverIp,
		ServerPort: serverPort,
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

//2.1 定义全局变量参数，使用命令行传参
var serverIp string
var serverPort int

//2.2 在main函数前使用init加载命令行输入的参数
func init()  {
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
	select {}
}