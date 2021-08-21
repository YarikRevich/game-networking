package workers

import (
	"net"

	"github.com/YarikRevich/game-networking/server/internal/table"
	"github.com/YarikRevich/game-networking/server/tools/buffer"
)

type WorkerManager struct {
	count int //Count of workers

	tab *table.Table

	buff *buffer.Buffer
	conn *net.UDPConn
	
	err chan error
}

func (wm *WorkerManager) Run(){
	for i := 0; i < wm.count; i++{
		go wm.worker()
	}
}

func (wm *WorkerManager) worker(){
	for {
		buff, ok := wm.buff.Get().([]byte)
		if !ok{
			continue
		}

		_, addr, err := wm.conn.ReadFrom(buff)
		if err != nil{
			wm.err <- err
		}
		
		wm.tab.Add(addr.String(), buff)

		wm.buff.Put(buff[:0])
	}
}

func (wm *WorkerManager) Error() error {
	return <-wm.err
}


func New(count int) *WorkerManager{
	return &WorkerManager{
		count: count,
		tab: table.New(),
		buff: buffer.New(),
		err: make(chan error, count),
	}
}