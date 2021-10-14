package server

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/YarikRevich/game-networking/pkg/protocol"

	"github.com/sirupsen/logrus"

	"github.com/YarikRevich/game-networking/tools/buffer"
	"github.com/YarikRevich/game-networking/tools/creators"
)

type establisher struct {
	buffer *buffer.Buffer

	addr *net.UDPAddr
	conn *net.UDPConn

	handlers map[string]func(data interface{}) (interface{}, error)

	closeC  chan int
	surveyC *time.Ticker
}

func (e *establisher) establishListening() {
	conn, err := net.ListenUDP("udp", e.addr)
	if err != nil {
		logrus.Fatal(err)
	}
	e.conn = conn
}

func (e *establisher) setConfig(conf config.Config) {
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil {
		logrus.Fatal(err)
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	e.addr = addr
}

func (e *establisher) run() {
	for {
		select {
		case <-e.closeC:
			return
		case <-e.surveyC.C:
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
				if err != nil{
					logrus.Fatal(err)
				}

				p.Msg = r

				b, err := json.Marshal(p)
				if err != nil {
					continue
				}

				if _, err = e.conn.WriteTo(b, addr); err != nil {
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

func (e *establisher) AddHandler(name string, callback func(data interface{}) (interface{}, error)) {
	e.handlers[name] = callback
}

func (e *establisher) CallHandler(name string, data interface{}) (interface{}, error) {
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

func NewEstablisher(conf config.Config) Listener {
	e := &establisher{
		buffer:   buffer.New(),
		handlers: make(map[string]func(data interface{}) (interface{}, error)),
		closeC:   make(chan int, 1),
		surveyC:  time.NewTicker(time.Microsecond * 350),
	}
	e.setConfig(conf)
	e.establishListening()
	go e.run()
	return e
}
