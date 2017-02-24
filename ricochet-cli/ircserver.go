package main

import (
	"fmt"
	"github.com/ricochet-im/ricochet-go/rpc"
	"net"
)

type IRCServer struct {
	Port int
}

func (s *IRCServer) ReceivedMessage(msg *ricochet.Message) {
	fmt.Printf("Should be sending to irc client...\n")
}

func (s *IRCServer) BeginListening() {
	go s.SpawnServer()
}

func (s *IRCServer) SpawnServer() {
	service := ":6667"

	listener, err := net.Listen("tcp", service)
	if err != nil {
		// :(
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.HandleClient(conn)
	}
}

func (s *IRCServer) HandleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}
