package client

import (
	"sync"
)

type scheduler struct {
	sync.Mutex
	confirmations int
	capacity         int
	scheduled chan func()
}

type IScheduler interface {
	Schedule(func())
	CountConfirmations() int
	IncConfirmations()
	DecConfirmations()
}

func (s *scheduler) loop(){
	for {
		select {
		case c := <- s.scheduled:
			c()
		}
	}
}

func (s *scheduler) Schedule(c func()) {
	if s.capacity == 0{
		c()
		s.Lock()
		s.capacity--
		s.Unlock()
		return
	}

	s.scheduled <- c

	// s.Lock()
	// s.capacity--
	// s.Unlock()
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	wg.Done()
	// 	for {
	// 		if s.confirmations == 0 {
	// 			c()
	// 			s.Lock()
	// 			s.confirmations--
	// 			s.Unlock()
	// 			break
	// 		}
	// 	}
	// }()
	// wg.Wait()
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
	s := &scheduler{capacity: capacity}
	go s.loop()
	return s
}
