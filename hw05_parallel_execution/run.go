package hw05parallelexecution

import (
	"errors"
	"log"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func singleTaskWorker(workerNum int, taskInput <-chan Task, taskOutput chan<- error, done <-chan struct{}) {

	for {
		select {
		case <-done:
			log.Println("Got signal for finishing")
			return
		case task := <-taskInput:
			err := task()
			log.Printf("Sending error %v", err)
			taskOutput <- err
		}

	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorsCh := make(chan error, len(tasks))
	tasksCh := make(chan Task, len(tasks))
	wg := &sync.WaitGroup{}
	finishChannels := make([]chan struct{}, 0, n)
	for workerNum := 1; workerNum < n+1; workerNum++ {
		wg.Add(1)
		finishCh := make(chan struct{}, 1)
		go func() {
			defer wg.Done()
			singleTaskWorker(workerNum, tasksCh, errorsCh, finishCh)
		}()
		finishChannels = append(finishChannels, finishCh)
	}
	log.Printf("Finish channels list is %v", finishChannels)

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
		log.Printf("Errors number is %d", numErrors)
		if numErrors == m && m > 0 {
			log.Printf("Errors number is %d, finishing workers with %d channels", m, len(finishChannels))
			for _, ch := range finishChannels {
				log.Printf("Sending finish signal for channel %v", ch)
				ch <- struct{}{}
			}
			return ErrErrorsLimitExceeded
		}
		numFinishedTasks++
	}

	wg.Wait()
	close(errorsCh)

	return nil
}
