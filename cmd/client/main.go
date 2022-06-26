package main

import (
	"flag"
	"fmt"
	"github.com/Pandalad1n/powwow/hashcash"
	pb "github.com/Pandalad1n/powwow/proto"
	"github.com/Pandalad1n/powwow/tcp"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	host := flag.String("host", "", "server host")
	port := flag.Int("port", 8888, "server port")
	flag.Parse()

	c, err := tcp.Dial(fmt.Sprintf("%s:%v", *host, *port))
	if err != nil {
		log.Fatalln(err)
	}
	// We could bring another abstraction.
	// But we keep it simple.
	buf, err := c.ReadMessage()
	if err != nil {
		log.Fatalln(err)
	}
	var serverMsg pb.Message
	err = proto.Unmarshal(buf, &serverMsg)
	if err != nil {
		log.Fatalln(err)
	}
	cPayload, ok := serverMsg.Payload.(*pb.Message_Challenge)
	if !ok {
		log.Fatalln("Wrong message")
	}
	challenge := hashcash.Challenge{
		Digest:     cPayload.Challenge.Digest,
		Difficulty: cPayload.Challenge.Difficulty,
	}
	solution, err := proto.Marshal(
		&pb.Message{
			Payload: &pb.Message_Solution{
				Solution: &pb.Solution{
					Solution: challenge.Solve(),
				},
			},
		},
	)
	if err != nil {
		return
	}
	err = c.WriteMessage(solution)
	if err != nil {
		return
	}

	buf, err = c.ReadMessage()
	if err != nil {
		log.Fatalln(err)
	}
	err = proto.Unmarshal(buf, &serverMsg)
	if err != nil {
		log.Fatalln(err)
	}
	wPayload, ok := serverMsg.Payload.(*pb.Message_Wisdom)
	if !ok {
		log.Fatalln("Wrong message")
	}
	fmt.Println(wPayload.Wisdom.Text)
}
