package establisher

import (
	"net"
	"os"
	"os/signal"

	"github.com/YarikRevich/game-networking/server/internal/workers"
)

type EstablishAwaiter interface {
	WaitForInterrupt() error
	Close() error
}

type Establisher struct {
	addr     *net.UDPAddr
	conn     *net.UDPConn
	wmanager *workers.WorkerManager

	close chan int
}

func (e *Establisher) EstablishListening() error {
	conn, err := net.ListenUDP("udp", e.addr)
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

func (e *Establisher) InitWorkers() {
	e.wmanager = workers.New(e.GetConn())
	e.wmanager.Run()
}

func (e *Establisher) GetConn() *net.UDPConn {
	return e.conn
}

func (e *Establisher) WaitForInterrupt() error {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	for {
		select {
		case <- e.close:
			return nil
		case <- sig:
			return e.Close()
		}
	}
}

func (e *Establisher) Close() error {
	e.close <- 1
	return e.conn.Close()
}

func New(addr *net.UDPAddr) *Establisher {
	return &Establisher{
		addr: addr,
		close: make(chan int, 1),
	}
}
