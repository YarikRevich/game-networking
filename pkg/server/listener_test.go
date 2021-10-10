package server

import (
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/franela/goblin"
)

func TestServer(t *testing.T){
	g := goblin.Goblin(t)

	g.Describe("TestServer", func() {
		g.It("TestListener", func(){
			c, err := Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
			g.Assert(err).IsNil()

			go func(){
				time.Sleep(3 * time.Second)
				g.Assert(c.Close()).IsNil()
			}()
			g.Assert(c.WaitForInterrupt()).IsNil()
		})
	})
}