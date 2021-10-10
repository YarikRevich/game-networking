package client

import (
	"sync"
	"time"
)

type scheduler struct {
	sync.Mutex
	scheduled     int
	confirmations int
	capacity         int
}

type IScheduler interface {
	Schedule(func())
	CountScheduled() int
	CountConfirmations() int
	IncConfirmations()
	DecConfirmations()
}

func (s *scheduler) Schedule(c func()) {
	if s.capacity == 0{
		c()
		s.Lock()
		s.capacity--
		s.Unlock()
		return
	}

	s.Lock()
	s.scheduled++
	s.capacity--
	s.Unlock()
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			if s.confirmations == 0 {
				go c()
				s.Lock()
				s.confirmations--
				s.Unlock()
			}
		}
	}()
}

func (s *scheduler) CountScheduled() int {
	return s.scheduled
}

func (s *scheduler) CountConfirmations() int {
	return s.confirmations
}

func (s *scheduler) IncConfirmations() {
	s.Lock()
	s.confirmations++
	s.Unlock()
}

func (s *scheduler) DecConfirmations() {
	s.Lock()
	s.confirmations--
	s.Unlock()
}

func NewScheduler(capacity int) IScheduler {
	return &scheduler{capacity: capacity}
}
