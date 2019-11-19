package tcp

import (
	"fmt"
	"net"
)

type Socket struct {
	Id string

	Conn *net.TCPConn
}

type TCP struct {
	// 监听地址
	Addr *net.TCPAddr

	// listener句柄
	Listener *net.TCPListener

	// 客户端链接数组
	Connections map[string]*Socket

	/*
		// 全局消息队列
		MessageQueue [][]byte
	*/
}

func Create(listen_addr string) (*TCP, error) {
	tcp_addr, err := net.ResolveTCPAddr("tcp", listen_addr)
	if err != nil {
		return nil, err
	}

	tcp_listener, err := net.ListenTCP("tcp", tcp_addr)
	if err != nil {
		return nil, err
	}

	tcp := TCP{
		Addr:     tcp_addr,
		Listener: tcp_listener,
	}

	return &tcp, nil
}

func (this *TCP) AddConn(conn *net.TCPConn) {
	// TODO: session_id(uuid)
	socket := Socket{
		Id:   "",
		Conn: conn,
	}

	this.Connections[socket.SID] = &socket

	var buffer = make([]byte, 1024)
	var err error
	var n int
	for {
		n, err = conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}

		DataFilter(socket.Id, buffer[:n], n)
		buffer = buffer[:0]
	}
}
