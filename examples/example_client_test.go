package examples_test

import (
	"fmt"
	"syscall"
	"time"

	"github.com/YarikRevich/game-networking/server/pkg/handlers"
	"github.com/YarikRevich/game-networking/server/pkg/config"
	"github.com/YarikRevich/game-networking/server/pkg/listener"
)

func ExampleConnect() {
	conn, _ := listener.Listen(config.Config{
		IP:   "127.0.0.1",
		Port: "9090",
	})

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
