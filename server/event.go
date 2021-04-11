package main

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

type event struct {
	client *client
	event  *chaton.Event
}
