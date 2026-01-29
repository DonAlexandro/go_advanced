package advancedworkerpool

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work to be processed by the worker pool.
type Task struct {
	ID       int
	Value    int
	Priority int
	Result   chan<- Result
}

// Result represents the result of a task processed by a worker.
type Result struct {
	TaskID int
	Output int
	Err    error
}

// Worker represents a worker that processes tasks from the task queue.
type Worker struct {
	ID          int
	TaskQueue   chan Task
	ResultQueue chan<- Result
	wg          *sync.WaitGroup
}

// NewWorker creates a new worker with the given ID, task queue, and result queue.
func NewWorker(id int, taskQueue chan Task, resultQueue chan<- Result, wg *sync.WaitGroup) Worker {
	return Worker{
		ID:          id,
		TaskQueue:   taskQueue,
		ResultQueue: resultQueue,
		wg:          wg,
	}
}

// Start initiates the worker to start processing tasks.
func (w Worker) Start() {
	go func() {
		for task := range w.TaskQueue {
			result := w.processTask(task)
			w.ResultQueue <- result
			w.wg.Done()
		}
	}()
}

// processTask processes a single task and returns the result.
func (w Worker) processTask(task Task) Result {
	// Simulate varying processing times based on priority
	time.Sleep(time.Millisecond * time.Duration(task.Priority*10))

	output := task.Value * 2 // Example processing: doubling the value
	fmt.Printf("Worker %d processed task %d with value %d\n", w.ID, task.ID, task.Value)

	return Result{TaskID: task.ID, Output: output, Err: nil}
}

// Dispatcher manages the worker pool and distributes tasks to workers.
type Dispatcher struct {
	WorkerPool  chan Worker
	MaxWorkers  int
	TaskQueue   chan Task
	ResultQueue chan Result
	wg          sync.WaitGroup
}

// NewDispatcher creates a new dispatcher with the given number of workers.
func NewDispatcher(maxWorkers int, taskQueue chan Task, resultQueue chan Result) *Dispatcher {
	return &Dispatcher{
		WorkerPool:  make(chan Worker, maxWorkers),
		MaxWorkers:  maxWorkers,
		TaskQueue:   taskQueue,
		ResultQueue: resultQueue,
	}
}

// Run initializes the worker pool and starts the workers.
func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i+1, d.TaskQueue, d.ResultQueue, &d.wg)
		d.WorkerPool <- worker
		worker.Start()
	}
}

// Dispatch adds tasks to the task queue and increments the wait group counter.
func (d *Dispatcher) Dispatch(tasks []Task) {
	for _, task := range tasks {
		d.wg.Add(1)
		d.TaskQueue <- task
	}

	// Wait for all tasks to be done in a separate goroutine and then close ResultQueue
	go func() {
		d.wg.Wait()
		close(d.ResultQueue)
	}()
}

// CollectResults collects results from the result queue.
func (d *Dispatcher) CollectResults() {
	for result := range d.ResultQueue {
		if result.Err != nil {
			fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Err)
		} else {
			fmt.Printf("Task %d result: %d\n", result.TaskID, result.Output)
		}
	}
}

func main() {
	taskQueue := make(chan Task, 10)
	resultQueue := make(chan Result, 10)

	dispatcher := NewDispatcher(3, taskQueue, resultQueue)
	dispatcher.Run()

	tasks := []Task{
		{ID: 1, Value: 100, Priority: 1, Result: resultQueue},
		{ID: 2, Value: 200, Priority: 2, Result: resultQueue},
		{ID: 3, Value: 300, Priority: 3, Result: resultQueue},
		{ID: 4, Value: 400, Priority: 4, Result: resultQueue},
		{ID: 5, Value: 500, Priority: 5, Result: resultQueue},
	}

	dispatcher.Dispatch(tasks)
	dispatcher.CollectResults()

	fmt.Println("All tasks processed")
}
