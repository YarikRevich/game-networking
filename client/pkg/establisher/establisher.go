package establisher

import (
	"net"
	"time"

	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/internal/workers"
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
)

type Establisher struct {
	addr     *net.UDPAddr
	conn     *net.UDPConn
	timeout  *timeout.Timeout
	wmanager *workers.WorkerManager
}

func (e *Establisher) EstablishConnection() error {
	conn, err := net.DialUDP("udp", nil, e.addr)
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

func (e *Establisher) InitTimeouts() error {
	return e.timeout.EstimateProperTimout()
}

func (e *Establisher) InitWorkers(count uint32) {
	e.wmanager = workers.New(count, e.GetConn())
	e.wmanager.Run()
}

func (e *Establisher) SetReadDeadLine() {
	rt := e.timeout.GetReadTimeout()
	e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(rt)))
}

func (e *Establisher) SetWriteDeadLine() {
	wt := e.timeout.GetWriteTimeout()
	e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(wt)))
}

func (e *Establisher) WorkerManager() *workers.WorkerManager {
	return e.wmanager
}

func (e *Establisher) Close() error {
	return e.conn.Close()
}

func (e *Establisher) GetConn() *net.UDPConn {
	return e.conn
}

func (e *Establisher) Ping() error {
	// return e.wmanager.Ping()
	return nil
}

func New(addr *net.UDPAddr, timeout *timeout.Timeout) *Establisher {
	return &Establisher{
		addr:    addr,
		timeout: timeout,
	}
}
