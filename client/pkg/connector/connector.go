package connector

import (
	"net"

	"github.com/YarikRevich/game-networking/client/internal/establisher"
	"github.com/YarikRevich/game-networking/client/internal/timeout"
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

	return establisher.NewConnector(
		addr, timeout.NewTimeout(conf.PingerAddr)), nil
}