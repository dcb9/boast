package ws

import (
	"github.com/dcb9/boast/transaction"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
		case tx := <-transaction.AddChannel:
			tss := []*transaction.Tx{tx}
			for client := range h.clients {
				select {
				case client.send <- tss:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
