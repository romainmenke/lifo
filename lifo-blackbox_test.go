package lifo_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/romainmenke/lifo"
)

func TestLifo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stack := lifo.New(50)

	for i := 1; i < 101; i++ {
		stack.Push(fmt.Sprint(i))
	}

	for i := 100; i > 50; i-- {
		item, err := stack.Pop(ctx)
		if err != nil {
			t.Fatal(err, i)
		}

		str, ok := item.(string)
		if !ok {
			t.Fatal("expected string type")
		}

		if str != fmt.Sprint(i) {
			t.Fatal("expected : " + fmt.Sprint(i) + " got : " + str)
		}
	}
}

func TestLifo_Concurrency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stack := lifo.New(50)

	wg := &sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			stack.Push(fmt.Sprint(i))
			time.Sleep(time.Millisecond * 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			stack.Push(fmt.Sprint(i))
			time.Sleep(time.Millisecond * 5)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_, err := stack.Pop(ctx)
			if err != nil {
				panic(err)
			}

		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_, err := stack.Pop(ctx)
			if err != nil {
				panic(err)
			}

		}
	}()

	wg.Wait()
}

func TestLifo_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	stack := lifo.New(0)

	_, err := stack.Pop(ctx)
	if err != ctx.Err() {
		panic(err)
	}
}

func benchmarkLifoPush(size int, b *testing.B) {
	stack := lifo.New(size)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		stack.Push(struct{}{})
	}
}

func BenchmarkLifoPush1(b *testing.B)   { benchmarkLifoPush(1, b) }
func BenchmarkLifoPush10(b *testing.B)  { benchmarkLifoPush(10, b) }
func BenchmarkLifoPush100(b *testing.B) { benchmarkLifoPush(100, b) }
func BenchmarkLifoPush1K(b *testing.B)  { benchmarkLifoPush(1000, b) }
func BenchmarkLifoPush10K(b *testing.B) { benchmarkLifoPush(1000*10, b) }

func benchmarkLifoPop(size int, b *testing.B) {
	stack := lifo.New(size)
	for i := 0; i < size; i++ {
		stack.Push(struct{}{})
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		stack.Push(struct{}{})
		_, err := stack.Pop(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLifoPop1(b *testing.B)   { benchmarkLifoPop(1, b) }
func BenchmarkLifoPop10(b *testing.B)  { benchmarkLifoPop(10, b) }
func BenchmarkLifoPop100(b *testing.B) { benchmarkLifoPop(100, b) }
func BenchmarkLifoPop1K(b *testing.B)  { benchmarkLifoPop(1000, b) }
func BenchmarkLifoPop10K(b *testing.B) { benchmarkLifoPop(1000*10, b) }
