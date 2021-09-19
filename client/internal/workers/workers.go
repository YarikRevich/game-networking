package workers

import (
	"bytes"
	// "encoding/json"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/YarikRevich/game-networking/client/internal/states"
	// "github.com/YarikRevich/game-networking/protocol/internal/id"
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
)

type WorkerManager struct {
	conn  *net.UDPConn
	state *states.State

	workersNum uint32

	send    chan []byte
	receive chan []byte
	ping    chan []byte
	exit    chan os.Signal
	err     chan error
}

func (wm *WorkerManager) Run() {
	signal.Notify(wm.exit, os.Interrupt)

	for i := 0; i <= int(wm.workersNum); i++ {
		go wm.worker()
	}
	// go wm.pingWorker()
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
				if _, err := wm.conn.Write(<-wm.send); err != nil {
					wm.err <- err
				}
			}
		}
	}
}

// func (wm *WorkerManager) pingWorker() {
// loop:
// 	for {
// 		select {
// 		case <-wm.exit:
// 			break loop
// 		default:
// 			switch wm.state.GetCurrState() {
// 			case states.PING:

// 			}
// 		}
// 	}
// }

func (wm *WorkerManager) Error() error {
	return <-wm.err
}

func (wm *WorkerManager) Read() ([]byte, error) {
	wm.state.SetCurrState(states.RECEIVE)

	// var m models.Msg
	// if err := json.Unmarshal(<-wm.receive, &m); err != nil {
	// return models.Msg{}, err
	// }
	return nil, nil
}

func (wm *WorkerManager) Write([]byte) {
	// request.CompleteRequestWithID(wm.lri.GetID(), &src)

	
	wm.state.SetCurrState(states.SEND)
	// wm.send <- request.FormatRequestToJSON(src)
}

func (wm *WorkerManager) Ping() error {
	// wm.ping <- request.New().CreateReq(nil, true)
	return nil
}

func New(workersNum uint32, conn *net.UDPConn) *WorkerManager {
	return &WorkerManager{
		state:      states.New(),
		workersNum: workersNum,
		conn:       conn,
		// lri:     id.New(),
		send:    make(chan []byte, workersNum),
		receive: make(chan []byte, workersNum),
		ping:    make(chan []byte, workersNum),
		exit:    make(chan os.Signal),
		err:     make(chan error, workersNum),
	}
}
