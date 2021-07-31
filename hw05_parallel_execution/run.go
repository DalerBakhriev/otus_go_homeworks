package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, len(tasks))
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	if n > len(tasks) {
		n = len(tasks)
	}
	var errorsNum int
	for workerNum := 0; workerNum < n; workerNum++ {
		wg.Add(1)
		go func(numErrors *int) {
			defer wg.Done()
			for task := range tasksCh {
				err := task()
				mu.Lock()
				if *numErrors == m && m > 0 {
					mu.Unlock()
					return
				}
				if err != nil {
					*numErrors++
				}
				mu.Unlock()
			}
		}(&errorsNum)
	}

	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)
	wg.Wait()

	if errorsNum == m && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
