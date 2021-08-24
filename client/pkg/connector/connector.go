package connector

import (
	"net"

	"github.com/YarikRevich/game-networking/client/internal/establisher"
	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/pkg/config"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
)

func Connect(conf config.Config) (*establisher.Establisher, error){
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil{
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil{
		return nil, err
	}

	return establisher.New(
		addr, timeout.NewTimeout(conf.PingerAddr)), nil
}