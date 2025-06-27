package main

import (
	"distributed-kv-store/rpc"
	"distributed-kv-store/store"
	"flag"
	"fmt"
	"log"
	"net"
	stdrpc "net/rpc"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	peersStr := flag.String("peers", "", "Comma-separated list of peer addresses (e.g., localhost:8081,localhost:8082)")
	flag.Parse()

	// Initialize the local key-value store
	localStore := store.NewStore()

	// Register the RPC server
	kvRPC := rpc.NewKVStore(localStore)
	stdrpc.Register(kvRPC)

	// Start listening for RPC connections
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()

	log.Printf("Node listening on port %d\n", *port)

	// Handle peer connections (simplified replication)
	peers := []string{}
	if *peersStr != "" {
		peers = strings.Split(*peersStr, ",")
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}
			go stdrpc.ServeConn(conn)
		}
	}()

	// Simple client for testing (can be expanded into a separate client CLI)
	if len(flag.Args()) > 0 {
		command := flag.Args()[0]
		switch command {
		case "put":
			if len(flag.Args()) < 3 {
				fmt.Println("Usage: go run main.go -port <port> put <key> <value>")
				os.Exit(1)
			}
			key := flag.Args()[1]
			value := flag.Args()[2]

			// Simulate replication to peers
			for _, peerAddr := range peers {
				client, err := stdrpc.Dial("tcp", peerAddr)
				if err != nil {
					log.Printf("Error connecting to peer %s: %v\n", peerAddr, err)
					continue
				}
				defer client.Close()

				args := rpc.Args{Key: key, Value: value}
				var reply rpc.Reply
				err = client.Call("KVStore.Put", args, &reply)
				if err != nil {
					log.Printf("Error replicating Put to %s: %v\n", peerAddr, err)
				} else {
					log.Printf("Replicated Put to %s: %s\n", peerAddr, reply.Value)
				}
			}
			// Also put to local store
			localStore.Put(key, value)
			fmt.Printf("Put %s=%s locally and replicated.\n", key, value)

		case "get":
			if len(flag.Args()) < 2 {
				fmt.Println("Usage: go run main.go -port <port> get <key>")
				os.Exit(1)
			}
			key := flag.Args()[1]
			val, err := localStore.Get(key)
			if err != nil {
				fmt.Printf("Error getting key %s: %v\n", key, err)
			} else {
				fmt.Printf("Get %s=%s\n", key, val)
			}
		default:
			fmt.Println("Unknown command. Use 'put' or 'get'.")
		}
	}

	// Keep the main goroutine alive for the server to run
	select {}
}