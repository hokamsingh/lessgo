package concurrency

import (
	"context"
	"errors"
	"sync"
)

// TaskFunc defines the type for the task function that returns a result and an error.
type TaskFunc func(ctx context.Context) (interface{}, error)

// Task represents an individual task.
type Task struct {
	fn TaskFunc
}

// NewTask creates a new Task.
func NewTask(fn TaskFunc) *Task {
	return &Task{fn: fn}
}

// Execute runs the task function and returns the result or an error.
func (t *Task) Execute(ctx context.Context) (interface{}, error) {
	return t.fn(ctx)
}

// ExecutionMode defines whether tasks should run in parallel or sequentially.
type ExecutionMode int

const (
	Parallel   ExecutionMode = 0
	Sequential ExecutionMode = 1
)

// WorkerPool manages a fixed number of workers to process tasks concurrently.
type WorkerPool struct {
	taskChan    chan *Task
	resultChan  chan result
	workerCount int
	wg          sync.WaitGroup
}

type result struct {
	index  int
	output interface{}
	err    error
}

// NewWorkerPool initializes a worker pool with the specified number of workers.
func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		taskChan:    make(chan *Task),
		resultChan:  make(chan result),
		workerCount: workerCount,
	}
}

// Run starts the workers in the pool.
func (wp *WorkerPool) Run(ctx context.Context) {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for task := range wp.taskChan {
				select {
				case <-ctx.Done():
					return
				default:
					output, err := task.Execute(ctx)
					wp.resultChan <- result{output: output, err: err}
				}
			}
		}()
	}
}

// Stop waits for all workers to finish.
func (wp *WorkerPool) Stop() {
	close(wp.taskChan)
	wp.wg.Wait()
	close(wp.resultChan)
}

// Submit adds a task to the task channel.
func (wp *WorkerPool) Submit(task *Task, index int) {
	wp.taskChan <- task
}

// Results returns the result channel to collect task outputs and errors.
func (wp *WorkerPool) Results() <-chan result {
	return wp.resultChan
}

// TaskManager manages and executes tasks concurrently or sequentially.
type TaskManager struct {
	tasks       []*Task
	mode        ExecutionMode
	workerCount int
}

// NewTaskManager creates a new TaskManager with the specified execution mode and optional worker count.
func NewTaskManager(mode ExecutionMode, workerCount int) *TaskManager {
	if workerCount <= 0 {
		workerCount = 10
	}
	return &TaskManager{mode: mode, workerCount: workerCount}
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

// runParallel executes all tasks concurrently using a worker pool and collects the results.
func (tm *TaskManager) runParallel(ctx context.Context) ([]interface{}, error) {
	pool := NewWorkerPool(tm.workerCount)
	results := make([]interface{}, len(tm.tasks))
	errChan := make(chan error, 1) // Buffer size 1 for first error

	// Start worker pool
	pool.Run(ctx)

	// Submit tasks to the worker pool
	for i, task := range tm.tasks {
		go func(index int, task *Task) {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
			default:
				pool.Submit(task, index)
			}
		}(i, task)
	}

	// Collect results
	go func() {
		for res := range pool.Results() {
			if res.err != nil {
				select {
				case errChan <- res.err: // pass only first error
				default:
				}
			} else {
				results[res.index] = res.output
			}
		}
	}()

	// Stop the worker pool and wait for results
	pool.Stop()
	close(errChan)

	// Check for errors
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
func NewTaskBuilder(mode ExecutionMode, workerCount int) *TaskBuilder {
	return &TaskBuilder{tm: NewTaskManager(mode, workerCount)}
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
