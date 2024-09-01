
---

## Package `scheduler`

The `scheduler` package provides an interface and an implementation for scheduling recurring jobs using a cron-like syntax.

### Types

#### `Scheduler`

The `Scheduler` interface defines methods for scheduling jobs.

```go
type Scheduler interface {
    AddJob(schedule string, job func()) error
    Start()
    Stop()
}
```

#### `CronScheduler`

The `CronScheduler` struct implements the `Scheduler` interface using the `robfig/cron/v3` package for cron-based job scheduling.

```go
type CronScheduler struct {
    cron *cron.Cron
}
```

### Functions

#### `NewCronScheduler`

```go
func NewCronScheduler() *CronScheduler
```

Creates a new instance of `CronScheduler`. This scheduler can be used to schedule jobs with cron expressions.

**Example:**

```go
s := scheduler.NewCronScheduler()
err := s.AddJob("* * * * *", func() {
    fmt.Println("Hello, World!")
})
if err != nil {
    log.Fatalf("Failed to add job: %v", err)
}
s.Start()
time.Sleep(10 * time.Minute)
s.Stop()
```

#### `AddJob`

```go
func (s *CronScheduler) AddJob(schedule string, job func()) error
```

Adds a job to the scheduler with a specified cron schedule. The job will be executed according to the provided cron expression.

**Parameters:**

- `schedule`: The cron expression to determine the job schedule.
- `job`: The function to execute.

**Example:**

```go
err := s.AddJob("0 0 * * *", func() {
    fmt.Println("It's midnight!")
})
if err != nil {
    log.Fatalf("Failed to add job: %v", err)
}
```

#### `Start`

```go
func (s *CronScheduler) Start()
```

Begins the execution of scheduled jobs. This method should be called to start the scheduler after all jobs have been added.

**Example:**

```go
s.Start()
```

#### `Stop`

```go
func (s *CronScheduler) Stop()
```

Halts the execution of scheduled jobs. This method stops the scheduler and prevents any further job executions.

**Example:**

```go
s.Start()
time.Sleep(1 * time.Hour)
s.Stop()
```

---