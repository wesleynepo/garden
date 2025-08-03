package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Task string
}

type Result struct {
	JobID  int
	Output string
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)

		// Simulate work
		time.Sleep(time.Second)

		// Send result
		results <- Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Job %d completed by worker %d", job.ID, id),
		}
	}
}

func main() {
	numWorkers := 3
	numJobs := 5

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Task: fmt.Sprintf("task-%d", j)}
	}
	close(jobs)

	go func() {
		wg.Wait()
		fmt.Printf("Released the beast")
		close(results)
	}()

	for result := range results {
		time.Sleep(time.Second)
		fmt.Println(result.Output)
	}
}
