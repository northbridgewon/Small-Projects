package main
import (
	"fmt"
	"math/rand"
	"time"
)

func statements() []string {
	return []string{
		"Hello, World!",
		"Go is awesome!",
		"Random statements are fun!",
		"Keep coding!",
		"Stay curious!",
		"Learning is a journey.",
		"Embrace challenges!",
		"Code is poetry.",
		"Make it happen!",
		"Believe in yourself!",
		"Stay positive!",
		"Keep pushing forward!",
		"Success is a journey, not a destination.",
		"Every day is a new opportunity.",
		"Dream big, work hard.",
		"Stay hungry, stay smart.",
		"Your only limit is your mind.",
		"Success is not for the lazy.",
		"Don't watch the clock; do what it does. Keep going.",
	}
}

func randomStatement() string {
	rand.Seed(time.Now().UnixNano())
	statementsList := statements()
	randomIndex := rand.Intn(len(statementsList))
	return statementsList[randomIndex]
}

func printRandomStatement() {
	for {
		fmt.Println(randomStatement())
		time.Sleep(0.5 * time.Second) // Wait for 2 seconds before printing the next statement
	}
}

func main() {
	fmt.Println("Random Statements Generator")
	fmt.Println("Press Ctrl+C to stop.")
	printRandomStatement()
	// This will never be reached because of the infinite loop in printRandomStatement
	fmt.Println("Goodbye!")
}