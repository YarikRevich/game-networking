package establisher

import (
	"net"
	"os"
	"os/signal"

	"github.com/YarikRevich/game-networking/server/internal/workers"
)

type Establisher struct {
	addr     string
	conn     net.PacketConn
	wmanager *workers.WorkerManager
}

func (e *Establisher) EstablishListening() error {
	conn, err := net.ListenPacket("udp", e.addr)
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

func (e *Establisher) InitWorkers(count int) {
	e.wmanager = workers.New(count, e.GetConn())
	e.wmanager.Run()
}

func (e *Establisher) GetConn() net.PacketConn {
	return e.conn
}

func (e *Establisher) WaitForInterrupt()error{
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	for range sig{
		return e.close()
	}
	return nil
}

func (e *Establisher) close() error {
	return e.conn.Close()
}

func New(addr string) *Establisher {
	return &Establisher{
		addr: addr,
	}
}
