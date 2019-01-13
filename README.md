[![](https://godoc.org/github.com/romainmenke/lifo?status.svg)](http://godoc.org/github.com/romainmenke/lifo)


# LIFO

- Push as many items from multiple go routines as you like.
- Pop blocks until there are items, also save for concurrency.
- Pop takes a context for cancelation

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/romainmenke/lifo"
)

func ExampleStack() {
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
```
