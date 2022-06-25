package main

import (
	"context"
	"fmt"
	"github.com/Pandalad1n/powwow/internal/tcp"
	server "github.com/Pandalad1n/powwow/proto"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	ctx := context.Background()
	err := tcp.ListenAndServe(ctx, ":8888", tcp.HandleFunc(func(conn tcp.Conn) {
		for {
			b, err := conn.ReadMessage()
			if err != nil {
				log.Fatalln(err)
			}
			m := server.Request{}
			err = proto.Unmarshal(b, &m)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(m.String())
		}
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
