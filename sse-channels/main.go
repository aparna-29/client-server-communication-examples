package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Client represents an SSE client connection
type Client struct {
	// Each client has its own channel to receive messages
	messageChan chan string
}

// Broker manages all connected clients and broadcasts messages
type Broker struct {
	// Register new clients
	clients map[*Client]bool

	// Channel for new client registrations
	register chan *Client

	// Channel for client unregister (disconnect)
	unregister chan *Client

	// Channel to broadcast messages to all clients
	broadcast chan string
}

func newBroker() *Broker {
	return &Broker{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan string),
	}
}

func (b *Broker) run() {
	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
			log.Println("Client registered. Total clients:", len(b.clients))

		case client := <-b.unregister:
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.messageChan)
				log.Println("Client unregistered. Total clients:", len(b.clients))
			}

		case message := <-b.broadcast:
			// Send message to all connected clients
			for client := range b.clients {
				select {
				case client.messageChan <- message:
				default:
					// If client channel is blocked, unregister it
					delete(b.clients, client)
					close(client.messageChan)
				}
			}
		}
	}
}

func sseHandler(b *Broker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		// Create new client and register it
		client := &Client{messageChan: make(chan string)}
		b.register <- client

		// Remove client when handler exits
		defer func() {
			b.unregister <- client
		}()

		// Notify client about connection
		fmt.Fprintf(w, "data: Connected\n\n")
		flusher.Flush()

		// Listen to connection close by client
		notify := w.(http.CloseNotifier).CloseNotify()

		for {
			select {
			case msg, ok := <-client.messageChan:
				if !ok {
					// Channel closed
					return
				}

				// Write message as SSE
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()

			case <-notify:
				// Client disconnected
				return
			}
		}
	}
}

func main() {
	broker := newBroker()
	go broker.run()

	// Broadcast a message every 2 seconds
	go func() {
		counter := 0
		for {
			message := fmt.Sprintf("Broadcast message %d at %s", counter, time.Now().Format(time.RFC3339))
			broker.broadcast <- message
			counter++
			time.Sleep(2 * time.Second)
		}
	}()

	http.HandleFunc("/events", sseHandler(broker))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
