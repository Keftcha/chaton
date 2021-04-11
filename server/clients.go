package main

import (
	"fmt"
)

type clients []*client

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

func (cs *clients) remove(client *client) {
	// Remove the client of our list
	for i, c := range *cs {
		if c == client {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			break
		}
	}
}

func (cs *clients) add(client *client) {
	*cs = append(*cs, client)
}

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
