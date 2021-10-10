package server

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"

	"github.com/YarikRevich/game-networking/tools/buffer"
	"github.com/YarikRevich/game-networking/tools/creators"
)

type establisher struct {
	buffer *buffer.Buffer

	addr *net.UDPAddr
	conn *net.UDPConn

	handlers map[string]func(data interface{}) ([]byte, error)

	closeC chan int
	surveyC *time.Ticker
}

func (e *establisher) establishListening() error {
	conn, err := net.ListenUDP("udp", e.addr)
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

func (e *establisher) setConfig(conf config.Config) error {
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil {
		return err
	}
	e.addr = addr
	return nil
}

func (e *establisher) run() {
	for {
		select {
		case <-e.closeC:
			return
		case <- e.surveyC.C:
			buff := e.buffer.GetFromBuffer().([]byte)

			_, addr, err := e.conn.ReadFromUDP(buff)
			if err != nil {
				continue
			}

			if addr != nil {
				buff = bytes.Trim(buff, "\x00")

				var p protocol.Protocol

				if err := json.Unmarshal(buff, &p); err != nil {
					continue
				}

				r, err := e.CallHandler(p.Procedure, p.Msg)
				
				p.Msg = string(r)
				p.Error = err

				b, err := json.Marshal(p)
				if err != nil {
					continue
				}

				if _, err := e.conn.WriteTo(b, addr); err != nil {
					continue
				}
			}

			if cap(buff) <= 60*1024 {
				e.buffer.PutToBuffer(buff[:0])
			}
		}
	}
}

func (e *establisher) WaitForInterrupt() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	for {
		select {
		case <-e.closeC:
			return nil
		case <-sig:
			return e.close()
		}
	}
}

func (e *establisher) AddHandler(name string, callback func(data interface{}) ([]byte, error)) {
	e.handlers[name] = callback
}

func (e *establisher) CallHandler(name string, data interface{}) ([]byte, error) {
	if v, ok := e.handlers[name]; ok {
		return v(data)
	}
	return nil, nil
}

func (e *establisher) close() error {
	e.closeC <- 1
	e.surveyC.Stop()
	return e.conn.Close()
}

func NewEstablisher(conf config.Config) (Listener, error) {
	e := &establisher{
		buffer: buffer.New(),
		handlers: make(map[string]func(data interface{}) ([]byte, error)),
		closeC:  make(chan int, 1),
		surveyC: time.NewTicker(time.Microsecond * 350),
	}
	if err := e.setConfig(conf); err != nil {
		return nil, err
	}
	if err := e.establishListening(); err != nil {
		return nil, err
	}
	go e.run()
	return e, nil
}
