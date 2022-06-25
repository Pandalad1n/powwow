package main

import (
	"github.com/Pandalad1n/powwow/internal/tcp"
	"log"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	c := tcp.NewConnection(conn)
	for {
		time.Sleep(1 * time.Second)
		msg := "test"
		err := c.WriteMessage([]byte(msg))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
