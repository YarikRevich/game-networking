package tests

import (
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/server/pkg/config"
	"github.com/YarikRevich/game-networking/server/pkg/listener"
	"github.com/franela/goblin"
)

func TestServer(t *testing.T){
	g := goblin.Goblin(t)

	g.Describe("TestServer", func() {
		g.It("TestListener", func(){
			c, err := listener.Listen(config.Config{IP: "127.0.0.1", Port: "8090", WorkersNum: 4})
			g.Assert(err).IsNil()

			go func(){
				time.Sleep(3 * time.Second)
				g.Assert(c.Close()).IsNil()
			}()
			g.Assert(c.WaitForInterrupt()).IsNil()
		})
	})
}