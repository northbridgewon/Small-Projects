package main
import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a slice of strings
	statements := []string{
		"我爱你",
		"I love you 9000!",
		"We do not give up, we do not give in!",
		"Love is more than an emotion, it is a promise, a commitment and a choice.",
		"Love is not about possession, it's about appreciation.",
		"There is no right or wrong in love, only the truth of your heart.",
		"These are things that trancend the boundaries of species.",
	}

	// Print a random statement from the slice
	for {
		fmt.Println(statements[rand.Intn(len(statements))])
		time.Sleep(1 * time.Second) // Wait for a second before printing the next statement
	}
}