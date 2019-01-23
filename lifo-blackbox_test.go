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

func BenchmarkLifoPushSize1(b *testing.B)   { benchmarkLifoPush(1, b) }
func BenchmarkLifoPushSize10(b *testing.B)  { benchmarkLifoPush(10, b) }
func BenchmarkLifoPushSize100(b *testing.B) { benchmarkLifoPush(100, b) }
func BenchmarkLifoPushSize1K(b *testing.B)  { benchmarkLifoPush(1000, b) }
func BenchmarkLifoPushSize10K(b *testing.B) { benchmarkLifoPush(1000*10, b) }

func benchmarkLifoPopSize(size int, b *testing.B) {
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

func BenchmarkLifoPopSize1(b *testing.B)   { benchmarkLifoPopSize(1, b) }
func BenchmarkLifoPopSize10(b *testing.B)  { benchmarkLifoPopSize(10, b) }
func BenchmarkLifoPopSize100(b *testing.B) { benchmarkLifoPopSize(100, b) }
func BenchmarkLifoPopSize1K(b *testing.B)  { benchmarkLifoPopSize(1000, b) }
func BenchmarkLifoPopSize10K(b *testing.B) { benchmarkLifoPopSize(1000*10, b) }

// IO is the expected bottleneck
// Ensure that large numbers of poppers and pushers don't cause run away locking or cross signalling between callers.
// 2 x size Pushes
// 1 x size Poppers
// 2 x size Stack
func benchmarkLifoPopCallers(size int, b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stack := lifo.New(size * 2)

	for i := 0; i < (size * 2); i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				stack.Push(struct{}{})
			}
		}()
	}

	for i := 0; i < size; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				func() {
					ctxPop, cancelPop := context.WithTimeout(context.Background(), time.Microsecond*5)
					defer cancelPop()

					stack.Pop(ctxPop)
				}()
			}
		}()
	}

	time.Sleep(time.Second)

	b.ResetTimer()

	// Bench how long it takes to Pop 1 Item
	for n := 0; n < b.N; n++ {
		func() {

			func() {
				ctx2, cancel2 := context.WithTimeout(context.Background(), time.Millisecond*500)
				defer cancel2()

				_, err := stack.Pop(ctx2)
				if err != nil {
					b.Fatal(err)
				}
			}()

		}()
	}
}

func BenchmarkLifoPopCallers3(b *testing.B)   { benchmarkLifoPopCallers(1, b) }
func BenchmarkLifoPopCallers30(b *testing.B)  { benchmarkLifoPopCallers(10, b) }
func BenchmarkLifoPopCallers300(b *testing.B) { benchmarkLifoPopCallers(100, b) }
func BenchmarkLifoPopCallers3K(b *testing.B)  { benchmarkLifoPopCallers(1000, b) }
