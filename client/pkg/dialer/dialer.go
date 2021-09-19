package dialer

import (
	"net"

	"github.com/YarikRevich/game-networking/protocol/pkg/protocol"
	"github.com/YarikRevich/game-networking/common"
	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/pkg/config"
	"github.com/YarikRevich/game-networking/client/internal/establisher"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
)

//, pm models.ProtocolManager
func Dial(conf config.Config, p protocol.Protocol) (common.Conn, error) {
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil {
		return nil, err
	}
	//pm

	conn := establisher.New(
		addr, timeout.NewTimeout(conf.PingerAddr), p)

	if err = conn.EstablishConnection(); err != nil {
		return nil, err
	}

	if err = conn.InitTimeouts(); err != nil {
		return nil, err
	}

	conn.InitWorkers(conf.WorkersNum)

	return conn, nil
}
