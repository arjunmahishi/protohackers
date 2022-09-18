package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

		go func(c net.Conn) {
			defer c.Close()

			bts, err := ioutil.ReadAll(c)
			if err != nil {
				panic(err)
			}

			c.Write(bts)
			log.Println(string(bts))
		}(conn)
	}
}

func handleKillSig() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	<-sigc
	os.Exit(0)
}
