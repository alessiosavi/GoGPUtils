package stack

import (
	"errors"
	"fmt"
	"sync"
)

type Stack struct {
	stack []float64
	m     sync.RWMutex
}

func (s *Stack) Push(f float64) {
	s.m.Lock()
	s.stack = append(s.stack, f)
	s.m.Unlock()
}

func (s *Stack) Pop() (float64, error) {
	if len(s.stack) > 0 {
		s.m.Lock()
		val := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		s.m.Unlock()
		return val, nil

	}
	return -1, errors.New("stack empty")
}

func (s *Stack) Stack() string {
	s.m.RLock()
	str := fmt.Sprint(s.stack)
	s.m.RUnlock()
	return str
}
