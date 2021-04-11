package types

import (
	"fmt"
)

// Clients represente a slice of *Client
type Clients []*Client

// Broadcasting send the event to all clients
func (cs *Clients) Broadcasting(e Event) {
	for _, c := range *cs {
		// Put the sender name as the message author
		e.Event.Msg.Author = e.Client.Nick
		// Remove client if there is an error sending him message
		if err := c.Stream.Send(e.Event); err != nil {
			cs.Remove(e.Client)
		}
	}
}

// Remove a client from our *Client slice
func (cs *Clients) Remove(client *Client) {
	// Remove the client of our list
	for i, c := range *cs {
		if c == client {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			break
		}
	}
}

// Add a client from our *Client slice
func (cs *Clients) Add(client *Client) {
	*cs = append(*cs, client)
}

// ListClients and their status
func (cs *Clients) ListClients() string {
	lst := ""
	for i, c := range *cs {
		if i != 0 {
			lst += "\n"
		}
		lst += fmt.Sprintf(
			"- %s: %s",
			c.Nick,
			c.Status,
		)
	}
	return lst
}
