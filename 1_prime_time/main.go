package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

var (
	connCount = 0
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

	buf := bufio.NewReader(c)
	for {
		req, err := buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Println(err)
			return
		}

		start := time.Now()
		resp, err := handleRequest(req)
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			return
		}

		c.Write(append(resp, []byte("\n")...))
		log.Printf("[%d] Req: %s", connID, string(req))
		log.Printf("[%d] Resp: %s %s", connID, string(resp), time.Since(start))
	}
}

func handleRequest(inp []byte) ([]byte, error) {
	var req request
	if err := json.Unmarshal(inp, &req); err != nil {
		return nil, err
	}

	if req.Method == nil || req.Number == nil {
		return nil, errors.New("bad request")
	}

	if *req.Method != "isPrime" {
		return nil, errors.New("method not supported")
	}

	return json.Marshal(response{
		Method: *req.Method,
		Prime:  isPrime(*req.Number),
	})
}

func isPrime(value float64) bool {
	if float64(int64(value)) != value {
		return false
	}

	for i := 2.0; i <= math.Sqrt(value); i++ {
		if math.Mod(value, i) == 0 {
			return false
		}
	}

	return value > 1
}

func handleKillSig() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	<-sigc
	os.Exit(0)
}
