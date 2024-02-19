//ReaderHeartbeatWarriors main.go

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	serverHeartbeatFile = "/app/shared/server_heartbeat.txt"
	readerHeartbeatFile = "/app/shared/reader_heartbeat.txt"
	heartbeatInterval   = 20 * time.Second
	restartDelay        = 100 * time.Second
	serverProgramPath   = "/app/ServerHeartbeatWarriors" // Adjusted for consistency
)

// Task represents a generic task with an ID
type Task struct {
	ID int
}

// Process simulates task processing
func (t *Task) Process() {
	fmt.Printf("Processing Task ID: %d\n", t.ID)
	// Simulate processing time
	time.Sleep(time.Second)
}

// Worker simulates a worker processing tasks
func Worker(tasks <-chan Task, workerID int, wg *sync.WaitGroup) {
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", workerID, task.ID)
		task.Process()
		wg.Done()
	}
}

func main() {

	time.Sleep(40 * time.Second)

	setupHeartbeatMonitoring(readerHeartbeatFile, serverHeartbeatFile, "Server", serverProgramPath)

	// Setup for task queue
	tasks := make(chan Task, 10)
	var wg sync.WaitGroup

	// Simulate adding tasks to the queue
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		tasks <- Task{ID: i}
	}

	// Start workers
	for i := 1; i <= 2; i++ {
		go Worker(tasks, i, &wg)
	}

	// Wait for all tasks to be processed
	wg.Wait()
	close(tasks)
	fmt.Println("All tasks processed, monitoring continues.")
}

func setupHeartbeatMonitoring(heartbeatFile, partnerFile, partnerName, partnerProgramPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	// Initial heartbeat update to ensure file exists
	updateHeartbeat(heartbeatFile)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateHeartbeat(heartbeatFile)
			if !checkHeartbeat(partnerFile) {
				fmt.Printf("%s is not responding, attempting to restart...\n", partnerName)
				restartProgram(partnerProgramPath)
			}
		}
	}
}

func updateHeartbeat(filename string) {
	ioutil.WriteFile(filename, []byte(fmt.Sprint(time.Now().Unix())), 0644)
}

func checkHeartbeat(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}

	lastHeartbeat, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil {
		return false
	}

	return time.Now().Unix()-lastHeartbeat < int64(heartbeatInterval.Seconds()*2)
}

func restartProgram(programPath string) {
	cmd := exec.Command(programPath)
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start %s: %v", programPath, err)
	} else {
		log.Printf("%s restarted successfully.", programPath)
	}
	time.Sleep(restartDelay)
}

func setupSignalHandler(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Received shutdown signal.")
		cancelFunc()
	}()
}
