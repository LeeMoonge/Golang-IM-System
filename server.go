package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}

func (s *Server) Handle(conn net.Conn)  {
	//...当前连接的业务
	fmt.Println("建立链接成功！")
}

//启动服务的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil{
		fmt.Println("net.Listen Error ==", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept Error:", err)
			continue
		}

		//do handle
		s.Handle(conn)
	}
}