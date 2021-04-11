package main

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
)

// client represent a client connetion to the server
type client struct {
	stream chaton.Chaton_JoinServer
	nick   string
	status string
}

// clients represente a slice of *Client
type clients []*client

// broadcasting send the event to all clients
func (cs *clients) broadcasting(e event) {
	for _, c := range *cs {
		// Put the sender name as the message author
		e.event.Msg.Author = e.client.nick
		// Remove client if there is an error sending him message
		if err := c.stream.Send(e.event); err != nil {
			cs.remove(e.client)
		}
	}
}

// remove a client from our *Client slice
func (cs *clients) remove(client *client) {
	// Remove the client of our list
	for i, c := range *cs {
		if c == client {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			break
		}
	}
}

// add a client from our *Client slice
func (cs *clients) add(client *client) {
	*cs = append(*cs, client)
}

// listClients and their status
func (cs *clients) listClients() string {
	lst := ""
	for i, c := range *cs {
		if i != 0 {
			lst += "\n"
		}
		lst += fmt.Sprintf(
			"- %s: %s",
			c.nick,
			c.status,
		)
	}
	return lst
}
