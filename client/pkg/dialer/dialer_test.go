package dialer

import (
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/config"
	"github.com/YarikRevich/game-networking/server/pkg/listener"
	"github.com/franela/goblin"
)

func TestDialer(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("TestDialer", func() {
		go func(){
			c, err := listener.Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
			c.AddHandler("test", func (m interface{})([]byte, error)  {
					return []byte("itworks"), nil
			})
			g.Assert(err).IsNil()
	
			go func(){
				time.Sleep(15 * time.Second)
				g.Assert(c.Close()).IsNil()
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
			g.Assert(d.Call("test", src, &dst, func(e error) { panic(err) }, false)).IsNil()
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")

		})

		g.It("dial call, test ank", func() {
			var dst string
			src := "itworks"
			g.Assert(d.Call("test", src, &dst, func(e error) { panic(err) }, true)).IsNil()
			
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")

			g.Assert(d.Call("test", src, &dst, func(e error) { panic(err) }, true)).IsNil()
			
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")
			
		})
	})
}
