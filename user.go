package main

import (
	"net"
	"strings"
)

//2.创建用户类

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
		server: server,
	}

	//启动监听当前User channel消息的goroutine
	go user.ListenMessage()

	return user
}

//监听当前User channel的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}

//4. 对User类进行逻辑梳理和封装
func (u *User) Online() {
	//用户上线，将用户加入到onlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播当前用户上线消息
	u.server.BroadCast(u, "上线了！")
}

//4. 
func (u *User)Offline()  {
	//用户下线，将用户从onlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播当前用户下线消息
	u.server.BroadCast(u, "下线了！")
}

//4. 
func (u *User)DoMessage(msg string)  {
	//5. 用户处理消息的业务
	if msg == "who" {
		//6. 查询当前在线用户都有哪些

		u.server.mapLock.Lock() //!!!!!!!有上锁一定要解锁
		for _, user := range u.server.OnlineMap {
			onlionMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMessage(onlionMsg)
		}
		u.server.mapLock.Unlock() //!!!!!!!!有上锁一定要解锁
	//7. 修改用户名
	} else if len(msg) >= 7 && msg[:7] == "rename|" {
		//7.1 消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		//7.2 判断名称是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMessage("当前名称已被占用！\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMessage("已成功更新用户名：" + newName + "\n")
		}
	//9. 用户私信功能
	} else if len(msg) >= 4 && msg[:3] == "to|" {
		//9.1 消息格式： to|张三|消息内容

		//9.1.1 获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMessage("消息格式不正确，请使用\"to|张三|你好啊\"格式。\n")
			return
		}

		//9.1.2 根据用户名 得到对方的User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMessage("该用户名不存在或不在线，请重新发送！\n")
			return
		}

		//9.1.3 获取消息内容，通过对方的User对象将消息发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMessage("无效发送内容，请重发！\n")
			return
		}
		remoteUser.SendMessage(u.Name + "对你说:" + content)

	} else {
		u.server.BroadCast(u, msg)
	}
}

func (u *User)SendMessage(msg string)  {
	u.conn.Write([]byte(msg))
}