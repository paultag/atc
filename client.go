package main

import (
	"io"
	"log"
)

func NewClient(w io.WriteCloser) *Client {
	return &Client{Writer: w}
}

type Client struct {
	Writer io.WriteCloser
}

func NewClients() *Clients {
	return &Clients{Clients: map[*Client]bool{}}
}

type Clients struct {
	Clients map[*Client]bool
}

func (c *Clients) Register(client *Client) {
	c.Clients[client] = true
}

func (c *Clients) Unregister(client *Client) {
	if _, ok := c.Clients[client]; ok {
		delete(c.Clients, client)
		client.Writer.Close()
	}
}

func (c *Clients) Broadcast(msg []byte) {
	for client := range c.Clients {
		_, err := client.Writer.Write(msg)
		if err != nil {
			log.Printf("Error: %s", err)
			c.Unregister(client)
		}
	}
}
