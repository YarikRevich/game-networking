package client

import "github.com/YarikRevich/game-networking/config"

func Dial(conf config.Config) (Dialer, error) {
	return NewEstablisher(conf)
}
