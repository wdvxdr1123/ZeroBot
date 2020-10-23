package main

import (
	"fmt"
	"sync"
)

// 遍历时删除所有的偶数,结果:确实删除了所有的偶数
func fun1() {
	x := sync.Map{}
	// 构建
	for i := 0; i < 100; i++ {
		x.Store(i, i)
	}
	// 遍历时删除偶数
	x.Range(func(k, v interface{}) bool {
		if k.(int)%2 == 0 {
			x.Delete(k)
		}
		return true
	})
	// 遍历打印剩下的
	cout := 0
	x.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		cout++
		return true
	})
	// 会发现是50个,说明删除了所有的偶数
	fmt.Println("删除偶数后,剩余元素数,cout:", cout)
}

// 遍历时删除所有元素,结果:确实删除了所有的元素
func fun2() {
	x := sync.Map{}
	// 构建
	for i := 0; i < 100; i++ {
		x.Store(i, i)
	}
	// 遍历时删除偶数
	x.Range(func(k, v interface{}) bool {
		x.Delete(k)
		return true
	})
	// 遍历打印剩下的
	cout := 0
	x.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		cout++
		return true
	})
	// 会发现是0个,说明删除了所有的元素
	fmt.Println("全部删除后,剩余元素数,cout:", cout)
}
func main() {
	// 遍历时删除一半
	fun1()

	// 遍历时删除所有元素
	fun2()
}
