package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Job holds details about a job.
type Job struct {
	ID          int64
	Type        string
	Description string
	Status      atomic.Value // Manages the job's status atomically.
}

type JobWatcher struct {
	nextJobID int64
	jobs      atomic.Value // Stores map[int64]*AtomicJob
}

// NewJob initializes a new job with specified details.
func NewJob(id int64, jobType, description string) *Job {
	job := &Job{ID: id, Type: jobType, Description: description}
	job.Status.Store("new") // Initializes status to "new".
	return job
}

// AtomicJob provides an atomic reference to a Job.
type AtomicJob struct {
	value atomic.Value
}

func (aj *AtomicJob) Store(job *Job) {
	aj.value.Store(job)
}

func (aj *AtomicJob) Load() *Job {
	if v := aj.value.Load(); v != nil {
		return v.(*Job)
	}
	return nil
}

func NewJobWatcher() *JobWatcher {
	jw := &JobWatcher{}
	jw.jobs.Store(make(map[int64]*AtomicJob))
	return jw
}

func (jw *JobWatcher) AddJob(jobType, description string) {
	jobID := atomic.AddInt64(&jw.nextJobID, 1)
	job := &Job{ID: jobID, Type: jobType, Description: description}
	job.Status.Store("new")
	atomicJob := &AtomicJob{}
	atomicJob.Store(job)

	// Safely update the map
	jobsMap := jw.jobs.Load().(map[int64]*AtomicJob)
	newMap := make(map[int64]*AtomicJob)
	for k, v := range jobsMap {
		newMap[k] = v
	}
	newMap[jobID] = atomicJob
	jw.jobs.Store(newMap)

	fmt.Printf("Added job: ID=%d, Type=%s\n", jobID, jobType)
}

func SetupRoutes(app *fiber.App, watcher *JobWatcher) {
	app.Post("/newjob", func(c *fiber.Ctx) error {
		jobType := c.FormValue("type")
		description := c.FormValue("description")
		watcher.AddJob(jobType, description)
		return c.JSON(fiber.Map{"message": "Job added successfully"})
	})
}

func AutonomousObserver(watcher *JobWatcher, condition string, action func(*Job)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		jobsMap := watcher.jobs.Load().(map[int64]*AtomicJob)
		for _, aj := range jobsMap {
			if job := aj.Load(); job != nil && job.Status.Load() == condition {
				action(job)
			}
		}
	}
}

func main() {
	app := fiber.New()
	watcher := NewJobWatcher()

	SetupRoutes(app, watcher)
	ProcessNewJobs(watcher)
	PerformQualityAssurance(watcher)

	log.Fatal(app.Listen(":8080"))
}

// Example usage of AutonomousObserver for new job processing and QA.
func ProcessNewJobs(watcher *JobWatcher) {
	action := func(job *Job) {
		fmt.Printf("Processing job: ID=%d\n", job.ID)
		// Replace with actual processing logic
		job.Status.Store("review")
	}
	go AutonomousObserver(watcher, "new", action)
}

// PerformQualityAssurance conducts QA on jobs marked for review.
func PerformQualityAssurance(watcher *JobWatcher) {
	qaAction := func(job *Job) {
		fmt.Printf("QA on job: ID=%d\n", job.ID)
		// Replace with actual QA logic
		job.Status.Store("approved")
	}
	go AutonomousObserver(watcher, "review", qaAction)
}

/*

Refined Function Template
Here is a polished version of the function template that adheres to Go's idiomatic practices:


func WatchAndProcessJobs(watcher *JobWatcher, triggerStatus string, process func(job *Job)) {
    processJob := func(job *Job) {
        // Explain what processing will occur here, specific to the triggerStatus.
        // For example, "Processing job for QA checks" for a triggerStatus of "review".
        fmt.Printf("Triggered processing on job: ID=%d, Type=%s, Status=%s\n", job.ID, job.Type, triggerStatus)

        // Simulate some processing time.
        time.Sleep(2 * time.Second)

        // Update job's status after processing.
        // For instance, setting status to "completed" after a "new" job is processed.
        job.Status.Store("completed")
        fmt.Printf("Processed job: ID=%d, new Status=%s\n", job.ID, "completed")
    }

    // Start the AutonomousObserver goroutine to monitor and act upon jobs with the specified triggerStatus.
    go AutonomousObserver(watcher, triggerStatus, processJob)
}
*/
