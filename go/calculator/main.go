package main

import (
	"bufio" // For buffered I/O, like reading from the console
	"fmt"   // For formatted I/O, like printing to the console
	"os"    // For interacting with the operating system, like handling signals
)

func additinon(a, b int) int {
	// This function takes two integers and returns their sum.
	// It is a simple addition function.
	return a + b
}

func subtract(a, b int) int {
	// This function takes two integers and returns their difference.
	// It is a simple subtraction function.
	return a - b
}

func multiply(a, b int) int {
	// This function takes two integers and returns their product.
	// It is a simple multiplication function.
	return a * b
}

func divide(a, b int) int {
	// This function takes two integers and returns their quotient.
	// It is a simple division function.
	if b == 0 {
		fmt.Println("Error: Division by zero")
		return 0
	}
	return a / b
}

func main() {

	// This is the main function where the program starts executing.
	// It will prompt the user for two numbers and an operator, then perform the calculation.

	// Create a new scanner to read input from the console.
	scanner := bufio.NewScanner(os.Stdin)

	// Prompt the user for the first number.
	fmt.Print("Enter first number: ")
	scanner.Scan() // Read the input
	var num1 int
	fmt.Sscan(scanner.Text(), &num1) // Convert the input to an integer

	// Prompt the user for the second number.
	fmt.Print("Enter second number: ")
	scanner.Scan() // Read the input
	var num2 int
	fmt.Sscan(scanner.Text(), &num2) // Convert the input to an integer

	// Prompt the user for an operator.
	fmt.Print("Enter operator (+, -, *, /): ")
	scanner.Scan()             // Read the input
	operator := scanner.Text() // Get the operator as a string

	// Perform the calculation based on the operator.
	switch operator {
	case "+":
		result := additinon(num1, num2)
		fmt.Printf("Result: %d\n", result)
	case "-":
		result := subtract(num1, num2)
		fmt.Printf("Result: %d\n", result)
	case "*":
		result := multiply(num1, num2)
		fmt.Printf("Result: %d\n", result)
	case "/":
		result := divide(num1, num2)
		fmt.Printf("Result: %d\n", result)

	default:
		fmt.Println("Invalid operator")
	}

}
