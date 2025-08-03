package main

import (
	"context"
	"fmt"
	"time"
)

func slowOperation(result chan<- string, delay time.Duration) {
	time.Sleep(delay)
	result <- "operation complete"
}

func timeoutExample() {
	ch := make(chan string)

	// Start slow operation
	go slowOperation(ch, 2*time.Second)

	// Wait with timeout
	select {
	case result := <-ch:
		fmt.Println("Got:", result)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout!")
	}
}

func timeoutExampleCtx() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ch := make(chan string)

	select {
	case result := <-ch:
		fmt.Println("Got:", result)
	case <-ctx.Done():
		fmt.Println("Timeout! Error:", ctx.Err())
	}
}

func main() {
	timeoutExample()
	timeoutExampleCtx()
}
