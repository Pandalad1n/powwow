package main

import (
	"context"
	"fmt"
	"github.com/Pandalad1n/powwow/internal/tcp"
	"log"
)

func main() {
	ctx := context.Background()
	err := tcp.ListenAndServe(ctx, ":8080", tcp.HandleFunc(func(conn tcp.Conn) {
		for {
			b, err := conn.ReadMessage()
			if err != nil {
				log.Fatalln(err)
			}
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(b))
		}
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
