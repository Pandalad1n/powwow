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
	difficulty := flag.Int("difficulty", 10, "pow difficulty")
	flag.Parse()

	ctx := context.Background()
	fmt.Println("Server starting")
	err := tcp.ListenAndServe(ctx, fmt.Sprintf("%s:%v", *listen, *port), rpc.Handler{Difficulty: uint32(*difficulty)})
	if err != nil {
		log.Fatalln(err)
	}
}
