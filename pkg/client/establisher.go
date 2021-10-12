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
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/game-networking/tools/buffer"
	"github.com/YarikRevich/game-networking/tools/creators"
	"github.com/sirupsen/logrus"
)

var poolBuff = buffer.New()

type establisher struct {
	sync.Mutex

	addr *net.UDPAddr
	conn *net.UDPConn
}

func (e *establisher) establishConnection() {
	conn, err := net.DialUDP("udp", nil, e.addr)
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

func (e *establisher) Call(procedure string, src interface{}, dst interface{}) {
	dstVal := reflect.ValueOf(dst)
	if dstVal.IsValid() && (dstVal.Kind() != reflect.Ptr && dstVal.Kind() != reflect.Slice) {
		logrus.Fatal(errors.New("dst should be a pointer or nil"))
	}
	hash, err := GenerateHashSum(src)
	if err != nil {
		logrus.Fatal(err)
	}

	m := protocol.Protocol{Procedure: procedure, Msg: src, HashSum: hash}
	b, err := json.Marshal(m)
	if err != nil {
		logrus.Fatal(err)
	}

main:
	for {
	write:
		for {
			if err := e.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
				logrus.Fatal(err)
			}

			_, err = e.conn.Write(b)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue write
			}
			if err != nil {
				logrus.Fatal(err)
			}

			break write
		}

		buff := poolBuff.GetFromBuffer().([]byte)

	read:
		for {
			if err := e.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
				logrus.Fatal(err)
			}

			n, err := e.conn.Read(buff)
			if n == 0 {
				continue main
			}
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue read
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

		if dstVal.IsValid() && p.Msg != nil {
			ParseToDst(p.Msg, dstVal)
		}

		break main
	}
}

func (e *establisher) Close() error {
	return e.conn.Close()
}

func NewEstablisher(conf config.Config) Dialer {
	e := new(establisher)
	e.setConfig(conf)
	e.establishConnection()
	return e
}
