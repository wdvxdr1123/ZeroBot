package async

import (
	"fmt"
	"testing"
)

func TestAsync(t *testing.T) {
	a := NewAsync(0)
	for i := 1; i <= 10; i++ {
		var index = i
		a.AddTask(func() (interface{}, error) {
			return index, nil
		})
	}
	for i := range a.Result() {
		fmt.Println(i.Value)
	}
}
