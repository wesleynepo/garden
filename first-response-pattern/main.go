package main

import (
	"fmt"
	"time"

	"math/rand/v2"
)

// Avoid timeout - Replicate servers and use first response
type Result string
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q", kind, query))
	}
}

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

func main() {
	start := time.Now()
	result := First("golang",
		fakeSearch("replica1"),
		fakeSearch("replica2"),
		fakeSearch("replica3"))
	elapsed := time.Since(start)
	fmt.Println(result)
	fmt.Printf("%.3fms\n", float64(elapsed.Nanoseconds())/1e6)
}
