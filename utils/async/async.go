// From: https://github.com/wechaty/go-wechaty
// Modified: https://github.com/wdvxdr1123

package async

import (
	"math"
	"runtime"
	"sync"
)

// DefaultWorkerNum default number of goroutines is twice the number of GOMAXPROCS
var DefaultWorkerNum = runtime.GOMAXPROCS(0) * 2

type (
	// Executor interface
	Executor[R any] interface {
		AddTask(task Task[R])
		Result() <-chan Result[R]
	}

	// Result ...
	Result[R any] struct {
		Value R
		Err   error
	}

	async[R any] struct {
		tasks        []Task[R]
		wg           sync.WaitGroup
		maxWorkerNum int
	}

	// Task task func
	Task[R any] func() (R, error)
)

// NewExec ...
func NewExec[R any](maxWorkerNum int) Executor[R] {
	if maxWorkerNum <= 0 {
		maxWorkerNum = DefaultWorkerNum
	}
	return &async[R]{
		maxWorkerNum: maxWorkerNum,
		wg:           sync.WaitGroup{},
	}
}

func (a *async[R]) AddTask(task Task[R]) {
	a.tasks = append(a.tasks, task)
}

func (a *async[R]) Result() <-chan Result[R] {
	taskChan := make(chan Task[R])
	resultChan := make(chan Result[R])
	taskNum := len(a.tasks)
	workerNum := int(math.Min(float64(taskNum), float64(a.maxWorkerNum)))
	a.wg.Add(taskNum)

	for i := 0; i < workerNum; i++ {
		go func() {
			for task := range taskChan {
				result := Result[R]{}
				result.Value, result.Err = task()
				resultChan <- result
				a.wg.Done()
			}
		}()
	}

	go func() {
		defer close(resultChan)
		defer close(taskChan)
		for _, v := range a.tasks {
			taskChan <- v
		}
		a.wg.Wait()
		a.tasks = nil
	}()

	return resultChan
}
