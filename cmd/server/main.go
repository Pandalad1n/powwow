package main

import (
	"context"
	"github.com/Pandalad1n/powwow/hashcash"
	pb "github.com/Pandalad1n/powwow/proto"
	"github.com/Pandalad1n/powwow/tcp"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	ctx := context.Background()
	err := tcp.ListenAndServe(ctx, ":8888", tcp.HandleFunc(func(conn tcp.Conn) {
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
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
