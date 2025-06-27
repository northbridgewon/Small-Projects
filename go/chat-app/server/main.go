package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Client represents a single chat client connected via WebSocket
type Client struct {
	ws *websocket.Conn
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered: %s\n", client.ws.RemoteAddr())
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered: %s\n", client.ws.RemoteAddr())
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.ws.Close()
	}()
	for {
		var msg string
		err := websocket.Message.Receive(c.ws, &msg)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Read error: %v\n", err)
			}
			break
		}
		log.Printf("Received: %s from %s\n", msg, c.ws.RemoteAddr())
		hub.broadcast <- []byte(msg)
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.ws.Close()
	}()
	for message := range c.send {
		_, err := c.ws.Write(message)
		if err != nil {
			log.Printf("Write error: %v\n", err)
			return
		}
	}
}

// wsHandler handles websocket connections
func wsHandler(hub *Hub, ws *websocket.Conn) {
	client := &Client{ws: ws, send: make(chan []byte, 256)}
	hub.register <- client

	go client.writePump()
	client.readPump(hub)
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		wsHandler(hub, ws)
	}))

	port := "8080"
	fmt.Printf("Chat server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
