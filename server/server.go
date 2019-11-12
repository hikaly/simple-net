package server

import (
	"simple-net/server/net/tcp"
)

type Global struct {
	// 模块handle
	Modules map[string]*Module

	// 消息处理handle
	Message map[string]MesageHandle

	// 最大链接数(tcp socket)
	MaxSocket int32

	// listen types
	ListenTypes []string

	// tcp
	TCPListener *tcp.TCP

	// http

	// udp
}

var global *Global

func Start() {
	for {
		tcp_conn, err := global.TCPListener.Listener.AcceptTCP()
		if err != nil {
			// todo:
		}

		if int(global.MaxSocket) == len(global.TCPListener.Connections) {
			// TODO:
			tcp_conn.Close()
			continue
		}
	}

}

func Init(conf *Config) error {
	global := &Global{}

	listener, err := tcp.Create("0.0.0.0:8864")
	if err != nil {
		return err
	}

	global.TCPListener = listener
}
