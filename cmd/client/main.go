package main

import (
	"github.com/Pandalad1n/powwow/internal/tcp"
	server "github.com/Pandalad1n/powwow/proto"
	"google.golang.org/protobuf/proto"
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
		msg := server.Request{Payload: &server.Request_Challenge_{Challenge: &server.Request_Challenge{}}}
		buf, err := proto.Marshal(&msg)
		if err != nil {
			log.Fatalln(err)
		}
		err = c.WriteMessage(buf)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
