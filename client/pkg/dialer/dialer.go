package dialer

import (
	"github.com/YarikRevich/game-networking/common"
	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/client/internal/establisher"
)

func Dial(conf config.Config) (common.Dialer, error) {
	return establisher.New(conf)
}
