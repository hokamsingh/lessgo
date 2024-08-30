/*
Package di provides a simple Dependency Injection (DI) container using the `dig` library.

This package defines a `Container` struct that wraps the `dig.Container` and provides methods for registering and invoking dependencies. It also includes specialized methods for registering and invoking a scheduler within the DI container.

Usage:

	import (
		"log"

		"github.com/hokamsingh/lessgo/pkg/lessgo/di"
		"github.com/hokamsingh/lessgo/pkg/lessgo/scheduler"
	)

	func main() {
		// Create a new DI container
		container := di.NewContainer()

		// Register dependencies
		if err := container.RegisterScheduler(); err != nil {
			log.Fatalf("Error registering scheduler: %v", err)
		}

		// Invoke a function with the scheduler dependency
		err := container.InvokeScheduler(func(sched scheduler.Scheduler) error {
			// Add a job to the scheduler
			if err := sched.AddJob("@every 1m", func() {
				log.Println("Job running every minute")
			}); err != nil {
				return err
			}

			// Start the scheduler
			sched.Start()

			// Stop the scheduler when the application shuts down
			defer sched.Stop()
			return nil
		})
		if err != nil {
			log.Fatalf("Error invoking scheduler: %v", err)
		}

		// Start your application logic here
	}
*/
package di

import (
	"fmt"
	"log"

	"github.com/hokamsingh/lessgo/internal/core/controller"
	scheduler "github.com/hokamsingh/lessgo/internal/core/job"
	"github.com/hokamsingh/lessgo/internal/core/module"
	"github.com/hokamsingh/lessgo/internal/core/router"
	"go.uber.org/dig"
)

// Container wraps the `dig.Container` and provides methods for registering and invoking dependencies.
// This struct serves as the main entry point for setting up and managing dependency injection within the application.
type Container struct {
	container *dig.Container
}

// NewContainer creates a new instance of `Container`.
// This method initializes the underlying `dig.Container`.
//
// Example:
//
//	container := di.NewContainer()
func NewContainer() *Container {
	return &Container{
		container: dig.New(),
	}
}

// Register adds a constructor or provider to the DI container.
// This method allows you to register dependencies that can later be resolved and injected where needed.
//
// Example:
//
//	container := di.NewContainer()
//	err := container.Register(func() MyService {
//		return NewMyService()
//	})
func (c *Container) Register(constructor interface{}) error {
	return c.container.Provide(constructor)
}

// Invoke resolves dependencies and invokes the specified function.
// This method allows you to execute a function with its dependencies automatically injected by the container.
//
// Example:
//
//	container := di.NewContainer()
//	err := container.Invoke(func(svc MyService) {
//		svc.DoSomething()
//	})
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}

// RegisterScheduler sets up and registers the scheduler in the DI container.
// This method ensures that the scheduler is available for dependency injection within your LessGo application.
//
// Example:
//
//	container := di.NewContainer()
//	err := container.RegisterScheduler()
//	if err != nil {
//		log.Fatalf("Error registering scheduler: %v", err)
//	}
func (c *Container) RegisterScheduler() error {
	sched := scheduler.NewCronScheduler()
	return c.Register(func() scheduler.Scheduler {
		return sched
	})
}

// InvokeScheduler provides access to the scheduler for initialization or configuration.
// This method invokes a function that takes the scheduler as a parameter, allowing you to configure it.
//
// Example:
//
//	container := di.NewContainer()
//	err := container.RegisterScheduler()
//	if err != nil {
//		log.Fatalf("Error registering scheduler: %v", err)
//	}
//
//	err = container.InvokeScheduler(func(sched scheduler.Scheduler) error {
//		// Configure the scheduler
//		return sched.AddJob("@hourly", func() {
//			log.Println("Job running every hour")
//		})
//	})
//
//	if err != nil {
//		log.Fatalf("Error invoking scheduler: %v", err)
//	}
func (c *Container) InvokeScheduler(fn func(scheduler.Scheduler) error) error {
	return c.container.Invoke(func(sched scheduler.Scheduler) {
		if err := fn(sched); err != nil {
			log.Fatalf("Error invoking scheduler: %v", err)
		}
	})
}

// RegisterDependencies registers dependencies into container
func RegisterDependencies(dependencies []interface{}) {
	container := NewContainer()
	for _, dep := range dependencies {
		if err := container.Register(dep); err != nil {
			log.Fatalf("Error registering dependencies: %v", err)
		}
	}
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	SkyBlue = "\033[36m"
)

// RegisterModules iterates over a slice of modules and registers their routes.
func RegisterModules(r *router.Router, modules []module.IModule) error {
	for _, module := range modules {
		controller.RegisterModuleRoutes(r, module)
		l := fmt.Sprintf("%sLessGo :: Registered module %s%s%s", Green, Yellow, module.GetName(), Reset)
		log.Println(l)
	}
	return nil
}
