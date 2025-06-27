package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

func checkWebsite(website string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	_, err := http.Get("https://" + website)
	if err != nil {
		results <- fmt.Sprintf("%s is down!", website)
		return
	}
	results <- fmt.Sprintf("%s is up!", website)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <website1> <website2> ...")
		return
	}

	websites := os.Args[1:]
	results := make(chan string)
	var wg sync.WaitGroup

	for _, website := range websites {
		wg.Add(1)
		go checkWebsite(website, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
