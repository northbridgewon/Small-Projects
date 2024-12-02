package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1] // Remove the newline character
	fmt.Printf("Hello, %s!\n", name)
}
