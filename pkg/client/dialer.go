package client

import "github.com/YarikRevich/game-networking/pkg/config"

func Dial(conf config.Config) (Dialer, error) {
	return NewEstablisher(conf)
}
