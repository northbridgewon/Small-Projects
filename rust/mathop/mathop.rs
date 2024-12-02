// Declare two variables

fn main() {
    let mut num1: i32 = 10;
    let num2: i32 = 5;

    let num3 = num1 + num2;
    println!("{} + {} = {}", num1, num2, num3);
    let num4 = num1 > num2;
    println!("{} > {} = {}", num1, num2, num4);
    let num5 = num1 > 0 && num2 > 0;
    println!("{} > 0 && {} > 0 = {}", num1, num2, num5);
    num1 *= 2;
    println!("{}",num1);
    
}
// Use operators to perform the following operations
// 1. Addition
// 2. Comparison (check if num1 is greater than num2)
// 3. Logical AND (check if both num1 and num2 are greater than 0)
// 4. Assignment (double the value of num1)

// Print the results of each operation