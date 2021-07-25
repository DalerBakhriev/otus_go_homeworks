package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func singleTaskWorker(workerNum int, taskInput <-chan Task, taskOutput chan<- error) {
	for task := range taskInput {
		err := task()
		taskOutput <- err
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorsCh := make(chan error, len(tasks))
	tasksCh := make(chan Task, len(tasks))
	wg := &sync.WaitGroup{}

	for workerNum := 1; workerNum < n+1; workerNum++ {
		wg.Add(1)
		go func(workerNumber int) {
			defer wg.Done()
			singleTaskWorker(workerNumber, tasksCh, errorsCh)
		}(workerNum)
	}

	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	numErrors := 0
	numFinishedTasks := 0
	for numFinishedTasks < len(tasks) {
		taskErr := <-errorsCh
		if taskErr != nil {
			numErrors++
		}
		if numErrors == m && m > 0 {
			return ErrErrorsLimitExceeded
		}
		numFinishedTasks++
	}

	wg.Wait()
	close(errorsCh)

	return nil
}
