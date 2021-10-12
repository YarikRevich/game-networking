package client

import (
	"os"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/sirupsen/logrus"
)

func init(){
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{
		"module": "client",
	}})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stderr)
}

func Dial(conf config.Config) Dialer {
	return NewEstablisher(conf)
}
