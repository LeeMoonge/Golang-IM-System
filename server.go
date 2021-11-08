package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int

	//2.在线用户的列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//2.消息广播的channel
	Message chan string
}

//1.创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,

		//2.新建连接时创建对应的列表和管道
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

//2. 广播当前用户上线消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

//2. 监听Message广播消息的channel的goroutine，一旦有消息就发送给全部在线的User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		//将msg发送给全部的在线User
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap{
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Handle(conn net.Conn)  {
	//2. 用户上线新建用户对象
	user := NewUser(conn, s)

	//1 ...当前连接的业务
	fmt.Println(user.Name + "建立链接成功！")


	//4. 优化代码，将逻辑封装在user.go中
	/*
	//2. 用户上线将用户加入到onlineMap中
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//2. 广播当前用户上线消息
	s.BroadCast(user, "上线了")
	*/
	//4.
	user.Online()

	//3. 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//4. 优化代码，将逻辑封装在user.go中
				/*s.BroadCast(user, "已下线！")*/
				//4.
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Error :", err)
				return
			}

			//3. 提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			//4. 优化代码，将逻辑封装在user.go中
			/*
			//3. 将得到的消息进行广播
			s.BroadCast(user, msg)
			*/
			//4. 用户针对msg进行消息处理
			user.DoMessage(msg)
		}
	}()

	//2. 当前handle阻塞
	select {}
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

	//2. 启动监听Message的goroutine
	go s.ListenMessage()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept Error:", err)
			continue
		}

		//do handle
		go s.Handle(conn)
	}
}