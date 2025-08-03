package main

import (
	"fmt"
	"time"

	"math/rand/v2"
)

func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.IntN(1e3)) * time.Millisecond)
	}
}

func main() {
	// unbuffered channel
	c := make(chan string)
	go boring("boring!", c)
	for range 10 {
		// <-c blocks the execution until
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")
}
