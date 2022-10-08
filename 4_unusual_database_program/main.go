package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync"
)

type Packet struct {
	addr    net.Addr
	message []byte
}

var (
	unusualDB    = map[string]string{"version": "unusual DB v0.1"}
	unusualMutex = sync.RWMutex{}
)

const (
	hostAddr = ":9000"
	maxBytes = 1000
)

func main() {
	conn, err := net.ListenPacket("udp", hostAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Listening for UPD connections on %s...", hostAddr)
	for {
		buff := make([]byte, 1024)
		nBytes, addr, err := conn.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}

		if nBytes > maxBytes {
			continue
		}

		go handleRequest(conn, Packet{
			addr:    addr,
			message: bytes.TrimSpace(buff[:nBytes]),
		})
	}

}

func handleRequest(conn net.PacketConn, packet Packet) {
	switch {
	case bytes.Contains(packet.message, []byte("=")): // insert
		msgSplit := bytes.SplitN(packet.message, []byte("="), 2)
		key, value := string(msgSplit[0]), string(msgSplit[1])
		if key == "version" {
			break
		}

		insert(key, value)
		log.Printf("req (%d): %s", len(packet.message), string(packet.message))
	default: // retrieve
		resp := retrieve(string(packet.message))
		log.Printf(
			"req (%d): %s, resp: %s",
			len(packet.message), string(packet.message), resp,
		)
		conn.WriteTo([]byte(resp), packet.addr)
	}
}

func insert(key, value string) {
	unusualMutex.Lock()
	defer unusualMutex.Unlock()

	unusualDB[key] = value
}

func retrieve(key string) string {
	unusualMutex.RLock()
	defer unusualMutex.RUnlock()

	value := unusualDB[key]
	respFormat := "%s=%s"

	return fmt.Sprintf(respFormat, key, value)
}
