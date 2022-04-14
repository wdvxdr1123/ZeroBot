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
	// IAsync interface
	IAsync[Result any] interface {
		AddTask(task Task[Result])
		Result() <-chan AsyncResult[Result]
	}

	// AsyncResult result struct
	AsyncResult[Result any] struct {
		Value Result
		Err   error
	}

	async[Result any] struct {
		tasks        []Task[Result]
		wg           sync.WaitGroup
		maxWorkerNum int
	}

	// Task task func
	Task[Result any] func() (Result, error)
)

// NewAsync ...
func NewAsync[Result any](maxWorkerNum int) IAsync[Result] {
	if maxWorkerNum <= 0 {
		maxWorkerNum = DefaultWorkerNum
	}
	return &async[Result]{
		maxWorkerNum: maxWorkerNum,
		wg:           sync.WaitGroup{},
	}
}

func (a *async[Result]) AddTask(task Task[Result]) {
	a.tasks = append(a.tasks, task)
}

func (a *async[Result]) Result() <-chan AsyncResult[Result] {
	taskChan := make(chan Task[Result])
	resultChan := make(chan AsyncResult[Result])
	taskNum := len(a.tasks)
	workerNum := int(math.Min(float64(taskNum), float64(a.maxWorkerNum)))
	a.wg.Add(taskNum)

	for i := 0; i < workerNum; i++ {
		go func() {
			for task := range taskChan {
				result := AsyncResult[Result]{}
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
