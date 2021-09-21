package dialer

import (
	"encoding/json"
	"testing"
	// "time"

	// "time"

	"github.com/YarikRevich/game-networking/config"
	// "github.com/YarikRevich/game-networking/server/pkg/handlers"
	// "github.com/YarikRevich/game-networking/server/pkg/listener"
	"github.com/franela/goblin"
)

func TestDialer(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("TestDialer", func() {
		// go func(){
		// 	c, err := listener.Listen(config.Config{IP: "127.0.0.1", Port: "8090"})
		// 	handlers.AddHandler("test", func (interface{})[]byte  {
		// 		return []byte("itworks")
		// 	})
		// 	g.Assert(err).IsNil()
	
		// 	go func(){
		// 		time.Sleep(15 * time.Second)
		// 		g.Assert(c.Close()).IsNil()
		// 	}()
		// 	g.Assert(c.WaitForInterrupt()).IsNil()
		// }()

		// time.Sleep(5 * time.Second)
	
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

			r := d.Call("test", "IT WORKS")
			t.Log(r)
			g.Assert(r).Eql("itworks")
			// d.Send("{}", []string{"q"})

			g.Assert(d.Error()).IsNil()

			defer func() {
				g.Assert(d.Close())
			}()
		})

		// g.It("dial read", func() {
		// 	d, err := Dial(clientConfig)
		// 	g.Assert(err).IsNil()
		// 	r, err := d.Read()
		// 	g.Assert(err).IsNil()
		// 	g.Assert(r).IsNotZero()
		// 	defer func() {
		// 		g.Assert(d.Close())
		// 	}()
		// })

	
	})
}
