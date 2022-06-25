package main

import (
	"context"
	"github.com/Pandalad1n/powwow/cmd/server/rpc"
	"github.com/Pandalad1n/powwow/tcp"
	"log"
)

func main() {
	ctx := context.Background()
	err := tcp.ListenAndServe(ctx, "0.0.0.0:8888", rpc.Handler{})
	if err != nil {
		log.Fatalln(err)
	}
}
