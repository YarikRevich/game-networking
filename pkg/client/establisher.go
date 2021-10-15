package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/YarikRevich/game-networking/pkg/protocol"
	"github.com/YarikRevich/game-networking/tools/buffer"
	"github.com/YarikRevich/game-networking/tools/creators"
	"github.com/sirupsen/logrus"
)

var poolBuff = buffer.New()

type establisher struct {
	sync.Mutex

	addr *net.UDPAddr
	conn net.Conn

	repeatCount int
}

func (e *establisher) establishConnection() error {
	conn, err := net.Dial("udp", e.addr.String())
	if err != nil {
		return err
	}
	e.conn = conn
	return err
}

func (e *establisher) ping() bool {
	m := protocol.Protocol{Procedure: "ping"}
	b, err := json.Marshal(m)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = e.conn.Write(b)
	if err != nil {
		return false
	}

	buff := make([]byte, len("ping"))
	_, err = e.conn.Read(buff)
	return err == nil
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

func (e *establisher) Call(procedure string, src interface{}, dst interface{}) {
	dstVal := reflect.ValueOf(dst)

	if dstVal.IsValid() && (dstVal.Kind() != reflect.Ptr && dstVal.Kind() != reflect.Slice) {
		logrus.Fatal(errors.New("dst should be a pointer or nil"))
	}

	hash, err := GenerateHashSum(src)
	if err != nil {
		logrus.Fatal(err)
	}

	m := protocol.Protocol{Msg: src, HashSum: hash, Procedure: procedure}
	b, err := json.Marshal(m)
	if err != nil {
		logrus.Fatal(err)
	}

main:
	for {
		e.repeatCount++

		if e.repeatCount == 20 {
			return
		}
	write:
		for {
			if err = e.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
				logrus.Fatal(err)
			}

			_, err = e.conn.Write(b)

			var e net.Error
			if errors.As(err, &e) {
				if e.Timeout() {
					continue write
				}
			}

			if err != nil {
				logrus.Fatal(err)
			}

			break write
		}

		buff := poolBuff.GetFromBuffer().([]byte)

	read:
		for {
			if err = e.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
				logrus.Fatal(err)
			}

			var n int
			n, err = e.conn.Read(buff)
			if err != nil{
				logrus.Fatal(err)
			}
			if n == 0{
				continue main
			}

			var e net.Error
			if errors.As(err, &e) {
				if e.Timeout() {
					continue read
				}
			}

			if err != nil {
				logrus.Fatal(err)
			}

			break read
		}

		buff = bytes.Trim(buff, "\x00")

		var p protocol.Protocol
		err = json.Unmarshal(buff, &p)
		if err != nil {
			logrus.Fatal(err)
		}
		if p.HashSum != hash || p.Procedure != procedure {
			continue main
		}

		if cap(buff) <= 40*1024 {
			poolBuff.PutToBuffer(buff[:0])
		}

		var m struct{ Msg json.RawMessage }
		if err := json.Unmarshal(buff, &m); err != nil {
			logrus.Fatal(err)
		}

		if err := json.Unmarshal(m.Msg, &dst); err != nil {
			logrus.Fatal(err)
		}

		break main
	}
}

func (e *establisher) Close() error {
	return e.conn.Close()
}

func (e *establisher) IsConnected() bool{
	return e.ping()
}

func NewEstablisher(conf config.Config) (Dialer, error) {
	e := &establisher{repeatCount: -1}
	e.setConfig(conf)
	if err := e.establishConnection(); err != nil {
		return nil, err
	}
	return e, nil
}
