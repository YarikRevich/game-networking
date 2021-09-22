package establisher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/YarikRevich/game-networking/common"
	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/game-networking/tools/pkg/buffer"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
	"github.com/YarikRevich/wrapper/pkg/wrapper"
)

var (
	poolBuff = buffer.New()
)

type establisher struct {
	sync.Mutex

	addr    *net.UDPAddr
	conn    *net.UDPConn
	wrapper wrapper.Wrapper
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
	e.wrapper.SetDecoder(conf.Decoder)
	e.wrapper.SetEncoder(conf.Encoder)

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

// func (e *establisher) runWorkers() {
// 	signal.Notify(e.signalExit, os.Interrupt)

// 	ctx, close := context.WithCancel(context.Background())

// 	go func() {
// 		for range e.signalExit {
// 			close()
// 		}
// 	}()

// 	for i := 0; i < runtime.NumCPU(); i++ {
// 		go func(ctx context.Context) {
// 			var buffer bytes.Buffer
// 			if _, err := io.Copy(&buffer, e.conn); err != nil {
// 				e.errorMutex.Lock()
// 				e.internalError = err
// 				e.errorMutex.Unlock()
// 			}
// 			var msg protocol.Protocol
// 			if err := json.Unmarshal(buffer.Bytes(), &msg); err != nil {
// 				e.errorMutex.Lock()
// 				e.internalError = err
// 				e.errorMutex.Unlock()
// 			}

// 			e.msgMutex.Lock()
// 			e.table[msg.Procedure] = msg
// 			e.msgMutex.Unlock()
// 		}(ctx)
// 		go func(ctx context.Context) {
// 			if _, err := e.conn.Write(e.send); err != nil {
// 				e.errorMutex.Lock()
// 				e.internalError = err
// 				e.errorMutex.Unlock()
// 			}
// 		}(ctx)
// 	}
// }

func (e *establisher) Call(procName string, src interface{}, dst interface{}, ec chan error)error {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr{
		return errors.New("dst should be a pointer")
	}
	go func(){
		m := protocol.Protocol{Procedure: procName, Msg: src}
		b, err := json.Marshal(m)
		if err != nil{
			ec <- err
			return
		}
	
		_, err = e.conn.Write(b)
		if err != nil{
			ec <- err
			return
		}

		buff := poolBuff.GetFromBuffer().([]byte)
	
		_, err = e.conn.Read(buff)
		if err != nil{
			ec <- err
			return
		}

		buff = bytes.Trim(buff, "\x00")

		var p protocol.Protocol
		err = json.Unmarshal(buff, &p)
		if err != nil{
			ec <- err
			return
		}
		if cap(buff) <= 20 * 1024{
			poolBuff.PutToBuffer(buff[:0])
		}

		if p.Msg == nil{
			ec <- fmt.Errorf("message is empty: %w", err)
			return
		}

		dstVal.Elem().Set(reflect.ValueOf(p.Msg))
	}()
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

func New(conf config.Config) (common.Dialer, error) {
	e := &establisher{
		wrapper: wrapper.UseWrapper(),
	}
	if err := e.setConfig(conf); err != nil {
		return nil, err
	}
	if err := e.establishConnection(); err != nil {
		return nil, err
	}
	return e, nil
}
