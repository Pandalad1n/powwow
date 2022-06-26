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
