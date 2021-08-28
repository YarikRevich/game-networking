package connector

import (
	"net"

	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/pkg/config"
	"github.com/YarikRevich/game-networking/client/pkg/establisher"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
)

func Connect(conf config.Config) (*establisher.Establisher, error) {
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil {
		return nil, err
	}

	conn := establisher.New(
		addr, timeout.NewTimeout(conf.PingerAddr))

	if err = conn.EstablishConnection(); err != nil {
		return nil, err
	}

	if err = conn.InitTimeouts(); err != nil {
		return nil, err
	}

	conn.InitWorkers(conf.Workers)

	return conn, nil
}
