package dialer

import (
	"encoding/json"
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
	
		clientConfig := config.Config{
			IP: "127.0.0.1", 
			Port: "8090", 
			Decoder: json.Unmarshal,
			Encoder: json.Marshal,
		}

		g.It("dial", func() {
			d, err := Dial(clientConfig)
			g.Assert(err).IsNil()

			defer func() {
				g.Assert(d.Close())
			}()
		})

		g.It("dial call", func() {
			d, err := Dial(clientConfig)
			g.Assert(err).IsNil()

			var dst string
			src := "itworks"
			g.Assert(d.Call("test", src, &dst, make(chan error, 32))).IsNil()
			time.Sleep(2 * time.Second)
			g.Assert(dst).Eql("itworks")

			defer func() {
				g.Assert(d.Close())
			}()
		})
	})
}
