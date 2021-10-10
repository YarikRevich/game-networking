package server

import (
	"syscall"
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/franela/goblin"
)

func TestServer(t *testing.T){
	g := goblin.Goblin(t)

	g.Describe("TestServer", func() {
			c, err := Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
			if err != nil{
				t.Fatal(err)
			}

			go func(){
				time.Sleep(150 * time.Second)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()

			t.Fatal(c.WaitForInterrupt())
	})
}