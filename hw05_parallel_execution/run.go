package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, len(tasks))
	wg := &sync.WaitGroup{}
	if n > len(tasks) {
		n = len(tasks)
	}
	var errorsNum int32
	for workerNum := 0; workerNum < n; workerNum++ {
		wg.Add(1)
		go func(numErrors *int32) {
			defer wg.Done()
			for task := range tasksCh {
				err := task()
				if atomic.LoadInt32(numErrors) == int32(m) && m > 0 {
					return
				}
				if err != nil {
					atomic.AddInt32(numErrors, 1)
				}
			}
		}(&errorsNum)
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorsNum) == int32(m) && m > 0 {
			break
		}
		tasksCh <- task
	}
	close(tasksCh)
	wg.Wait()

	if errorsNum == int32(m) && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
