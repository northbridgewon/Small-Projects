package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	dir := flag.String("dir", ".", "Directory to serve static files from")

	flag.Parse()

	// Check if the directory exists
	info, err := os.Stat(*dir)
	if os.IsNotExist(err) {
		log.Fatalf("Error: Directory %q does not exist.\n", *dir)
	} else if !info.IsDir() {
		log.Fatalf("Error: %q is not a directory.\n", *dir)
	}

	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Serving static files from directory %q on port %d\n", *dir, *port)
	log.Fatal(http.ListenAndServe(addr, nil))
}