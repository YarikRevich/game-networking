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
			c, err := server.Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
			c.AddHandler("1", func(m []byte) (interface{}, error) {
				return ResultStub{Result: "1"}, nil
			})
			c.AddHandler("2", func(m []byte) (interface{}, error) {
				return ResultStub{Result: "2"}, nil
			})
			c.AddHandler("3", func(m []byte) (interface{}, error) {
				
				return ResultStub{Result: "3"}, nil
			})
			g.Assert(err).IsNil("Connection refused")

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

		d, err := Dial(clientConfig)
		g.Assert(err).IsNil()

		g.After(func(){
			time.Sleep(time.Second * 5)
			g.Assert(d.Close()).IsNil()
		})

		g.It("dial call, test ank", func() {
			

			var first ResultStub
			d.Call("1", nil, &first, func(e error) { t.Fatal(err) }, true)

			var second ResultStub
			d.Call("2", nil, &second, func(e error) { t.Fatal(e) }, true)

			var stub ResultStub
			d.Call("3", nil, &stub, func(e error) { t.Fatal(e) }, true)
			d.Call("3", nil, &stub, func(e error) { t.Fatal(e) }, true)
			d.Call("3", nil, &stub, func(e error) { t.Fatal(e) }, true)
			d.Call("3", nil, &stub, func(e error) { t.Fatal(e) }, true)

			var third ResultStub
			d.Call("3", nil, &third, func(e error) { t.Fatal(e) }, true)

			t.Log(first, second, third)
			g.Assert(first.Result).Eql("1")
			g.Assert(second.Result).Eql("2")
			g.Assert(third.Result).Eql("3")
		})
	})
}
