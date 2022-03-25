package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/keftcha/chaton/grpc/chaton"
)

// Client represent a client connetion to the server
type Client struct {
	Stream chaton.Chaton_JoinServer
	Nick   string
	Status string
}

// Clients represente connected client
type Clients map[*Client]struct{}

// Broadcasting send the event to all clients
func Broadcasting(cs Clients, e Event) {
	// Put the sender name as the message author
	e.Event.Msg.Author = e.Client.Nick

	for c := range cs {
		if err := c.Stream.Send(e.Event); err != nil {
			log.Printf("Err Sending event to Client `%#v` %s\n", c, err.Error())
		}
	}
}

// RemoveClient a client from our *Client slice
func RemoveClient(cs Clients, c *Client) {
	delete(cs, c)
}

// AddClient a client from our *Client slice
func AddClient(cs Clients, c *Client) {
	cs[c] = struct{}{}
}

// ListClients and their status
func ListClients(cs Clients) string {
	lst := make([]string, len(cs))
	i := 0
	for c := range cs {
		lst[i] = fmt.Sprintf(
			"- %s: %s",
			c.Nick,
			c.Status,
		)
		i++
	}
	return strings.Join(lst, "\n")
}
