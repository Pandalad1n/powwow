package rpc

import (
	"context"
	"github.com/Pandalad1n/powwow/hashcash"
	pb "github.com/Pandalad1n/powwow/proto"
	"github.com/Pandalad1n/powwow/tcp"
	"github.com/Pandalad1n/powwow/wisdom"
	"google.golang.org/protobuf/proto"
	"log"
)

type Handler struct {
	Difficulty uint32
}

func (h Handler) ServeTCP(ctx context.Context, c tcp.Conn) {
	challenge := hashcash.NewChallenge(h.Difficulty)
	cMessage := pb.Message{Payload: &pb.Message_Challenge{Challenge: &pb.Challenge{Digest: challenge.Digest, Difficulty: challenge.Difficulty}}}
	pMessage, err := proto.Marshal(&cMessage)
	if err != nil {
		return
	}
	err = c.WriteMessage(pMessage)
	if err != nil {
		return
	}
	resp, err := c.ReadMessage()
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
		if challenge.Verify(solution.Solution) {
			wMessage := pb.Message{Payload: &pb.Message_Wisdom{Wisdom: &pb.Wisdom{Text: wisdom.Wisdom()}}}
			pMessage, err = proto.Marshal(&wMessage)
			if err != nil {
				return
			}
			err = c.WriteMessage(pMessage)
			if err != nil {
				return
			}
		}
	}
}
