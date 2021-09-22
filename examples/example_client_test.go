package examples_test

import (
	"fmt"
	"syscall"
	"time"

	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/server/pkg/listener"
)

func ExampleConnect() {
	conn, _ := listener.Listen(config.Config{
		IP:   "127.0.0.1",
		Port: "9090",
	})

	conn.AddHandler("ping", func(data interface{}) ([]byte, error){
		return []byte("ping"), nil
	})

	go func(){
		time.Sleep(3 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	conn.WaitForInterrupt()

	fmt.Println("It works")
	//Output: It works
}
