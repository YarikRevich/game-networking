package workers

import (
	"bytes"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/YarikRevich/game-networking/client/internal/request"
	"github.com/YarikRevich/game-networking/client/internal/states"
)

type WorkerManager struct {
	state *states.State

	count int // Count of workers

	conn *net.UDPConn

	send    chan []byte
	receive chan []byte
	ping chan []byte

	exit chan os.Signal
	err  chan error
}

func (wm *WorkerManager) Run() {
	signal.Notify(wm.exit, os.Interrupt)

	for i := 0; i <= wm.count; i++ {
		go wm.worker()
	}
	go wm.pingWorker()
}

func (wm *WorkerManager) worker() {
loop:
	for {
		select {
		case <-wm.exit:
			break loop
		default:
			switch wm.state.GetCurrState() {
			case states.RECEIVE:
				var buffer bytes.Buffer
				if _, err := io.Copy(&buffer, wm.conn); err != nil {
					wm.err <- err
				}
				wm.receive <- buffer.Bytes()
			case states.SEND:
				wm.conn.Write(<-wm.send)
			}
		}
	}
}

func (wm *WorkerManager) pingWorker() {
loop:
	for {
		select {
		case <-wm.exit:
			break loop
		default:
			switch wm.state.GetCurrState(){
			case states.PING:

			}
		}
	}
}

func (wm *WorkerManager) Error() error {
	return <-wm.err
}

func (wm *WorkerManager) Read() interface{} {
	wm.state.SetCurrState(states.RECEIVE)
	return <-wm.receive
}

func (wm *WorkerManager) Write(src []byte) {
	wm.state.SetCurrState(states.SEND)
	wm.send <- src
}

func (wm *WorkerManager) Ping()error{
	wm.ping <- request.New().CreateReq(nil, true)
	return nil
}

func New(count int, conn *net.UDPConn) *WorkerManager {
	return &WorkerManager{
		state:   states.New(),
		count:   count,
		conn:    conn,
		send:    make(chan []byte, count),
		receive: make(chan []byte, count),
		ping: make(chan []byte, count),
		exit:    make(chan os.Signal),
		err:     make(chan error, count),
	}
}
