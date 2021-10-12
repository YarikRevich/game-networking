package server

import (
	"os"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/sirupsen/logrus"
)

func init(){
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{
		"module": "server",
	}})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stderr)
}

func Listen(conf config.Config) Listener{
	return NewEstablisher(conf)
}