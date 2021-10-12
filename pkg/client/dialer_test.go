package client

import (
	"syscall"
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/YarikRevich/game-networking/pkg/server"
	"github.com/franela/goblin"
)

type ResultStub struct {
	Result string
}

func TestDialer(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("TestDialer", func() {
		go func() {
			c := server.Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
			c.AddHandler("1", func(m interface{}) (interface{}, error) {
				return ResultStub{Result: "1"}, nil
			})
			c.AddHandler("2", func(m interface{}) (interface{}, error) {
				return ResultStub{Result: "2"}, nil
			})
			c.AddHandler("3", func(m interface{}) (interface{}, error) {
				return ResultStub{Result: "3"}, nil
			})

			go func() {
				time.Sleep(15 * time.Second)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()
			g.Assert(c.WaitForInterrupt()).IsNil()
		}()

		time.Sleep(2 * time.Second)

		clientConfig := config.Config{
			IP:   "127.0.0.1",
			Port: "8090",
		}

		d := Dial(clientConfig)

		g.After(func(){
			time.Sleep(time.Second * 5)
			g.Assert(d.Close()).IsNil()
		})

		g.It("dial call, test ank", func() {
			var first ResultStub
			d.Call("1", nil, &first)


			var second ResultStub
			d.Call("2", nil, &second)

			var stub ResultStub
			d.Call("3", nil, &stub)
			d.Call("3", nil, &stub)

			var third ResultStub
			d.Call("3", nil, &third)

			g.Assert(first.Result).Eql("1")
			g.Assert(second.Result).Eql("2")
			g.Assert(third.Result).Eql("3")
		})
	})
}
