package main

import (
	"fmt"
	"github.com/Pandalad1n/powwow/hashcash"
	pb "github.com/Pandalad1n/powwow/proto"
	"github.com/Pandalad1n/powwow/tcp"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "powwow_server:8888")
	if err != nil {
		log.Fatalln(err)
	}
	c := tcp.NewConnection(conn)
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}
		pmsg := pb.Message{}
		err = proto.Unmarshal(msg, &pmsg)
		if err != nil {
			log.Fatalln(err)
		}
		switch pmsg.Payload.(type) {
		case *pb.Message_Challenge:
			pChallenge := pmsg.GetChallenge()
			challenge := hashcash.Challenge{Digest: pChallenge.Digest, Difficulty: pChallenge.Difficulty}
			solution := challenge.Solve()
			sMessage := pb.Message{Payload: &pb.Message_Solution{Solution: &pb.Solution{Solution: solution}}}
			pMessage, err := proto.Marshal(&sMessage)
			if err != nil {
				return
			}
			err = c.WriteMessage(pMessage)
			if err != nil {
				return
			}
		case *pb.Message_Wisdom:
			fmt.Println(pmsg.GetWisdom().String())
			return
		}

	}

}
