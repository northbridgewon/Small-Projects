package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Question represents a multiple-choice question
type Question struct {
	Text    string
	Options []string
	Answer  string
}

func main() {
	// Define the 5 questions
	questions := []Question{
		{
			Text: "What is the capital of France?",
			Options: []string{
				"Berlin",
				"Paris",
				"London",
				"Rome",
			},
			Answer: "Paris",
		},
		{
			Text: "Which planet is known as the Red Planet?",
			Options: []string{
				"Earth",
				"Mars",
				"Jupiter",
				"Saturn",
			},
			Answer: "Mars",
		},
		{
			Text: "Who painted the Starry Night?",
			Options: []string{
				"Leonardo da Vinci",
				"Vincent van Gogh",
				"Pablo Picasso",
				"Claude Monet",
			},
			Answer: "Vincent van Gogh",
		},
		{
			Text: "What is the largest planet in our solar system?",
			Options: []string{
				"Jupiter",
				"Saturn",
				"Uranus",
				"Neptune",
			},
			Answer: "Jupiter",
		},
		{
			Text: "Which programming language are we using?",
			Options: []string{
				"Python",
				"Java",
				"C++",
				"Go",
			},
			Answer: "Go",
		},
	}

	// Initialize score
	score := 0

	// Loop through each question
	for i, q := range questions {
		fmt.Printf("\nQuestion %d: %s\n", i+1, q.Text)
		for j, option := range q.Options {
			fmt.Printf("%d. %s\n", j+1, option)
		}

		// Get user's answer
		var userAnswer string
		var selectedOption string
		fmt.Print("Enter the number of your answer: ")
		fmt.Scanln(&userAnswer)
		// Validate input (simple, does not handle non-integer inputs)
		if num, err := strconv.Atoi(userAnswer); err == nil {
			if num >= 1 && num <= len(q.Options) {
				selectedOption = q.Options[num-1]
			} else {
				fmt.Println("Invalid option. Moving to the next question.")
				continue
			}
		} else {
			fmt.Println("Invalid input. Please enter a number. Moving to the next question.")
			continue
		}

		// Check if the answer is correct
		if strings.TrimSpace(selectedOption) == strings.TrimSpace(q.Answer) {
			fmt.Println("Correct!")
			score++
		} else {
			fmt.Printf("Incorrect. The answer was %s.\n", q.Answer)
		}
	}

	// Display final score
	fmt.Printf("\nQuiz Over! Your final score is %d out of %d\n", score, len(questions))
}
