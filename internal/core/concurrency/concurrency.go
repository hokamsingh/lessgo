package concurrency

import (
	"context"
	"sync"
)

// TaskFunc defines the type for the task function that returns a result and an error.
type TaskFunc func(ctx context.Context) (interface{}, error)

// Task represents an individual task.
type Task struct {
	fn TaskFunc
}

// TaskFunc defines the type for the task function that returns a result and an error.
// @callback TaskFunc
// @param {context.Context} ctx - The context for the task, used for cancellation and deadlines.
// @returns {interface{}, error} The result of the task and an error, if any.
//
// Example:
//
//	taskFunc := func(ctx context.Context) (interface{}, error) {
//	    return "task result", nil
//	}
func NewTask(fn TaskFunc) *Task {
	return &Task{fn: fn}
}

// Execute runs the task function and returns the result or an error.
// Example:
//
//	 result, err := task.Execute(context.Background())
//		if err != nil {
//		    log.Fatal(err)
//		}
//
// fmt.Println(result)
func (t *Task) Execute(ctx context.Context) (interface{}, error) {
	return t.fn(ctx)
}

// ExecutionMode defines whether tasks should run in parallel or sequentially.
type ExecutionMode int

const (
	Parallel   ExecutionMode = 0
	Sequential ExecutionMode = 1
)

// TaskManager manages and executes tasks concurrently or sequentially.
type TaskManager struct {
	tasks []*Task
	mode  ExecutionMode
}

// NewTaskManager creates a new TaskManager with the specified execution mode.
func NewTaskManager(mode ExecutionMode) *TaskManager {
	return &TaskManager{mode: mode}
}

// AddTask adds a task to the manager.
func (tm *TaskManager) AddTask(task *Task) {
	tm.tasks = append(tm.tasks, task)
}

// Run executes tasks based on the specified execution mode.
func (tm *TaskManager) Run(ctx context.Context) ([]interface{}, error) {
	if tm.mode == Parallel {
		return tm.runParallel(ctx)
	}
	return tm.runSequential(ctx)
}

// runParallel executes all tasks concurrently and collects the results.
func (tm *TaskManager) runParallel(ctx context.Context) ([]interface{}, error) {
	var wg sync.WaitGroup
	results := make([]interface{}, len(tm.tasks))
	errChan := make(chan error, len(tm.tasks))

	for i, task := range tm.tasks {
		wg.Add(1)
		go func(i int, t *Task) {
			defer wg.Done()
			result, err := t.Execute(ctx)
			if err != nil {
				errChan <- err
				return
			}
			results[i] = result
		}(i, task)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}
	return results, nil
}

// runSequential executes all tasks one by one and collects the results.
func (tm *TaskManager) runSequential(ctx context.Context) ([]interface{}, error) {
	results := make([]interface{}, len(tm.tasks))

	for i, task := range tm.tasks {
		result, err := task.Execute(ctx)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// TaskBuilder allows building and executing tasks in a chainable manner.
type TaskBuilder struct {
	tm *TaskManager
}

// NewTaskBuilder creates a new TaskBuilder with the specified execution mode.
func NewTaskBuilder(mode ExecutionMode) *TaskBuilder {
	return &TaskBuilder{tm: NewTaskManager(mode)}
}

// Add adds a new TaskFunc to the builder.
func (tb *TaskBuilder) Add(fn TaskFunc) *TaskBuilder {
	tb.tm.AddTask(NewTask(fn))
	return tb
}

// Run executes all tasks and returns the results or an error.
func (tb *TaskBuilder) Run(ctx context.Context) ([]interface{}, error) {
	return tb.tm.Run(ctx)
}
