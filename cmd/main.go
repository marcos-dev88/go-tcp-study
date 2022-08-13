package main

import (
	"log"
	"net"

	server "github.com/marcos-dev88/go-tcp-study/tcp_server"
)

const (
	networkType = "tcp"
	streamPort  = ":9000"
)

func main() {
	serverStream, err := net.Listen(networkType, streamPort)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer func(serverStream net.Listener) {
		err := serverStream.Close()
		if err != nil {

		}
	}(serverStream)

	for {
		conn, err := serverStream.Accept()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		go server.Handle(conn)
	}
}
