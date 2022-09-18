package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	connCount = 0
)

const (
	I = 73
	Q = 81
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

		connCount++
		go handleConn(connCount, conn)
	}
}

func handleConn(connID int, c net.Conn) {
	defer c.Close()

	db := &DB{}

	for {
		req := make([]byte, 9)
		_, err := io.ReadFull(c, req)
		if err != nil {
			log.Println(err)
			return
		}

		resp, err := handleRequest(connID, db, req)
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			return
		}

		if len(resp) > 0 {
			c.Write(resp)
		}

	}
}

func handleRequest(connID int, db *DB, cmd []byte) ([]byte, error) {
	start := time.Now()
	isInsert, arg1, arg2, err := decodeCmd(cmd)
	if err != nil {
		return nil, err
	}

	log.Printf("[%d] Req: %t %d %d", connID, isInsert, arg1, arg2)

	if isInsert {
		db.insert(arg1, arg2)
		log.Printf("[%d] Processed %s", connID, time.Since(start))
		return nil, nil
	}

	mean, err := db.query(arg1, arg2)
	if err != nil {
		return nil, err
	}

	resp := make([]byte, 4)
	binary.BigEndian.PutUint32(resp, uint32(mean))

	log.Printf("[%d] Resp: %d %s", connID, mean, time.Since(start))
	return resp, nil
}

func decodeCmd(msg []byte) (bool, int32, int32, error) {
	switch msg[0] {
	case I:
		return true, int32(binary.BigEndian.Uint32(msg[1:5])),
			int32(binary.BigEndian.Uint32(msg[5:])), nil
	case Q:
		return false, int32(binary.BigEndian.Uint32(msg[1:5])),
			int32(binary.BigEndian.Uint32(msg[5:])), nil
	default:
		return false, 0, 0, errors.New("unknown command type")
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
