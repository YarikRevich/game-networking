package establisher

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/YarikRevich/game-networking/config"
	// "github.com/YarikRevich/game-networking/server/internal/workers"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
)

type EstablishAwaiter interface {
	WaitForInterrupt() error
	Close() error
}

type establisher struct {
	addr     *net.UDPAddr
	conn     *net.UDPConn

	close chan int
}

func (e *establisher) establishListening() error {
	conn, err := net.ListenUDP("udp", e.addr)
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

func (e *establisher) setConfig(conf config.Config)error{
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil{
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil{
		return err
	}
	e.addr = addr
	return nil
}

// func (e *Establisher) InitWorkers() {
// 	e.wmanager = workers.New(e.GetConn())
// 	e.wmanager.Run()
// }

// func (e *Establisher) GetConn() *net.UDPConn {
// 	return e.conn
// }
func (e *establisher) run(){
	for {
		select {
		case <- e.close:
			return
		default:
			buff := make([]byte, 32)

			_, addr, err := e.conn.ReadFromUDP(buff)
			if err != nil{
				continue
			}
			if addr != nil{
				fmt.Println(addr.String())
			}
		}
	}
}


func (e *establisher) WaitForInterrupt() error {
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

func (e *establisher) Close() error {
	e.close <- 1
	return e.conn.Close()
}

func New(conf config.Config) (EstablishAwaiter, error) {
	e := &establisher{close: make(chan int, 1)}
	if err := e.setConfig(conf); err != nil {
		return nil, err
	}
	if err := e.establishListening(); err != nil {
		return nil, err
	}
	e.run()
	return e, nil
}
