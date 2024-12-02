import java.util.InputMismatchException;
import java.util.Scanner;

public class Quiz {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        int score = 0;

        // Define the questions, options, and correct answers
        String[] questions = {
            "What is the largest organ in the human body?",
            "What is the largest country by land area?",
            "What is 2 + 2?",
            "Who wrote 'To Kill a Mockingbird'?",
            "What is the chemical symbol for water?"
        };

        String[][] options = {
            {"1) Mesentery", "2) Skin", "3) Brain", "4) Heart"},
            {"1) China", "2) Russia", "3) Canada", "4) USA"},
            {"1) 3", "2) 4", "3) 5", "4) 6"},
            {"1) Mark Twain", "2) Harper Lee", "3) Ernest Hemingway", "4) F. Scott Fitzgerald"},
            {"1) H2O", "2) CO2", "3) O2", "4) NaCl"}
        };

        int[] answers = {2, 2, 2, 2, 1};

        // Loop through each question
        for (int i = 0; i < questions.length; i++) {
            System.out.println("\nQuestion " + (i + 1) + ": " + questions[i]);
            for (String option : options[i]) {
                System.out.println(option);
            }
            System.out.print("Your answer (1, 2, 3, or 4): ");
            
            int answer = 0;
            boolean validInput = false;
            
            while (!validInput) {
                try {
                    answer = scanner.nextInt();
                    if (answer < 1 || answer > 4) {
                        System.out.print("Invalid input. Please enter a number between 1 and 4: ");
                    } else {
                        validInput = true;
                    }
                } catch (InputMismatchException e) {
                    System.out.print("Invalid input. Please enter a number between 1 and 4: ");
                    scanner.next(); // Clear the invalid input
                }
            }

            // Check if the answer is correct
            if (answer == answers[i]) {
                System.out.println("Correct!");
                score++;
            } else {
                System.out.println("Incorrect. The correct answer is " + answers[i] + ".");
            }
        }

        // Display the final score
        System.out.println("\nYou scored " + score + " out of " + questions.length + ".");
        scanner.close();
    }
}