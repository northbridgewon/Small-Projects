import java.util.Scanner;

public class StudentRecordSystem {

    // Static variables to store student data
    private static final int MAX_STUDENTS = 100;
    private static String[] studentIds = new String[MAX_STUDENTS];
    private static String[] names = new String[MAX_STUDENTS];
    private static int[] ages = new int[MAX_STUDENTS];
    private static String[] grades = new String[MAX_STUDENTS];
    private static int studentCount = 0;

    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        int choice;

        do {
            displayMenu();
            choice = scanner.nextInt();
            scanner.nextLine();  // Consume newline

            switch (choice) {
                case 1:
                    addStudent(scanner);
                    break;
                case 2:
                    updateStudent(scanner);
                    break;
                case 3:
                    viewStudent(scanner);
                    break;
                case 4:
                    viewAllStudents();
                    break;
                case 5:
                    System.out.println("Exiting...");
                    break;
                default:
                    System.out.println("Invalid choice. Please try again.");
            }
        } while (choice != 5);

        scanner.close();
    }

    // Method to display the menu
    private static void displayMenu() {
        System.out.println("\nStudent Management System");
        System.out.println("1. Add Student");
        System.out.println("2. Update Student");
        System.out.println("3. View Student");
        System.out.println("4. View All Students");
        System.out.println("5. Exit");
        System.out.print("Enter your choice: ");
    }

    // Method to add a new student
    private static void addStudent(Scanner scanner) {
        if (studentCount >= MAX_STUDENTS) {
            System.out.println("Cannot add more students. Maximum capacity reached.");
            return;
        }
        System.out.print("Enter Student ID: ");
        String studentId = scanner.nextLine();
        System.out.print("Enter Student Name: ");
        String name = scanner.nextLine();
        System.out.print("Enter Student Age: ");
        int age = scanner.nextInt();
        scanner.nextLine();  // Consume newline
        System.out.print("Enter Student Grade: ");
        String grade = scanner.nextLine();
        studentIds[studentCount] = studentId;
        names[studentCount] = name;
        ages[studentCount] = age;
        grades[studentCount] = grade;
        studentCount++;
        System.out.println("Student added successfully.");
    }

    // Method to update an existing student's information
    private static void updateStudent(Scanner scanner) {
        System.out.print("Enter Student ID to update: ");
        String studentId = scanner.nextLine();
        for (int i = 0; i < studentCount; i++) {
            if (studentIds[i].equals(studentId)) {
                System.out.print("Enter Student Name: ");
                String name = scanner.nextLine();
                System.out.print("Enter Student Age: ");
                int age = scanner.nextInt();
                scanner.nextLine();  // Consume newline
                System.out.print("Enter Student Grade: ");
                String grade = scanner.nextLine();
                names[i] = name;
                ages[i] = age;
                grades[i] = grade;
                System.out.println("Student updated successfully.");
                return;
            }
        }
        System.out.println("Student not found.");
    }

    // Method to view a specific student's details
    private static void viewStudent(Scanner scanner) {
        System.out.print("Enter Student ID to view: ");
        String studentId = scanner.nextLine();
        for (int i = 0; i < studentCount; i++) {
            if (studentIds[i].equals(studentId)) {
                System.out.println("Student ID: " + studentIds[i] + ", Name: " + names[i] + ", Age: " + ages[i] + ", Grade: " + grades[i]);
                return;
            }
        }
        System.out.println("Student not found.");
    }

    // Method to view all students' details
    private static void viewAllStudents() {
        if (studentCount == 0) {
            System.out.println("No students to display.");
            return;
        }
        for (int i = 0; i < studentCount; i++) {
            System.out.println("Student ID: " + studentIds[i] + ", Name: " + names[i] + ", Age: " + ages[i] + ", Grade: " + grades[i]);
        }
    }
}