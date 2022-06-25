package main

import (
	"context"
	"github.com/Pandalad1n/powwow/hashcash"
	pb "github.com/Pandalad1n/powwow/proto"
	"github.com/Pandalad1n/powwow/tcp"
	"github.com/Pandalad1n/powwow/wisdom"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	ctx := context.Background()
	err := tcp.ListenAndServe(ctx, "0.0.0.0:8888", tcp.HandleFunc(func(conn tcp.Conn) {
		c := hashcash.NewChallenge(10)
		cMessage := pb.Message{Payload: &pb.Message_Challenge{Challenge: &pb.Challenge{Digest: c.Digest, Difficulty: c.Difficulty}}}
		pMessage, err := proto.Marshal(&cMessage)
		if err != nil {
			return
		}
		err = conn.WriteMessage(pMessage)
		if err != nil {
			return
		}
		resp, err := conn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}
		pmsg := pb.Message{}
		err = proto.Unmarshal(resp, &pmsg)
		if err != nil {
			log.Fatalln(err)
		}
		switch pmsg.Payload.(type) {
		case *pb.Message_Solution:
			solution := pmsg.GetSolution()
			if c.Verify(solution.Solution) {
				wMessage := pb.Message{Payload: &pb.Message_Wisdom{Wisdom: &pb.Wisdom{Text: wisdom.GetWisdom()}}}
				pMessage, err = proto.Marshal(&wMessage)
				if err != nil {
					return
				}
				err = conn.WriteMessage(pMessage)
				if err != nil {
					return
				}
			}
		}
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
