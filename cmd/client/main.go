package main

import (
	"fmt"
	pb "github.com/Pandalad1n/powwow/proto"
	tcp2 "github.com/Pandalad1n/powwow/tcp"
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
	c := tcp2.NewConnection(conn)

	time.Sleep(1 * time.Second)
	msg, err := c.ReadMessage()
	if err != nil {
		log.Fatalln(err)
	}
	pmsg := pb.Message{}
	err = proto.Unmarshal(msg, &pmsg)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(pmsg.String())

}
