use std::io;

fn main() {
    let int = 84;
    let float = 19.84;
    let mut name = String::new();
    println!("What's your name?");
    io::stdin().read_line(&mut name).expect("Failed to read line");
    println!("{} {}",int,float);
    println!("Hello, {}!", name.trim());
}