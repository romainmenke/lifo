package lifo

import (
	"context"
	"sync"
)

type Stack interface {
	Push(item interface{})
	Pop(ctx context.Context) (interface{}, error)
}

func New(size int) *stack {
	if size < 1 {
		size = 1
	}

	s := &stack{
		size:  size,
		sema:  make(chan struct{}, size),
		stack: make([]interface{}, 0, size*2),
	}

	return s
}

type stack struct {
	mu sync.Mutex

	size  int
	stack []interface{}
	sema  chan struct{}
}

// Add never blocks, it just removes the oldest items from the queue
// Handling items that never get processed is up to the implementer
func (s *stack) Push(v interface{}) {
	s.mu.Lock()

	select {
	case s.sema <- struct{}{}:
		// always try to add another to the semaphore
	default:
		// semaphore already filled
	}

	if len(s.stack) >= s.size {
		s.stack[0] = nil
		s.stack = s.stack[1:]
	}

	s.stack = append(s.stack, v)

	s.mu.Unlock()
}

// Pop blocks until there is an item to pop
func (s *stack) Pop(ctx context.Context) (interface{}, error) {
	select {
	case <-s.sema:
		// wait for item in stack
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	s.mu.Lock()

	n := len(s.stack)

	v := s.stack[n-1]

	s.stack[n-1] = nil
	s.stack = s.stack[:n-1]

	s.mu.Unlock()

	return v, nil
}
