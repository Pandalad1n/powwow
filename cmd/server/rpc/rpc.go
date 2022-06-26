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

func (h Handler) ServeTCP(_ context.Context, c tcp.Conn) {
	challenge := hashcash.NewChallenge(h.Difficulty)
	buf, err := proto.Marshal(
		&pb.Message{
			Payload: &pb.Message_Challenge{
				Challenge: &pb.Challenge{
					Digest:     challenge.Digest,
					Difficulty: challenge.Difficulty,
				},
			},
		},
	)
	if err != nil {
		return
	}
	err = c.WriteMessage(buf)
	if err != nil {
		return
	}
	buf, err = c.ReadMessage()
	if err != nil {
		log.Fatalln(err)
	}
	var msg pb.Message
	err = proto.Unmarshal(buf, &msg)
	if err != nil {
		log.Fatalln(err)
	}
	sPayload, ok := msg.Payload.(*pb.Message_Solution)
	if !ok {
		log.Fatalln("Wrong message")
	}
	if !challenge.Verify(sPayload.Solution.Solution) {
		return
	}
	msg = pb.Message{
		Payload: &pb.Message_Wisdom{
			Wisdom: &pb.Wisdom{
				Text: wisdom.Wisdom(),
			},
		},
	}
	buf, err = proto.Marshal(&msg)
	if err != nil {
		return
	}
	err = c.WriteMessage(buf)
	if err != nil {
		return
	}
}
