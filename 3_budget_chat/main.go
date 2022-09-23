package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

var (
	alphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)

func main() {
	go handleKillSig()

	fmt.Println("listening on port 9000...")
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}

		go handleNewConn(conn)
	}
}

func handleNewConn(c net.Conn) {
	buf := bufio.NewReader(c)

	c.Write([]byte("Welcome to budgetchat! What shall we call you?\n"))
	username, err := buf.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		c.Close()
		return
	}

	username = bytes.TrimSpace(username)
	if len(username) == 0 {
		log.Println("invalid username")
		c.Close()
		return
	}

	if !alphaNumeric.Match(username) {
		log.Println("invalid username")
		c.Close()
		return
	}

	currentUsers := getChatRoom().getUserNames()
	user := &User{
		name: strings.TrimSpace(string(username)),
		conn: c,
		room: getChatRoom(),
	}

	if err := getChatRoom().addUser(user); err != nil {
		c.Write([]byte(err.Error()))
		c.Close()
		return
	}

	c.Write([]byte(fmt.Sprintf(
		"* people in this room: %s\n", strings.Join(currentUsers, ", "),
	)))
	fmt.Printf("* people in this room: %s\n", strings.Join(currentUsers, ", "))

	go user.startChat()
}

func handleKillSig() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	<-sigc
	os.Exit(0)
}
