package demo

import (
	"fmt"
	v1 "test/pkg"
)

func init() {
	v1.RegisterFilter("my-custom", myFilterBuilder)
}

func myFilterBuilder(next v1.Filter) v1.Filter {
	return func(c *v1.Context) {
		fmt.Println("假装这是我自定义的 filter")
		next(c)
	}
}
