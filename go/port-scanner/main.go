package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	host := flag.String("host", "localhost", "Host to scan (IP address or hostname)")
	startPort := flag.Int("start", 1, "Starting port for the scan")
	endPort := flag.Int("end", 1024, "Ending port for the scan")
	concurrency := flag.Int("c", 100, "Number of concurrent goroutines")

	flag.Parse()

	if *startPort < 1 || *endPort > 65535 || *startPort > *endPort {
		fmt.Println("Error: Invalid port range. Ports must be between 1 and 65535, and start port must be less than or equal to end port.")
		os.Exit(1)
	}

	fmt.Printf("Scanning %s from port %d to %d with %d concurrent connections...\n", *host, *startPort, *endPort, *concurrency)

	openPorts := make(chan int)
	go func() {
		for port := *startPort; port <= *endPort; port++ {
			go scanPort(*host, port, openPorts)
		}
	}()

	var results []int
	for i := *startPort; i <= *endPort; i++ {
		port := <-openPorts
		if port != 0 {
			results = append(results, port)
		}
	}

	sort.Ints(results)

	if len(results) > 0 {
		fmt.Println("\nOpen ports:")
		for _, port := range results {
			fmt.Printf("  %d\n", port)
		}
	} else {
		fmt.Println("\nNo open ports found in the specified range.")
	}
}

func scanPort(host string, port int, openPorts chan<- int) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		openPorts <- 0 // Port is closed or filtered
		return
	}
	conn.Close()
	openPorts <- port // Port is open
}
