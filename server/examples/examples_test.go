package examples_test

import (
	"log"

	"github.com/YarikRevich/game-networking/server/internal/handlers"
	"github.com/YarikRevich/game-networking/server/pkg/config"
	"github.com/YarikRevich/game-networking/server/pkg/connector"
)

func ExampleConnect() {
	conn, _ := connector.Connect(config.Config{
		IP:   "127.0.0.1",
		Port: "9999",
	})

	

	if err := conn.EstablishListening(); err != nil{
		log.Fatalln(err)
	}

	conn.InitWorkers(4)

	handlers.AddHandler("ping", func() []byte {
		return []byte("ping")
	})

	conn.WaitForInterrupt()

	//Output: Workable connector
}
