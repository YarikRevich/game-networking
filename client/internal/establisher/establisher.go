package establisher

import (
	"net"
	// "time"

	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/internal/workers"
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/wrapper/pkg/wrapper"
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
)

type Establisher struct {
	addr    *net.UDPAddr
	conn    *net.UDPConn
	timeout *timeout.Timeout

	wmanager *workers.WorkerManager
	p        protocol.Protocol
}

func (e *Establisher) EstablishConnection(workersNum uint32) error {
	conn, err := net.DialUDP("udp", nil, e.addr)
	if err != nil {
		return err
	}
	e.conn = conn

	e.initTimeouts()

	e.initWorkers(workersNum)
	return nil
}

func (e *Establisher) Read() (interface{}, error) {
	m, err := e.wmanager.Read()
	if err != nil {
		return nil, err
	}
	w := wrapper.UseWrapper()
	if err = w.Unmarshal(m); err != nil {
		return nil, err
	}
	return w.GetBase(), nil
}

func (e *Establisher) Write(m interface{}) error {
	w := wrapper.UseWrapper()
	w.SetBase(m)

	b, err := w.Marshal()
	if err != nil {
		return err
	}
	e.wmanager.Write(b)
	return nil
}

func (e *Establisher) Close() error {
	return e.conn.Close()
}

func (e *Establisher) initTimeouts() error {
	return e.timeout.EstimateProperTimout()
}

func (e *Establisher) initWorkers(count uint32) {
	e.wmanager = workers.New(count, e.conn)
	e.wmanager.Run()
}

// func (e *Establisher) SetReadDeadLine() error {
// 	rt := e.timeout.GetReadTimeout()
// 	return e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(rt)))
// }

// func (e *Establisher) SetWriteDeadLine() error {
// 	wt := e.timeout.GetWriteTimeout()
// 	return e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(wt)))
// }

// func (e *Establisher) WorkerManager() *workers.WorkerManager {
// 	return e.wmanager
// }

//, pm models.ProtocolManager

func New(addr *net.UDPAddr, timeout *timeout.Timeout, p protocol.Protocol) *Establisher {
	return &Establisher{
		addr:    addr,
		timeout: timeout,
		p:       p,
	}
}
