package listener

import (
	"github.com/YarikRevich/game-networking/common"
	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/server/internal/establisher"
)

func Listen(conf config.Config) (common.Listener, error){
	return establisher.New(conf)
}