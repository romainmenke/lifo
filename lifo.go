package lifo

import (
	"container/list"
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
		stack: list.New(),
	}

	return s
}

type stack struct {
	mu sync.Mutex

	size  int
	stack *list.List
	sema  chan struct{}
}

// Add never blocks, it just removes the oldest items from the queue
// Handling items that never get processed is up to the implementer
func (s *stack) Push(v interface{}) {
	s.mu.Lock()

	if s.stack.Len() >= s.size {
		e := s.stack.Back()
		if e != nil {
			s.stack.Remove(e)
		}
	} else {
		select {
		case s.sema <- struct{}{}:
			// always try to add another to the semaphore
		default:
			//
		}
	}

	s.stack.PushFront(v)

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

	e := s.stack.Front()
	if e != nil {
		s.stack.Remove(e)
	}

	s.mu.Unlock()

	return e.Value, nil
}
