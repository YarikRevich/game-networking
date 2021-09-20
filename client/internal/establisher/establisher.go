package establisher

import (
	"bytes"
	"context"
	// "crypto/sha256"
	"encoding/json"
	// "fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	// "time"

	// "github.com/YarikRevich/game-networking/client/internal/timeout"
	// "github.com/YarikRevich/game-networking/client/internal/workers"
	"github.com/YarikRevich/game-networking/common"
	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"

	// "github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/wrapper/pkg/wrapper"
	//
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
)

type workerManager struct {
	write chan []byte
	read  chan []byte
	exit  chan os.Signal
	err   chan error
}

type establisher struct {
	*sync.Mutex

	// workerManager

	signalExit    chan os.Signal
	internalError error

	send []byte

	table map[string]protocol.Protocol

	addr    *net.UDPAddr
	conn    *net.UDPConn
	wrapper wrapper.Wrapper
	// timeout *timeout.Timeout

	// wmanager *workers.WorkerManager
	// p        protocol.Protocol
}

func (e *establisher) establishConnection() error {
	conn, err := net.DialUDP("udp", nil, e.addr)
	if err != nil {
		return err
	}
	e.conn = conn

	e.runWorkers()
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

func (e *establisher) runWorkers() {
	signal.Notify(e.signalExit, os.Interrupt)

	ctx, close := context.WithCancel(context.Background())

	go func() {
		for range e.signalExit {
			close()
		}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(ctx context.Context) {
			var buffer bytes.Buffer
			if _, err := io.Copy(&buffer, e.conn); err != nil {
				e.Lock()
				e.internalError = err
				e.Unlock()
			}
			var msg protocol.Protocol
			if err := json.Unmarshal(buffer.Bytes(), &msg); err != nil {
				e.Lock()
				e.internalError = err
				e.Unlock()
			}

			e.table[msg.Procedure] = msg
		}(ctx)
		go func(ctx context.Context) {
			if _, err := e.conn.Write(e.send); err != nil {
				e.Lock()
				e.internalError = err
				e.Unlock()
			}
		}(ctx)
	}
}

func (e *establisher) Call(procName string) interface{} {
	after := time.After(1 * time.Second)
	for {
		select {
		case <-after:
			return nil
		default:
			k, ok := e.table[procName]
			if ok {
				return k.Msg
			}
		}

	}
	// if err := e.wrapper.Unmarshal(<-e.read); err != nil {
	// 	return nil, err
	// }

	// if e.wrapper.GetField("hash_sum").([32]byte) == e.wrapper.GetBase() {
	// 	//validation by hash_sum
	// }

	// return e.wrapper.GetBase().([]byte), nil
}

func (e *establisher) Send(procName string, msg interface{}) {

	// e.wrapper.SetBase(m)
	// e.wrapper.SetField("hash_sum", sha256.Sum256([]byte(fmt.Sprintf("%v", m))))

	b, err := json.Marshal(protocol.Protocol{Procedure: procName, Msg: msg})
	if err != nil {
		e.Lock()
		e.internalError = err
		e.Unlock()
	}
	// b, err := e.wrapper.Marshal()
	// if err != nil {
	// 	return err
	// }
	e.Lock()
	e.send = b
	e.Unlock()
}

func (e *establisher) Error() error {
	return e.internalError
}

func (e *establisher) Close() error {
	return e.conn.Close()
}

// func (e *Establisher) initTimeouts() error {
// return e.timeout.EstimateProperTimout()
// }

// func (e *Establisher) SetReadDeadLine() error {
// 	rt := e.timeout.GetReadTimeout()
// 	return e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(rt)))
// }

// func (e *Establisher) SetWriteDeadLine() error {
// 	wt := e.timeout.GetWriteTimeout()
// 	return e.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(wt)))
// }

//, pm models.ProtocolManager

func New(conf config.Config) (common.Conn, error) {
	e := &establisher{
		Mutex: new(sync.Mutex),
		// workerManager: workerManager{
		// 	write: make(chan []byte),
		// 	read:  make(chan []byte, runtime.NumCPU()),

		// 	err:   make(chan error),
		// },
		signalExit: make(chan os.Signal),

		table:   make(map[string]protocol.Protocol),
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
