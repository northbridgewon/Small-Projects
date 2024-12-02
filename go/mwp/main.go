package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	pst := now.Unix()
	fmt.Println("Current UNIX Time:", pst, now)
}
