package id

import (
	"sync"
)

type LocalRequestID struct {
	sync.Mutex
	id int
}

func (lri *LocalRequestID) SetID(d int) {
	lri.Lock()
	lri.id = d
	lri.Unlock()
}

func (lri *LocalRequestID) GetID() int {
	lri.Lock()
	defer lri.Unlock()
	return lri.id
}

func New() *LocalRequestID {
	return new(LocalRequestID)
}
