//queue

package main

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Process()
}

type SimpleTask struct {
	ID int
}

func (t *SimpleTask) Process() {
	fmt.Printf("Processing task %d\n", t.ID)
	time.Sleep(2 * time.Second)
}

type WorkerPool struct {
	TasksChan   chan Task
	Wg          sync.WaitGroup
}

func NewWorkerPool(concurrency int) *WorkerPool {
	return &WorkerPool{
		TasksChan:   make(chan Task),
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.Wg.Add(1)
	go func() {
		wp.TasksChan <- task
	}()
}

func (wp *WorkerPool) worker(workerID int) {
	for task := range wp.TasksChan {
		fmt.Printf("Worker %d: ", workerID)
		task.Process()
		wp.Wg.Done()
	}
}

func (wp *WorkerPool) Run(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go wp.worker(i + 1)
	}
}

func (wp *WorkerPool) Close() {
	close(wp.TasksChan)
	wp.Wg.Wait()
}

func main() {
	wp := NewWorkerPool(3)

	wp.Run(3) // Start 3 workers

	// Add tasks
	for i := 1; i <= 5; i++ {
		wp.AddTask(&SimpleTask{ID: i})
	}

	wp.Close()

	fmt.Println("All tasks have been processed!")
}
