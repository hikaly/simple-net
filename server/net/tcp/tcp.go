package tcp

import (
	"net"
)

type Socket struct {
	SID string

	Conn *net.TCPConn
}

type TCP struct {
	// 监听地址
	Addr *net.TCPAddr

	// listener句柄
	Listener *net.TCPListener

	// 客户端链接数组
	Connections map[string]*Socket
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
		addr:     tcp_addr,
		listener: tcp_listener,
	}

	return &tcp, nil
}

func (this *TCP) AddConn(conn *net.TCPConn) {
	socket := Socket{
		SID:  "",
		Conn: conn,
	}

	this.Connections[socket.SID] = &socket

	var buffer = make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}

	}
}
