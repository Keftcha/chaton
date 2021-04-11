package main

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

type client struct {
	stream chaton.Chaton_JoinServer
	nick   string
	status string
}
