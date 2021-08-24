package workers

import (
	"net"

	"github.com/YarikRevich/game-networking/protocol/pkg/id"
	"github.com/YarikRevich/game-networking/protocol/pkg/models"
	"github.com/YarikRevich/game-networking/server/internal/handlers"
	"github.com/YarikRevich/game-networking/server/internal/table"
	"github.com/YarikRevich/game-networking/server/tools/buffer"
)

type WorkerManager struct {
	count int //Count of workers

	tab *table.Table

	lri  *id.LocalRequestID

	buff *buffer.Buffer
	conn net.PacketConn

	err chan error
}

func (wm *WorkerManager) Run() {
	for i := 0; i < wm.count; i++ {
		go wm.worker()
	}
}

func (wm *WorkerManager) worker() {
	for {
		buff, ok := wm.buff.GetFromBuffer().([]byte)
		if !ok {
			continue
		}

		_, addr, err := wm.conn.ReadFrom(buff)
		if err != nil {
			wm.err <- err
		}

		if models.IsProtocolMsg(buff) {
			wm.tab.Add(addr.String(), buff)
		}

		msg, err := models.UnmarshalProtocol(buff)
		if err != nil {
			continue
		}

		res := handlers.CallHandler(msg.Procedure)
		if res != nil{
			continue
		}

		if _, err := wm.conn.WriteTo(res, addr); err != nil{
			continue
		}

		if cap(buff) < 1<<20 {
			wm.buff.PutToBuffer(buff[:0])
		}
	}
}

func (wm *WorkerManager) Error() error {
	return <-wm.err
}

func New(count int, conn net.PacketConn) *WorkerManager {
	return &WorkerManager{
		count: count,
		conn:  conn,
		lri: id.New(),
		tab:   table.New(),
		buff:  buffer.New(),
		err:   make(chan error, count),
	}
}
