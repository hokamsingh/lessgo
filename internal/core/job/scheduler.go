/*
Package scheduler provides an interface and implementation for scheduling recurring jobs using a cron-like syntax.

This package defines a Scheduler interface and provides a CronScheduler implementation that leverages the `robfig/cron/v3` package.

Usage:

	import "your/package/path/scheduler"

	func main() {
		s := scheduler.NewCronScheduler()

		// Add a job to print "Hello, World!" every minute
		err := s.AddJob("* * * * *", func() {
			fmt.Println("Hello, World!")
		})

		if err != nil {
			log.Fatalf("Failed to add job: %v", err)
		}

		// Start the scheduler
		s.Start()

		// Stop the scheduler after some time (e.g., 10 minutes)
		time.Sleep(10 * time.Minute)
		s.Stop()
	}
*/
package scheduler

import (
	"github.com/robfig/cron/v3"
)

// Scheduler defines the interface for scheduling jobs.
// Implementations of this interface should provide methods to add jobs, start the scheduler, and stop it.
type Scheduler interface {
	AddJob(schedule string, job func()) error
	Start()
	Stop()
}

// CronScheduler is an implementation of the Scheduler interface using the `robfig/cron/v3` package.
// This scheduler allows for scheduling jobs based on cron expressions.
type CronScheduler struct {
	cron *cron.Cron
}

// NewCronScheduler creates a new instance of CronScheduler.
// This scheduler can be used to schedule jobs with cron expressions.
//
// Example:
//
//	s := scheduler.NewCronScheduler()
//	err := s.AddJob("* * * * *", func() {
//		fmt.Println("Hello, World!")
//	})
//
//	if err != nil {
//		log.Fatalf("Failed to add job: %v", err)
//	}
//
//	s.Start()
//	time.Sleep(10 * time.Minute)
//	s.Stop()
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
	}
}

// AddJob adds a job to the scheduler with a specified cron schedule.
// The job will be executed according to the cron expression provided.
//
// Example:
//
//	s := scheduler.NewCronScheduler()
//	err := s.AddJob("0 0 * * *", func() {
//		fmt.Println("It's midnight!")
//	})
//
//	if err != nil {
//		log.Fatalf("Failed to add job: %v", err)
//	}
func (s *CronScheduler) AddJob(schedule string, job func()) error {
	_, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return err
	}
	return nil
}

// Start begins the execution of scheduled jobs.
// This method should be called to start the scheduler after all jobs have been added.
//
// Example:
//
//	s := scheduler.NewCronScheduler()
//	_ = s.AddJob("* * * * *", func() {
//		fmt.Println("Running every minute")
//	})
//	s.Start()
func (s *CronScheduler) Start() {
	s.cron.Start()
}

// Stop halts the execution of scheduled jobs.
// This method can be called to stop the scheduler, preventing any further job executions.
//
// Example:
//
//	s := scheduler.NewCronScheduler()
//	s.Start()
//	time.Sleep(1 * time.Hour)
//	s.Stop()
func (s *CronScheduler) Stop() {
	s.cron.Stop()
}
