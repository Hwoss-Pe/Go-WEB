package v1

import (
	"fmt"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter func(c *Context)

//这里写builder的时候，由于上面是一个方法，我这里也是一样参数和返回值就像等于实现了接口

var _ FilterBuilder = MetricFilterBuilder

func MetricFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		start := time.Now().Nanosecond()
		next(c)
		end := time.Now().Nanosecond()
		fmt.Printf("用了 %d 纳秒", end-start)
	}
}
