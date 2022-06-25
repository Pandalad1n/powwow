package main

import (
	"github.com/Pandalad1n/powwow/internal/tcp"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalln(err)
	}
	c := tcp.NewConnection(conn)
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		msg := "test"
		err := c.WriteMessage([]byte(msg))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
