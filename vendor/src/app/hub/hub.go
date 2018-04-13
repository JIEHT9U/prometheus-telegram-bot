// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import "app/logger"

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func New() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

//Run start ws loop
func (h *Hub) Run(stopCh <-chan struct{}, l *logger.Logger) {
	go func() {
		l.InfoEntry().Info("Run Ws hub")
		defer l.InfoEntry().Info("Graceful shutdown Ws hub")
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
				l.InfoEntry().Infof("Add new ws client %s in pull", client.conn.LocalAddr().String())
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					l.InfoEntry().Infof("Remove  ws client %s in pull", client.conn.LocalAddr().String())
					delete(h.clients, client)
					close(client.send)
				}
			case message := <-h.broadcast:
				for client := range h.clients {
					select {
					case client.send <- message:
						l.InfoEntry().Info("Send broadcast")
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			case <-stopCh:
				break
			}
		}
	}()
}
