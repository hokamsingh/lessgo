package concurrency_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hokamsingh/lessgo/internal/core/concurrency"
)

type TaskFunc = concurrency.TaskFunc

var Parallel = concurrency.Parallel
var Sequential = concurrency.Sequential
var NewTaskBuilder = concurrency.NewTaskBuilder

// Helper function to create a simple task that returns a result after a delay.
func createDelayedTask(result interface{}, delay time.Duration) TaskFunc {
	return func(ctx context.Context) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			return result, nil
		}
	}
}

// Helper function to create a task that returns an error.
func createErrorTask(err error) TaskFunc {
	return func(ctx context.Context) (interface{}, error) {
		return nil, err
	}
}

// Test TaskManager for parallel task execution with no errors.
func TestTaskManager_RunParallel_NoError(t *testing.T) {
	ctx := context.Background()

	taskBuilder := NewTaskBuilder(Parallel, 5)
	taskBuilder.Add(createDelayedTask("result1", 100*time.Millisecond)).
		Add(createDelayedTask("result2", 200*time.Millisecond)).
		Add(createDelayedTask("result3", 300*time.Millisecond))

	results, err := taskBuilder.Run(ctx)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	expectedResults := []interface{}{"result1", "result2", "result3"}
	for i, result := range results {
		if result != expectedResults[i] {
			t.Errorf("Expected result %v, but got %v", expectedResults[i], result)
		}
	}
}

// Test TaskManager for parallel task execution with an error in one task.
func TestTaskManager_RunParallel_WithError(t *testing.T) {
	ctx := context.Background()

	taskBuilder := NewTaskBuilder(Parallel, 5)
	taskBuilder.Add(createDelayedTask("result1", 100*time.Millisecond)).
		Add(createErrorTask(errors.New("task error"))).
		Add(createDelayedTask("result3", 200*time.Millisecond))

	results, err := taskBuilder.Run(ctx)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}

	if err.Error() != "task error" {
		t.Errorf("Expected 'task error', but got %v", err)
	}

	if results != nil {
		t.Errorf("Expected results to be nil on error, but got %v", results)
	}
}

// Test TaskManager for sequential task execution with no errors.
func TestTaskManager_RunSequential_NoError(t *testing.T) {
	ctx := context.Background()

	taskBuilder := NewTaskBuilder(Sequential, 0)
	taskBuilder.Add(createDelayedTask("result1", 100*time.Millisecond)).
		Add(createDelayedTask("result2", 200*time.Millisecond)).
		Add(createDelayedTask("result3", 300*time.Millisecond))

	results, err := taskBuilder.Run(ctx)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	expectedResults := []interface{}{"result1", "result2", "result3"}
	for i, result := range results {
		if result != expectedResults[i] {
			t.Errorf("Expected result %v, but got %v", expectedResults[i], result)
		}
	}
}

// Test TaskManager for sequential task execution with an error in one task.
func TestTaskManager_RunSequential_WithError(t *testing.T) {
	ctx := context.Background()

	taskBuilder := NewTaskBuilder(Sequential, 0)
	taskBuilder.Add(createDelayedTask("result1", 100*time.Millisecond)).
		Add(createErrorTask(errors.New("task error"))).
		Add(createDelayedTask("result3", 200*time.Millisecond))

	results, err := taskBuilder.Run(ctx)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}

	if err.Error() != "task error" {
		t.Errorf("Expected 'task error', but got %v", err)
	}

	if results != nil {
		t.Errorf("Expected results to be nil on error, but got %v", results)
	}
}

// Test for context cancellation in parallel execution.
func TestTaskManager_RunParallel_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	taskBuilder := NewTaskBuilder(Parallel, 5)
	taskBuilder.Add(createDelayedTask("result1", 100*time.Millisecond)).
		Add(createDelayedTask("result2", 200*time.Millisecond)).
		Add(createDelayedTask("result3", 300*time.Millisecond))

	results, err := taskBuilder.Run(ctx)
	if err == nil {
		t.Fatal("Expected context canceled error, but got none")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, but got %v", err)
	}

	if results != nil {
		t.Errorf("Expected results to be nil on cancellation, but got %v", results)
	}
}

// Test for context cancellation in sequential execution.
func TestTaskManager_RunSequential_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	taskBuilder := NewTaskBuilder(Sequential, 0)
	taskBuilder.Add(createDelayedTask("result1", 100*time.Millisecond)).
		Add(createDelayedTask("result2", 200*time.Millisecond)).
		Add(createDelayedTask("result3", 300*time.Millisecond))

	results, err := taskBuilder.Run(ctx)
	if err == nil {
		t.Fatal("Expected context canceled error, but got none")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, but got %v", err)
	}

	if results != nil {
		t.Errorf("Expected results to be nil on cancellation, but got %v", results)
	}
}
