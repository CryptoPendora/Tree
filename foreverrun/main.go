package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go printTestContinuously()
	gracefulTermination()
}

// printTestContinuously prints "test" followed by a number and resets if it overflows.
// It rewrites the same line to give the illusion of updating the message.
func printTestContinuously() {
	counter := 0
	for {
		if counter > 9999 { // Reset counter to avoid overflow, number can be adjusted as needed.
			counter = 0
		}
		fmt.Printf("\rtest %d  \t\thit ctrl+c to exit", counter) // '\r' returns the cursor to the beginning of the line.
		counter++
		time.Sleep(1 * time.Second) // Sleep for a second to avoid spamming too fast.
	}
}

// gracefulTermination handles graceful termination of the program.
func gracefulTermination() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Printf("\rProgram exiting gracefully, hit ctrl+c to exit\n")
}
