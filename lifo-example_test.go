package lifo_test

import (
	"context"
	"fmt"
	"time"

	"github.com/romainmenke/lifo"
)

func ExampleLifo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	lifoStack := lifo.New(2)

	lifoStack.Push("a")
	lifoStack.Push("b")
	lifoStack.Push("c")

	// Pop blocks until there are items on the stack
	x, err := lifoStack.Pop(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(x)
	// Output: c
}
