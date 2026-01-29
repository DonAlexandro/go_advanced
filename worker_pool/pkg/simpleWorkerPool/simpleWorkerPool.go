package simpleWorkerPool

import (
	"fmt"
	"sync"
)

// Task represents a unit of work to be processed by the worker pool.
type Task struct {
	ID    int
	Value int
}

// Worker represents a worker that processes tasks from the task queue.
type Worker struct {
	ID        int
	TaskQueue chan Task
	wg        *sync.WaitGroup
}

// NewWorker creates a new worker with the given ID and task queue.
func NewWorker(id int, taskQueue chan Task, wg *sync.WaitGroup) Worker {
	return Worker{
		ID:        id,
		TaskQueue: taskQueue,
		wg:        wg,
	}
}

// Start initiates the worker to start processing tasks.
func (w Worker) Start() {
	go func() {
		for task := range w.TaskQueue {
			fmt.Printf("Worker %d processing task %d with value %d\n", w.ID, task.ID, task.Value)
			w.wg.Done()
		}
	}()
}

// Dispatcher manages the worker pool and distributes tasks to workers.
type Dispatcher struct {
	WorkerPool chan Worker
	MaxWorkers int
	TaskQueue  chan Task
	wg         sync.WaitGroup
}

// NewDispatcher creates a new dispatcher with the given number of workers.
func NewDispatcher(maxWorkers int, taskQueue chan Task) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan Worker, maxWorkers),
		MaxWorkers: maxWorkers,
		TaskQueue:  taskQueue,
	}
}

// Run initializes the worker pool and starts the workers.
func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i+1, d.TaskQueue, &d.wg)
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
	close(d.TaskQueue)
	d.wg.Wait()
}

func main() {
	taskQueue := make(chan Task, 10)
	dispatcher := NewDispatcher(3, taskQueue)
	dispatcher.Run()

	tasks := []Task{
		{ID: 1, Value: 100},
		{ID: 2, Value: 200},
		{ID: 3, Value: 300},
		{ID: 4, Value: 400},
		{ID: 5, Value: 500},
	}

	dispatcher.Dispatch(tasks)
	fmt.Println("All tasks processed")
}
