package buffer

import (
	"sync"
)

type Buffer struct {
	*sync.Pool
}

func (b *Buffer) PutToBuffer(e interface{}) {
	b.Put(e)
}

func (b *Buffer) GetFromBuffer() interface{} {
	return b.Get()
}

func New() *Buffer {
	return &Buffer{
		&sync.Pool{
			New: func() interface{} { return make([]byte, 32 * 1024) },
		},
	}
}
