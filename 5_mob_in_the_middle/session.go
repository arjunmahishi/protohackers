package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type Session struct {
	sync.RWMutex
	chatServer net.Conn
	client     net.Conn
	kill       chan bool
}

func (s *Session) start() {
	go s.forwardToServer()
	go s.forwardToClient()
}

func (s *Session) forwardToServer() {
	for {
		msgFromClient, err := bufio.NewReader(s.client).ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.chatServer.Close()
				s.kill <- true
				return
			}
			log.Println("error:", err)
			return
		}

		msgFromClient = replaceToken(msgFromClient)

		if _, err := s.chatServer.Write([]byte(msgFromClient)); err != nil {
			log.Println("error:", err)
		}
	}
}

func (s *Session) forwardToClient() {
	for {
		select {
		case <-s.kill:
			return
		default:
			msgFromServer, err := bufio.NewReader(s.chatServer).ReadString('\n')
			if err != nil {
				log.Println("error:", err)
				return
			}

			log.Println(msgFromServer)
			msgFromServer = replaceToken(msgFromServer)
			if _, err := s.client.Write([]byte(msgFromServer)); err != nil {
				log.Println("error:", err)
			}
		}
	}
}
