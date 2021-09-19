package tests

import (
	"testing"

	"github.com/YarikRevich/game-networking/client/pkg/config"
	"github.com/YarikRevich/game-networking/client/pkg/dialer"
	"github.com/YarikRevich/game-networking/tests/mocks"
	"github.com/golang/mock/gomock"

	// co "github.com/YarikRevich/game-networking/server/pkg/config"
	// c "github.com/YarikRevich/game-networking/server/pkg/connector"
	"github.com/YarikRevich/game-networking/server/pkg/handlers"
	"github.com/franela/goblin"
)

func TestClient(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("TestClient", func() {
		p := mocks.NewMockProtocolManager(gomock.NewController(t))

		g.It("TestDialer(Connection)", func() {
			conn, err := dialer.Dial(config.Config{
				IP:         "127.0.0.1",
				Port:       "8080",
				PingerAddr: "https://www.google.com/",
			}, p)
			g.Assert(err).IsNil()
			g.Assert(conn).IsNotNil()
			g.Assert(conn.Close()).IsNil()
		})

		g.It("TestDialer(Read)", func() {
			conn, err := dialer.Dial(config.Config{
				IP:         "127.0.0.1",
				Port:       "8080",
				PingerAddr: "https://www.google.com/",
			}, p)
			g.Assert(err).IsNil()
			g.Assert(conn).IsNotNil()

		

			defer func(){
				g.Assert(conn.Close()).IsNil()
			}()
		})
	})
}

// func TestInits(t *testing.T) {
// 	g := goblin.Goblin(t)
// 	g.Describe("Testing inits", func() {
// 		g.It("Test inits", func() {
// 			conn, err := connector.Connect(config.Config{
// 				IP:         "127.0.0.1",
// 				Port:       "8080",
// 				PingerAddr: "https://www.google.com/",
// 			})
// 			g.Assert(err).IsNil()

// 			err = conn.EstablishConnection()
// 			g.Assert(err).IsNil()

// 			defer func() {
// 				g.Assert(conn.Close()).IsNil()
// 			}()

// 			g.Assert(conn).IsNotNil()
// 			g.Assert(err).IsNil()

// 			err = conn.InitTimeouts()
// 			g.Assert(err).IsNil()

// 			conn.InitWorkers(4)
// 		})
// 	})
// }

// func TestWorkers(t *testing.T) {
// 	g := goblin.Goblin(t)
// 	g.Describe("Testing workers", func() {
// 		g.Before(func() {
// 			go func() {
// 				conn, err := c.Listen(
// 					co.Config{IP: "127.0.0.1", Port: "9090"})

// 				g.Assert(err).IsNil()
// 				g.Assert(conn).IsNotNil()

// 				err = conn.EstablishListening()
// 				g.Assert(err).IsNil()

// 				conn.InitWorkers(4)

// 				handlers.AddHandler("papa", func() []byte {
// 					return []byte(`{"id": 10, "procedure": "papa", "data": 20}`)
// 				})

// 				err = conn.WaitForInterrupt()
// 				g.Assert(err).IsNil()

// 			}()
// 		})

// 		g.It("Test workers", func() {
// 			conn, err := connector.Connect(config.Config{
// 				IP:         "127.0.0.1",
// 				Port:       "9090",
// 				PingerAddr: "www.google.com",
// 			})

// 			g.Assert(err).IsNil()

// 			err = conn.EstablishConnection()
// 			g.Assert(err).IsNil()

// 			defer func() {
// 				g.Assert(conn.Close()).IsNil()
// 			}()

// 			g.Assert(conn).IsNotNil()
// 			g.Assert(err).IsNil()

// 			err = conn.InitTimeouts()
// 			g.Assert(err).IsNil()

// 			conn.InitWorkers(4)

// 			conn.WorkerManager().Write(request.NewRequest(
// 				"papa", "10",
// 			))

// 			m, err := conn.WorkerManager().Read()
// 			g.Assert(err).IsNil()
// 			g.Assert(m).IsNotNil()
// 		})
// 	})
// }
