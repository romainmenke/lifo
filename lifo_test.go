package lifo

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestStackSize(t *testing.T) {
	tcs := []int{
		1,
		5,
		10,
		100,
		1000,
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s := &stack{
				size:  tc,
				scond: sync.NewCond(&sync.Mutex{}),
				stack: list.New(),
			}

			for i := 0; i < (tc * 4); i++ {
				s.Push(fmt.Sprint(i))
			}

			if s.stack.Len() != tc {
				t.Fatal("wrong stack size for :", tc, "got :", s.stack.Len())
			}
		})
	}

}

func TestNew(t *testing.T) {
	tcs := []int{
		-1,
		0,
		1,
		10,
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s := New(tc)

			for i := 0; i < 10; i++ {
				s.Push(fmt.Sprint(i))
			}

			for i := 0; i < 1; i++ {
				x, err := s.Pop(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				if x == nil {
					t.Fatal("expected non nil popped value")
				}
			}
		})
	}

}

func TestDrain(t *testing.T) {
	tcs := []int{
		1,
		5,
		10,
		100,
		1000,
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s := New(tc)

			for i := 0; i < tc; i++ {
				s.Push(fmt.Sprint(i))
			}

			for i := 0; i < (tc + 20); i++ {
				func() {
					ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
					defer cancel()
					s.Pop(ctx)
				}()
			}

			if s.stack.Len() != 0 {
				t.Fatal("expected stack to be empty, got : ", s.stack.Len())
			}
		})
	}

}

func TestLargeNumberOfCallers(t *testing.T) {
	tcs := []struct {
		pushers int
		poppers int
		size    int
	}{}

	for i := 1; i < 6; i++ {
		x := i * i

		for j := 1; j < 6; j++ {
			y := j * j

			for k := 1; k < 6; k++ {
				z := k * k

				tcs = append(tcs, struct {
					pushers int
					poppers int
					size    int
				}{x, y, z})
			}
		}
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("pushers_%d__poppers_%d__size_%d", tc.pushers, tc.poppers, tc.size), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
			defer cancel()

			s := New(tc.size)

			wg := &sync.WaitGroup{}

			for i := 0; i < tc.pushers; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					for {
						select {
						case <-ctx.Done():
							return
						default:
							//
						}

						s.Push(struct{}{})
					}
				}()
			}

			for i := 0; i < tc.poppers; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					for {
						select {
						case <-ctx.Done():
							return
						default:
							//
						}

						ctxPop, cancelPop := context.WithTimeout(ctx, time.Millisecond*5)
						defer cancelPop()

						s.Pop(ctxPop)
					}
				}()
			}

			<-ctx.Done()
			wg.Wait()

			{
				for i := 0; i < tc.size; i++ {
					s.Push(struct{}{})
				}

				for i := 0; i < tc.size; i++ {
					func() {
						ctxPop, cancelPop := context.WithTimeout(context.Background(), time.Millisecond*10)
						defer cancelPop()

						_, err := s.Pop(ctxPop)
						if err != nil {
							t.Log("stack len : ", s.stack.Len())
							t.Fatal(err)
						}
					}()
				}

			}

			if s.stack.Len() != 0 {
				t.Fatal("expected stack to be empty, got : ", s.stack.Len())
			}
		})
	}

}
