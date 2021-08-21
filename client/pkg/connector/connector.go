package connector

import (
	"net"

	"github.com/YarikRevich/game-networking/client/internal/establisher"
	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/internal/workers"
	"github.com/YarikRevich/game-networking/client/pkg/config"
	"github.com/YarikRevich/game-networking/client/tools"
)

func Connect(conf config.Config) (*establisher.Connector, error){
	createdAddr, err := tools.CreateAddr(conf.IP, conf.Port)
	if err != nil{
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil{
		return nil, err
	}

	// wmanager := workers.New(conf.WorkersCount, conn.GetConn())

	conn := establisher.NewConnector(
		addr, timeout.NewTimeout(conf.PingerAddr))


	return conn, nil
}