package connector

import (
	"net"

	"github.com/YarikRevich/game-networking/server/internal/establisher"
	"github.com/YarikRevich/game-networking/server/pkg/config"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
)

func Listen(conf config.Config) (*establisher.Establisher, error){
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil{
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil{
		return nil, err
	}

	return establisher.New(addr.String()), nil
}