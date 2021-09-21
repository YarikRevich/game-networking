package listener

import (
	"net"

	"github.com/YarikRevich/game-networking/server/internal/establisher"
	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/tools/pkg/creators"
)

func Listen(conf config.Config) (establisher.EstablishAwaiter, error){
	createdAddr, err := creators.CreateAddr(conf.IP, conf.Port)
	if err != nil{
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", createdAddr)
	if err != nil{
		return nil, err
	}

	e := establisher.New(addr)

	if err := e.EstablishListening(); err != nil{
		return nil, err
	}

	e.InitWorkers()

	return e, nil
}