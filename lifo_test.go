package lifo

import (
	"context"
	"fmt"
	"testing"
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
				sema:  make(chan struct{}, tc),
				stack: make([]interface{}, tc),
			}

			for i := 0; i < (tc * 4); i++ {
				s.Push(fmt.Sprint(i))
			}

			if len(s.stack) != tc {
				t.Fatal("wrong stack size for :", tc, "got :", len(s.stack))
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
