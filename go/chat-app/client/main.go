package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

func main() {
	// Determine the WebSocket URL
	serviceURL := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	// Connect to the WebSocket server
	ws, err := websocket.Dial(serviceURL.String(), "", serviceURL.String())
	if err != nil {
		log.Fatalf("Dial failed: %v\n", err)
	}
	defer ws.Close()

	fmt.Println("Connected to chat server. Type your message and press Enter.")

	// Goroutine to read messages from the server
	go func() {
		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Printf("Error receiving message: %v\n", err)
				return
			}
			fmt.Printf("\rReceived: %s\n> ", msg) // \r to clear current input line
		}
	}()

	// Read input from stdin and send to server
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		message := strings.TrimSpace(input)

		if message == "/quit" || message == "/exit" {
			fmt.Println("Exiting chat.")
			return
		}

		if message != "" {
			err := websocket.Message.Send(ws, message)
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
				return
			}
		}
	}
}
