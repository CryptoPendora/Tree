// 1 การจัดสรรขนาดของแชนแนล (Channel Buffering):

// 2 การใช้ non-buffered channel (make(chan Task)) ทำให้การส่ง task ไปยัง channel ต้องรอจนกว่าจะมี goroutine อื่นรับข้อมูลจาก channel นั้น สิ่งนี้อาจทำให้เกิด latency ในการประมวลผลงานหากมีงานส่งเข้ามาพร้อมกันจำนวนมาก
// แนะนำ: พิจารณาใช้ buffered channel (เช่น make(chan Task, 100)) เพื่อให้สามารถรองรับการส่งงานเข้ามาได้หลายงานโดยไม่ต้อง block ทันที ซึ่งช่วยลด latency และช่วยให้ระบบมีประสิทธิภาพการทำงานที่ดีขึ้น
// การปรับขนาดของ Worker:

// 3 การปรับขนาด worker ให้เหมาะสมกับปริมาณงานเป็นสิ่งสำคัญ เพื่อให้สามารถประมวลผลงานได้อย่างมีประสิทธิภาพโดยไม่ให้ทรัพยากรว่างเปล่าหรือใช้ทรัพยากรมากเกินไป
// แนะนำ: พิจารณาเพิ่มกลไกสำหรับการปรับขนาดจำนวน worker โดยอัตโนมัติตามปริมาณงานที่มี เช่น เพิ่ม worker ในช่วงที่มีงานเข้ามามาก และลดจำนวน worker ลงเมื่องานน้อยลง
// การจัดการกับ Backpressure:

// เมื่อใช้ buffered channel, ยังคงมีความเป็นไปได้ที่ channel จะเต็มหากมีงานส่งเข้ามาในอัตราที่รวดเร็วกว่าที่ worker สามารถประมวลผลได้
// แนะนำ: พิจารณาใช้กลไก backpressure เช่น การจำกัดอัตราการส่งงานหรือใช้ queue ภายนอกสำหรับการจัดคิวงาน เพื่อป้องกันไม่ให้ระบบโอเวอร์โหลดและให้สามารถรับมือกับสถานการณ์ที่ channel เต็มได้

// queue//queue

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	TasksChan chan Task
	Wg        sync.WaitGroup
	ctx       context.Context
	cancel    func()
	shutdown  chan struct{}
}

func NewWorkerPool(concurrency int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		TasksChan: make(chan Task),
		ctx:       ctx,
		cancel:    cancel,
		shutdown:  make(chan struct{}),
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	select {
	case <-wp.shutdown:
		fmt.Println("No longer accepting new tasks.")
		return
	default:
		wp.Wg.Add(1)
		go func() {
			defer wp.Wg.Done()
			select {
			case wp.TasksChan <- task:
			case <-wp.ctx.Done():
			}
		}()
	}
}

func (wp *WorkerPool) worker(workerID int) {
	for {
		select {
		case task, ok := <-wp.TasksChan:
			if !ok {
				return // Exit the worker if the channel is closed
			}
			fmt.Printf("Worker %d: ", workerID)
			task.Process()
		case <-wp.ctx.Done():
			return // Exit the worker if context is cancelled
		}
	}
}

func (wp *WorkerPool) Run(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go wp.worker(i + 1)
	}
}

func (wp *WorkerPool) Close() {
	close(wp.shutdown)  // Signal that no more tasks should be added
	wp.Wg.Wait()        // Wait for all tasks to be processed before closing the channel
	wp.cancel()         // Cancel the context
	close(wp.TasksChan) // It's now safe to close this channel since all workers have exited and all tasks have been processed
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

// printTestContinuously()

// gracefulTermination()

// }

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

//untested

// func NewWorkerPool(concurrency int, bufferSize int) *WorkerPool {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	return &WorkerPool{
// 		TasksChan: make(chan Task, bufferSize), // ใช้ buffered channel
// 		ctx:       ctx,
// 		cancel:    cancel,
// 		shutdown:  make(chan struct{}),
// 	}
// }
