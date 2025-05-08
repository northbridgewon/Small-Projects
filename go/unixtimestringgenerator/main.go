// main.go
package main

import (
	"fmt"      // For formatted I/O, like printing to the console
	"os"       // For interacting with the operating system, like handling signals
	"os/signal" // For specific signal handling capabilities
	"strconv"  // For string conversions, like int to string in different bases
	"syscall"  // Contains OS-level system calls, like signal types
	"time"     // For time-related functions
)

// stringGeneratorFunction is a type alias for functions that take a Unix timestamp (int64)
// and return a string.
type stringGeneratorFunction func(timestamp int64) string

// generateStringReadableUTC formats the Unix timestamp into a human-readable UTC date and time string.
func generateStringReadableUTC(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}

// generateStringSimple creates a basic string with "UTS_" prefix and the timestamp.
func generateStringSimple(timestamp int64) string {
	return fmt.Sprintf("UTS_%d", timestamp)
}

// generateStringBase16 converts the Unix timestamp (int64) to its base-16 (hexadecimal) string representation.
func generateStringBase16(timestamp int64) string {
	// strconv.FormatInt converts an integer (int64) to a string representation in the given base.
	// The second argument '16' specifies base-16 (hexadecimal).
	return strconv.FormatInt(timestamp, 16)
}

func main() {
	fmt.Println("Starting Go string generator (one per Unix second)...")
	fmt.Println("Press Ctrl+C to stop.")

	// --- 1. Choose your string generation function ---
	// You can now choose generateStringBase16 as well.
	// selectedGenerator := generateStringReadableUTC
	// selectedGenerator := generateStringSimple
	selectedGenerator := generateStringBase16 // <-- Select the new base-16 generator

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastProcessedSecond int64 = 0

	for {
		select {
		case sig := <-osSignals:
			fmt.Printf("\nSignal received: %s. Exiting gracefully...\n", sig)
			return

		case currentTime := <-ticker.C:
			currentUnixSecond := currentTime.Unix()

			if currentUnixSecond > lastProcessedSecond {
				generatedString := selectedGenerator(currentUnixSecond)
				fmt.Println(generatedString)
				lastProcessedSecond = currentUnixSecond
			}
		}
	}
}