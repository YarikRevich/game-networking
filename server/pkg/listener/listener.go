package listener

import (
	"github.com/YarikRevich/game-networking/server/internal/establisher"
	"github.com/YarikRevich/game-networking/config"
)

func Listen(conf config.Config) (establisher.EstablishAwaiter, error){
	return establisher.New(conf)
}