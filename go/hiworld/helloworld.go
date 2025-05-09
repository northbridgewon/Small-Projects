package main

import (
    "fmt"
    "time"
)

func produce(ch chan<- string) {
    msgs := []string{"Hello", "World", "from", "channels"}
    for _, msg := range msgs {
        ch <- msg // Send message to channel
        time.Sleep(time.Millisecond * 250) // Simulate work
    }
    close(ch) // Close channel once done sending
}

func consume(ch <-chan string) {
    for msg := range ch {
        fmt.Println(msg)
    }
}

func main() {
    ch := make(chan string) // Create a new channel

    go produce(ch) // Start producer goroutine
    go consume(ch) // Start consumer goroutine

    time.Sleep(time.Second) // Wait for goroutines to finish
    fmt.Println("Main program finished")
}
