package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) < n {
		n = len(tasks)
	}

	doneCh := make(chan struct{})
	taskCh := make(chan Task)
	errorCh := make(chan error, n)

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go runTask(&wg, taskCh, errorCh)
	}
	go countErrors(m, errorCh, doneCh)
	result := genTasks(tasks, doneCh, taskCh)

	close(taskCh)
	wg.Wait()
	close(doneCh)
	close(errorCh)
	return result
}

func runTask(wg *sync.WaitGroup, taskCh <-chan Task, errorCh chan<- error) {
	defer wg.Done()
	for task := range taskCh {
		err := task()
		if err != nil {
			errorCh <- err
		}
	}
}

func countErrors(m int, errorCh <-chan error, doneCh chan<- struct{}) {
	errorsCount := 0
	for range errorCh {
		errorsCount++
		if m > 0 && errorsCount == m {
			doneCh <- struct{}{}
		}
	}
}

func genTasks(tasks []Task, doneCh <-chan struct{}, taskCh chan<- Task) error {
	for _, task := range tasks {
		select {
		case <-doneCh:
			return ErrErrorsLimitExceeded
		case taskCh <- task:
		}
	}
	return nil
}
