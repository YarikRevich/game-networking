package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/game-networking/tools/pkg/buffer"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
	"github.com/YarikRevich/wrapper/pkg/wrapper"
)

var poolBuff = buffer.New()

type establisher struct {
	sync.Mutex

	scheduler IScheduler
	addr      *net.UDPAddr
	conn      *net.UDPConn
	wrapper   wrapper.Wrapper
}

func (e *establisher) establishConnection() error {
	conn, err := net.DialUDP("udp", nil, e.addr)
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

func (e *establisher) Call(procedure string, src interface{}, dst interface{}, errc func(error), ank bool) error {
	if errc == nil {
		return errors.New("error callback mustn't be nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr && !dstVal.IsNil() {
		return errors.New("dst should be a pointer")
	}

	call := func() {
		m := protocol.Protocol{Procedure: procedure, Msg: src}
		b, err := json.Marshal(m)
		if err != nil {
			errc(err)
			return
		}

		_, err = e.conn.Write(b)
		if err != nil {
			errc(err)
			return
		}

		buff := poolBuff.GetFromBuffer().([]byte)

		_, err = e.conn.Read(buff)
		if err != nil {
			errc(err)
			return
		}

		buff = bytes.Trim(buff, "\x00")

		var p protocol.Protocol
		err = json.Unmarshal(buff, &p)
		if err != nil {
			errc(err)
			return
		}
		if cap(buff) <= 20*1024 {
			poolBuff.PutToBuffer(buff[:0])
		}

		if p.Msg == nil {
			errc(fmt.Errorf("message is empty: %w", err))
			return
		}

		if ank {
			e.scheduler.DecConfirmations()
		}

		if !dstVal.IsNil() {
			dstVal.Elem().Set(reflect.ValueOf(p.Msg))
		}
	}

	if e.scheduler.CountConfirmations() != 0 {
		e.scheduler.Schedule(call)
		return nil
	}

	if ank {
		e.scheduler.IncConfirmations()
	}

	go call()

	return nil
	// if e.wrapper.GetField("hash_sum").([32]byte) == e.wrapper.GetBase() {
	// 	//validation by hash_sum
	// }

	// return e.wrapper.GetBase().([]byte), nil
}

func (e *establisher) Close() error {
	return e.conn.Close()
}

// func (e *Establisher) SetReadDeadLine() error {
// 	rt := e.timeout.GetReadTimeout()
// 	return e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(rt)))
// }

// func (e *Establisher) SetWriteDeadLine() error {
// 	wt := e.timeout.GetWriteTimeout()
// 	return e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(wt)))
// }

//, pm models.ProtocolManager

func NewEstablisher(conf config.Config) (Dialer, error) {
	e := &establisher{
		wrapper:   wrapper.UseWrapper(),
		scheduler: NewScheduler(20),
	}
	if err := e.setConfig(conf); err != nil {
		return nil, err
	}
	if err := e.establishConnection(); err != nil {
		return nil, err
	}
	return e, nil
}
