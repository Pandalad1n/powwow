package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Pandalad1n/powwow/cmd/server/rpc"
	"github.com/Pandalad1n/powwow/tcp"
	"log"
)

func main() {
	listen := flag.String("listen", "", "server listen")
	port := flag.Int("port", 8888, "server port")
	flag.Parse()

	ctx := context.Background()
	err := tcp.ListenAndServe(ctx, fmt.Sprintf("%s:%v", *listen, *port), rpc.Handler{})
	if err != nil {
		log.Fatalln(err)
	}
}
