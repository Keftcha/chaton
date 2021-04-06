package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	g "github.com/keftcha/chaton/generated"
)

type GreeterServer struct {
	g.UnimplementedGreeterServer
}

func (s *GreeterServer) Greeting(ctx context.Context, name *g.GreetingsRequest) (*g.GreetingsResponse, error) {
	fmt.Println(name, name.GetName())
	return &g.GreetingsResponse{Msg: "Hello " + name.GetName()}, nil
}

func newGreeterServer() *GreeterServer {
	return &GreeterServer{}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 21617))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server started on port " + lis.Addr().String())

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	g.RegisterGreeterServer(grpcServer, newGreeterServer())
	grpcServer.Serve(lis)
}
