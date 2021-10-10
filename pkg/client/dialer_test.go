package client

import (
	"syscall"
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/YarikRevich/game-networking/pkg/server"
	"github.com/franela/goblin"
)

func TestDialer(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("TestDialer", func() {
		go func(){
			c, err := server.Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
			c.AddHandler("test", func (m []byte)([]byte, error)  {
					return []byte("itworks"), nil
			})
			g.Assert(err).IsNil()
	
			go func(){
				time.Sleep(15 * time.Second)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()
			g.Assert(c.WaitForInterrupt()).IsNil()
		}()

		time.Sleep(2 * time.Second)

		clientConfig := config.Config{
			IP: "127.0.0.1", 
			Port: "8090", 
		}

		d, err := Dial(clientConfig)
		g.Assert(err).IsNil()

		g.After(func() {
			g.Assert(d.Close())
		})

		g.It("dial call", func() {
	

			var dst string
			src := "itworks"
			d.Call("test", src, &dst, func(e error) { t.Fatal(err) }, false)
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")

		})

		g.It("dial call, test ank", func() {
			var dst string
			src := "itworks"
			d.Call("test", src, &dst, func(e error) { t.Fatal(err) }, true)
			
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")

			dst = ""
			d.Call("test", src, &dst, func(e error) { t.Fatal(e) }, true)
			
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")
		})
	})
}
