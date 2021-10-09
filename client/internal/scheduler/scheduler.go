package scheduler

import (
	"sync"
	"time"
)

type scheduler struct {
	sync.Mutex
	scheduled     int
	confirmations int
	limit         int
}

type IScheduler interface {
	Schedule(func())
	CountScheduled() int
	CountConfirmations() int
	IncConfirmations()
	DecConfirmations()
}

func (s *scheduler) Schedule(c func()) {
	s.Lock()
	s.scheduled++
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

func NewScheduler(limit int) IScheduler {
	return &scheduler{limit: limit}
}
