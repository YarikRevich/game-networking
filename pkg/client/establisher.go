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
)

var poolBuff = buffer.New()

type establisher struct {
	sync.Mutex

	addr      *net.UDPAddr
	conn      *net.UDPConn
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

func (e *establisher) Call(procedure string, src interface{}, dst interface{}, errc func(error), ank bool) {
	if errc == nil {
		errc(errors.New("error callback mustn't be nil"))
		return
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.IsValid() && (dstVal.Kind() != reflect.Ptr && dstVal.Kind() != reflect.Slice){
		errc(errors.New("dst should be a pointer or nil"))
		return
	}
		hash, err := GenerateHashSum(src)
		if err != nil{
			errc(err)
			return
		}

		m := protocol.Protocol{Procedure: procedure, Msg: src, HashSum: hash}
		b, err := json.Marshal(m)
		if err != nil {
			errc(err)
			return
		}

	main:
		for {
		write:
			for {
				if err := e.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
					errc(err)
				}
			
				_, err = e.conn.Write(b)
				if err, ok := err.(net.Error); ok && err.Timeout() {
					continue write
				}
				if err != nil {
					errc(err)
				}
			

				break write
			}

			buff := poolBuff.GetFromBuffer().([]byte)

		read:
			for {
				if err := e.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
					errc(err)
				}
				
				n, err := e.conn.Read(buff)
				if n == 0{
					continue main
				}
				if err, ok := err.(net.Error); ok && err.Timeout() {
					continue read
				}
				if err != nil {
					errc(err)
				}

				break read
			}

			buff = bytes.Trim(buff, "\x00")

			var p protocol.Protocol
			err = json.Unmarshal(buff, &p)
			if err != nil {
				errc(err)
				return
			}
			if p.HashSum != hash || p.Procedure != procedure {
				continue main
			}
			
			if cap(buff) <= 40*1024 {
				poolBuff.PutToBuffer(buff[:0])
			}


			ParseToDst(p.Msg, dstVal)

			// fmt.Println(reflect.ValueOf(t).Convert(dstVal.Elem().Type()))
			// fmt.Println(p.Msg.(ResultStub))
			
			// if dstVal.IsValid() && p.Msg != nil {
			// 	dstVal.Elem().Set(reflect.ValueOf(p.Msg))
			// }
			
			
			break main
		}
}

func (e *establisher) Close() error {
	return e.conn.Close()
}

func NewEstablisher(conf config.Config) (Dialer, error) {
	e := new(establisher)
	if err := e.setConfig(conf); err != nil {
		return nil, err
	}
	if err := e.establishConnection(); err != nil {
		return nil, err
	}
	return e, nil
}
