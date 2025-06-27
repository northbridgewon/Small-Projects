package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	questions := map[int]string{
		0:  "What is the value of pi?",
		1:  "What is a mammal that lays eggs?",
		2:  "What color is the sky?",
		3:  "The best place to live in the USA?",
		4:  "What is the capital of France?",
		5:  "What is the largest planet in our solar system?",
		6:  "What is the chemical symbol for gold?",
		7:  "What is the speed of light in vacuum?",
		8:  "What is the tallest mountain in the world?",
		9:  "What is the largest ocean on Earth?",
		10: "What is the smallest prime number?",
	}

	answers := map[int]string{
		0:  "3.14159",
		1:  "Platypus",
		2:  "Blue",
		3:  "Tri State Area",
		4:  "Paris",
		5:  "Jupiter",
		6:  "Au",
		7:  "299,792,458 m/s",
		8:  "Mount Everest",
		9:  "Pacific Ocean",
		10: "2",
	}

	rand.Seed(time.Now().UnixNano())
	n := 10
	randomInt := rand.Intn(n)

	var userAnswer string
	fmt.Println("Welcome to the Quiz Game!")
	fmt.Println("Answer the following question:")
	fmt.Println(questions[randomInt])
	fmt.Println("Your answer is:")
	fmt.Scanln(&userAnswer)
	fmt.Println("You Answered:", userAnswer)
	fmt.Println("Correct answer is:", answers[randomInt])
}
