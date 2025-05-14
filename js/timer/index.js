/**
 * @function countdownTimer
 * @description Starts a countdown timer for a specified duration.
 * @param {number} durationInSeconds - The total duration of the timer in seconds.
 */
function countdownTimer(durationInSeconds) {
    // Validate the input
    if (isNaN(durationInSeconds) || durationInSeconds <= 0) {
        console.error("Error: Please provide a positive number for the duration in seconds.");
        console.log("Usage: node timer.js <seconds>");
        return;
    }

    let remainingTime = Math.floor(durationInSeconds);

    // Start message
    console.log(`Timer started for ${remainingTime} seconds.`);

    // The interval function will run every second (1000 milliseconds)
    const timerInterval = setInterval(() => {
        // Calculate minutes and seconds for display
        const minutes = Math.floor(remainingTime / 60);
        const seconds = remainingTime % 60;

        // Format the time string (e.g., 02:05)
        const displayMinutes = String(minutes).padStart(2, '0');
        const displaySeconds = String(seconds).padStart(2, '0');

        // Clear the previous line and write the new time
        // process.stdout.write moves the cursor to the beginning of the line and clears it
        process.stdout.write(`\rTime remaining: ${displayMinutes}:${displaySeconds}   `); // Extra spaces to clear previous longer lines

        // Decrement the time
        remainingTime--;

        // Check if the timer has finished
        if (remainingTime < 0) {
            clearInterval(timerInterval); // Stop the interval
            process.stdout.write("\rTimer finished!             \n"); // Clear the line and print final message
            // You can add a sound or other notification here if desired
            // For example, console.log('\x07'); // Beep sound (might not work on all terminals)
        }
    }, 1000); // 1000 milliseconds = 1 second
}

// --- Main execution ---

// Get the duration from command-line arguments
// process.argv is an array:
// process.argv[0] is 'node'
// process.argv[1] is the path to your script (e.g., 'timer.js')
// process.argv[2] is the first actual argument
const args = process.argv.slice(2); // Get all arguments after 'node' and script path

if (args.length === 0) {
    console.log("Usage: node timer.js <seconds>");
    console.log("Example: node timer.js 60  (for a 1-minute timer)");
} else {
    const duration = parseInt(args[0], 10);
    countdownTimer(duration);
}
