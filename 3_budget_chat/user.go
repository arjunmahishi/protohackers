package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type User struct {
	name string
	conn net.Conn
	room *ChatRoom
}

func (u *User) startChat() {
	defer u.conn.Close()

	buf := bufio.NewReader(u.conn)
	for {
		msg, err := buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				u.room.removeUser(u)
				return
			}

			log.Println(err)
			return
		}

		if err := u.send(msg); err != nil {
			log.Println(err)
			u.conn.Write([]byte(err.Error()))
			return
		}
	}
}

func (u *User) send(msg []byte) error {
	msg = bytes.TrimSpace(msg)
	if len(msg) == 0 {
		return nil
	}

	finalMsg := []byte(fmt.Sprintf(
		"[%s] %s\n", u.name, string(msg),
	))
	return u.room.broadcast(u, finalMsg)
}
