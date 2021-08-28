package states

import (
	"sync"
)

const (
	RECEIVE = iota
	SEND
	PING
)

type State struct {
	sync.RWMutex
	curr int
}

func (s *State) GetCurrState() int {
	s.RLock()
	defer s.RUnlock()
	return s.curr
}

func (s *State) SetCurrState(ns int) {
	s.Lock()
	defer s.Unlock()
	s.curr = ns
}

func New() *State {
	return new(State)
}
