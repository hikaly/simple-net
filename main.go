package main

import (
	"simple-net/server"
)

func main() {
	conf, err := server.ReadConfig()
	if err != nil {
		return
	}

	err = server.Init(conf)
	if err != nil {
		return
	}

	server.Start()
}
