package main

import (
	"fmt"
)

func main() {
	//fmt.Println("hello world")
	//var a = 10
	//b := "hello"
	//newString := fmt.Sprintf("%d + %s", a, b)
	//fmt.Println(newString)

	var a = "Runoob"
	fmt.Println(a)

	var b, c = 1, 2
	fmt.Println(b, c)
	var n [10]int /* n 是一个长度为 10 的数组 */
	var i, j int

	/* 为数组 n 初始化元素 */
	for i = 0; i < 10; i++ {
		n[i] = i + 100 /* 设置元素为 i + 100 */
	}

	/* 输出每个数组元素的值 */
	for j = 0; j < 10; j++ {
		fmt.Printf("Element[%d] = %d\n", j, n[j])
	}
	TestForRangeDelayBidding()

}
func TestForRangeDelayBidding() {
	list := []int{1, 2, 3, 4, 5}
	var funcList []func()
	for _, v := range list {
		f := func() {
			fmt.Println(v)
		}
		funcList = append(funcList, f)
	}
	for _, f := range funcList {
		f()
	}

	Delay()
}

func ReturnClosure(name string) func() string {
	return func() string {
		return "Hello, " + name
	}
}

func Delay() {
	fns := make([]func(), 0, 10)

	for i := 0; i < 10; i++ {
		fns = append(fns, func() {
			fmt.Printf("hello, this is : %d \n", i)
		})
	}

	for _, fn := range fns {
		fn()
	}
}
