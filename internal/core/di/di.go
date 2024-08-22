package di

import (
	"log"

	scheduler "github.com/hokamsingh/lessgo/internal/core/job"
	"go.uber.org/dig"
)

type Container struct {
	container *dig.Container
}

func NewContainer() *Container {
	return &Container{
		container: dig.New(),
	}
}

func (c *Container) Register(constructor interface{}) error {
	return c.container.Provide(constructor)
}

// // Provide is an alias for Register. It registers a constructor or provider in the container
// func (c *Container) Provide(constructor interface{}) error {
// 	return c.Register(constructor)
// }

func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}

/*
RegisterScheduler sets up and registers the scheduler in the DI container.

This method ensures that the scheduler is available for dependency injection within your LessGo application. It uses the `cron` package under the hood to provide scheduling capabilities.

### Example Usage

```go
package main

import (

	"log"

	"github.com/hokamsingh/lessgo/pkg/lessgo"
	"github.com/hokamsingh/lessgo/pkg/lessgo/scheduler"

)

	func main() {
	    // Create a new DI container
	    container := lessgo.NewContainer()

	    // Register the scheduler in the container
	    if err := container.RegisterScheduler(); err != nil {
	        log.Fatalf("Error registering scheduler: %v", err)
	    }

	    // Use the scheduler
	    err := container.InvokeScheduler(func(sched scheduler.Scheduler) error {
	        // Add a job to the scheduler
	        if err := sched.AddJob("@every 1m", func() {
	            log.Println("Job running every minute")
	        }); err != nil {
	            return err
	        }

	        // Start the scheduler
	        sched.Start()

	        // Optionally, stop the scheduler when your application shuts down
	        defer sched.Stop()
	        return nil
	    })
	    if err != nil {
	        log.Fatalf("Error invoking scheduler: %v", err)
	    }

	    // Start your application logic here
	}
*/
func (c *Container) RegisterScheduler() error {
	sched := scheduler.NewCronScheduler()
	return c.Register(func() scheduler.Scheduler {
		return sched
	})
}

// InvokeScheduler provides access to the scheduler for initialization or configuration
func (c *Container) InvokeScheduler(fn func(scheduler.Scheduler) error) error {
	return c.container.Invoke(func(sched scheduler.Scheduler) {
		if err := fn(sched); err != nil {
			log.Fatalf("Error invoking scheduler: %v", err)
		}
	})
}
