package LibrarySystem;
import java.util.HashMap;
import java.util.Map;
import java.util.Scanner;

public class LibrarySystem {
    private Map<String, Integer> books;

    public LibrarySystem() {
        this.books = new HashMap<>();
    }

    public void addBook(String title, int quantity) {
        if (books.containsKey(title)) {
            books.put(title, books.get(title) + quantity);
            System.out.println("Quantity of " + title + " updated.");
        } else {
            books.put(title, quantity);
            System.out.println(title + " added to the library.");
        }
    }

    public void borrowBook(String title, int quantity) {
        if (books.containsKey(title) && books.get(title) >= quantity) {
            books.put(title, books.get(title) - quantity);
            System.out.println(quantity + " copies of " + title + " borrowed successfully.");
        } else {
            System.out.println("Error: Not enough copies of " + title + " available.");
        }
    }

    public void returnBook(String title, int quantity) {
        if (books.containsKey(title)) {
            books.put(title, books.get(title) + quantity);
            System.out.println(quantity + " copies of " + title + " returned successfully.");
        } else {
            System.out.println("Error: " + title + " does not belong to this library.");
        }
    }

    public void displayMenu() {
        System.out.println("----------------------------------------------------------------------------------------------------------");
        System.out.println("Press 1 to Add new Book.");
        System.out.println("Press 2 to Borrow a Book.");
        System.out.println("Press 3 to Return a Book.");
        System.out.println("Press 4 to Exit.");
        System.out.println("----------------------------------------------------------------------------------------------------------");
    }

    public void run() {
        Scanner input = new Scanner(System.in);
        int choice;
        do {
            displayMenu();
            System.out.print("Enter your choice: ");
            choice = input.nextInt();
            input.nextLine(); // Consume the newline

            switch (choice) {
                case 1:
                    System.out.print("Enter Book Title: ");
                    String titleAdd = input.nextLine();
                    System.out.print("Enter Quantity: ");
                    int quantityAdd = input.nextInt();
                    addBook(titleAdd, quantityAdd);
                    break;
                case 2:
                    System.out.print("Enter Book Title: ");
                    String titleBorrow = input.nextLine();
                    System.out.print("Enter Quantity to Borrow: ");
                    int quantityBorrow = input.nextInt();
                    borrowBook(titleBorrow, quantityBorrow);
                    break;
                case 3:
                    System.out.print("Enter Book Title: ");
                    String titleReturn = input.nextLine();
                    System.out.print("Enter Quantity to Return: ");
                    int quantityReturn = input.nextInt();
                    returnBook(titleReturn, quantityReturn);
                    break;
                case 4:
                    System.out.println("Exiting the Library System.");
                    break;
                default:
                    System.out.println("Invalid choice. Please try again.");
            }
        } while (choice != 4);
        input.close();
    }

    public static void main(String[] args) {
        System.out.println("********************Welcome to the UoP Library!********************");
        LibrarySystem library = new LibrarySystem();
        library.run();
    }
}