package zero

import (
	"slices"
	"strconv"
	"testing"
)

// https://github.com/wdvxdr1123/ZeroBot/issues/82
func Test_sortMatcher(t *testing.T) {
	block := 10
	batch := 10
	curBatch := 0
	for i := 0; i < block*batch; i++ {
		c := i
		if c%batch == 0 {
			// 优先级从1开始，因为优先级为0的matcher会先注册先执行
			curBatch++
		}

		OnMessage().Handle(func(ctx *Ctx) {
			ctx.message = strconv.Itoa(c)
		}).SetPriority(batch)
	}

	ctx := &Ctx{}
	var result []int
	for _, m := range matcherList {
		m.Handler(ctx)
		number, err := strconv.Atoi(ctx.message)
		if err != nil {
			// should not happen
			t.Fatal(err)
		}
		result = append(result, number)
	}
	// 每个batch的matcher执行结果应该是有序的
	for i := 0; i < block*batch; i += block {
		batchRes := result[i : i+block]
		// 优先级从1开始的matcher先注册后执行，所以结果是逆序的
		slices.Reverse(batchRes)
		if !slices.IsSorted(batchRes) {
			t.Fatalf("matcherList is not sorted, sort func is not stable: %v", batchRes)
		}
	}

}
