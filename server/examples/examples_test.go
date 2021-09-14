package examples_test

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/YarikRevich/game-networking/server/pkg/handlers"
	"github.com/YarikRevich/game-networking/server/pkg/config"
	"github.com/YarikRevich/game-networking/server/pkg/connector"
)

func ExampleConnect() {
	conn, _ := connector.Listen(config.Config{
		IP:   "127.0.0.1",
		Port: "9090",
	})

	if err := conn.EstablishListening(); err != nil{
		log.Fatalln(err)
	}

	conn.InitWorkers(4)

	handlers.AddHandler("ping", func(data interface{}) []byte{
		return []byte("ping")
	})

	go func(){
		time.Sleep(3 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	conn.WaitForInterrupt()

	fmt.Println("It works")
	//Output: It works
}
