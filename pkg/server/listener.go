package server

import "github.com/YarikRevich/game-networking/pkg/config"

func Listen(conf config.Config) (Listener, error){
	return NewEstablisher(conf)
}