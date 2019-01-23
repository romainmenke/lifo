package lifo

import (
	"container/list"
	"context"
	"errors"
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
		stack: list.New(),
		scond: sync.NewCond(&sync.Mutex{}),
	}

	return s
}

type stack struct {
	scond *sync.Cond

	size  int
	stack *list.List
}

// Add never blocks, it just removes the oldest items from the queue
// Handling items that never get processed is up to the implementer
func (s *stack) Push(v interface{}) {
	s.scond.L.Lock()

	needsSignal := false
	if s.stack.Len() >= s.size {
		e := s.stack.Back()
		if e != nil {
			s.stack.Remove(e)
		}
	} else {
		needsSignal = true
	}

	s.stack.PushFront(v)

	if needsSignal {
		s.scond.Signal()
	}

	s.scond.L.Unlock()
}

// Pop blocks until there is an item to pop
func (s *stack) Pop(ctx context.Context) (interface{}, error) {
	s.scond.L.Lock()
	defer s.scond.L.Unlock()

	if s.stack.Len() > 0 {
		v, err := s.popLocked()
		return v, err
	}

	ctx2, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			s.scond.Broadcast()
			return
		case <-ctx2.Done():
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			//
		}

		if s.stack.Len() > 0 {
			break
		}

		s.scond.Wait()
	}

	v, err := s.popLocked()

	return v, err
}

// Pop blocks until there is an item to pop
func (s *stack) popLocked() (interface{}, error) {
	e := s.stack.Front()
	if e == nil {
		return nil, errors.New("lifo : tried to pop an empty stack")
	}

	s.stack.Remove(e)

	return e.Value, nil
}
