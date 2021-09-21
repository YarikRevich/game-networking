package workers

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	"fmt"
	"net"
	"os"

	// "os"
	"runtime"

	// "github.com/YarikRevich/game-networking/protocol/pkg/id"
	// "github.com/YarikRevich/game-networking/protocol/pkg/models"
	// "github.com/YarikRevich/game-networking/server/pkg/handlers"
	// "github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	// "github.com/YarikRevich/game-networking/server/pkg/handlers"
	// "github.com/YarikRevich/game-networking/server/internal/table"
	// "github.com/YarikRevich/game-networking/server/tools/buffer"
)

type WorkerManager struct {
	table map[string]map[string]protocol.Protocol

	// lri  *id.LocalRequestID

	// buff *buffer.Buffer
	conn *net.UDPConn

	err chan error
}

func (wm *WorkerManager) Run() {
	// for i := 0; i < wm.count; i++ {
	// 	go wm.worker()
	// }

	go wm.worker()
}

func (wm *WorkerManager) worker() {
	for {

		var buffer []byte

		// buff, ok := .GetFromBuffer().([]byte)
		// if !ok {
		// 	continue
		// }

		_, _, err := wm.conn.ReadFromUDP(buffer)
		if err != nil {
			wm.err <- err
		}

		// if n != 0{
			fmt.Fprintln(os.Stderr, err)
		// }
		


		// var p protocol.Protocol

		// json.Unmarshal(buffer, &p)

		// fmt.Fprintln(os.Stderr, addr.String())

		// wm.table[addr.String()][p.Procedure] = p



		// if models.IsProtocolMsg(buff) {
		// 	wm.tab.Add(addr.String(), buff)
		// }

		// msg, err := models.UnmarshalProtocol(buff)
		// if err != nil {
		// 	continue
		// }

		// res := handlers.CallHandler(p.Procedure, p.Msg)
		// if res != nil{
		// 	continue
		// }

		// wm.table[addr.String()][p.Procedure] = res


		// if _, err := wm.conn.WriteTo([]byte("itworks"), addr); err != nil{
		// 	continue
		// }

		// if cap(buff) < 1<<20 {
		// 	wm.buff.PutToBuffer(buff[:0])
		// }
	}
}

func (wm *WorkerManager) Error() error {
	return <-wm.err
}

func New(conn *net.UDPConn) *WorkerManager {
	return &WorkerManager{
		conn:  conn,
		// lri: id.New(),
		table:   make(map[string]map[string]protocol.Protocol),
		// buff:  buffer.New(),
		err:   make(chan error, runtime.NumCPU()),
	}
}
