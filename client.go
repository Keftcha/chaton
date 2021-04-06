package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	g "github.com/keftcha/chaton/generated"
)

func makeRequest(c g.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grtResp, err := c.Greeting(ctx, &g.GreetingsRequest{Name: name})
	fmt.Println(grtResp, err)
	fmt.Println(grtResp.GetMsg())
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial("localhost:21617", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := g.NewGreeterClient(conn)
	makeRequest(client, "Alban")
}
