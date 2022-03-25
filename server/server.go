package main

//go:generate protoc --go_out=../grpc/ --go-grpc_out=../grpc/ -I ../proto/ ../proto/chaton.proto

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"github.com/keftcha/chaton/grpc/chaton"
)

// HOST address for the server
var HOST string

// PORT the server listen on
var PORT int64

func init() {
	HOST = os.Getenv("HOST")
	// Set a default value
	if HOST == "" {
		HOST = "0.0.0.0"
	}

	var err error
	PORT, err = strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	// Set a default value
	if PORT == 0 && err.Error() == `strconv.ParseInt: parsing "": invalid syntax` {
		PORT = 21617
	} else if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialise the listening connections
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started: " + lis.Addr().String())

	// Create and start the gRPC server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chaton.RegisterChatonServer(grpcServer, newChatonServer())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
