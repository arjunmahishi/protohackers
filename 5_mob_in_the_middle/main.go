package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

const (
	chatAddr = "chat.protohackers.com:16963"
)

var (
	tokenPat  = regexp.MustCompile(`(?m)^7\S{25,34}$`)
	tonyToken = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

func main() {
	go handleKillSig()
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

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

func handleNewConn(client net.Conn) {
	chatServer, err := net.Dial("tcp", chatAddr)
	if err != nil {
		log.Println("Error:", err)
	}

	session := &Session{
		chatServer: chatServer,
		client:     client,
		kill:       make(chan bool),
	}

	session.start()
}

func replaceToken(msg string) string {
	msgSplit := strings.Split(msg, " ")
	newWords := []string{}

	for _, word := range msgSplit {
		newWords = append(newWords, tokenPat.ReplaceAllString(word, tonyToken))
	}

	return strings.Join(newWords, " ")
}

func handleKillSig() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	<-sigc
	os.Exit(0)
}
