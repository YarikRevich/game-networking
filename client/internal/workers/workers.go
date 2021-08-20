package workers

import (
	"bytes"
	// "encoding/json"
	"io"
	"net"
	"os"
	"os/signal"
	// "reflect"
)

type WorkerManager struct {
	count int // Count of workers

	conn *net.UDPConn

	exit chan os.Signal
	err  chan error
}

func (wm *WorkerManager) Run() {
	signal.Notify(wm.exit, os.Interrupt)

	for i := 0; i <= wm.count; i++ {
		go wm.worker()
	}
}

func (wm *WorkerManager) worker() {
loop:
	for {
		select {
		case <-wm.exit:
			break loop
		default:
			var buffer bytes.Buffer
			if _, err := io.Copy(&buffer, wm.conn); err != nil {
				wm.err <- err
			}
		}
	}
}

func (wm *WorkerManager) Error() error {
	return <-wm.err
}

//dst is the place where deserialised conf will be saved to
func (wm *WorkerManager) Read(dst interface{}) interface{} {
	return nil
}

func (wm *WorkerManager) Write(src interface{}) {

} 

func New(count int, conn *net.UDPConn) *WorkerManager {
	return &WorkerManager{
		count: count,
		conn: conn,
		exit:  make(chan os.Signal),
		err:   make(chan error, count),
	}
}
