package tests

import (
	"testing"

	"github.com/YarikRevich/game-networking/client/internal/request"
	"github.com/YarikRevich/game-networking/client/pkg/config"
	"github.com/YarikRevich/game-networking/client/pkg/connector"
	"github.com/franela/goblin"
)

func TestConnect(t *testing.T){
	g := goblin.Goblin(t)
	g.Describe("Testing connect function", func(){
		g.It("Test connector", func(){
			conn, err := connector.Connect(config.Config{
				IP: "127.0.0.1",
				Port: "8080",
				PingerAddr: "http://google.com",
			})
			conn.InitTimeouts()
			conn.InitWorkers(4)

			conn.WorkerManager().Write(request.NewRequest(
				"ping", nil,
			))
			
			// m, err := conn.WorkerManager().Read()
			g.Assert(err).IsNil(err)

			

			defer conn.Close()
			g.Assert(err).IsNil(err)
		})
	})
}

func TestWorkers(t *testing.T){
	g := goblin.Goblin(t)
	g.Describe("Testing workers", func(){
		g.It("Test connector", func(){
			// conn := client.Connect(client.Config{
			// 	IP: "127.0.0.1",
			// 	Port: "8080",
			// 	PingerAddr: "http://google.com",
			// })
			// defer conn.Close()

			// workerManager := workers.New(4, conn.GetConn())
			// workerManager.Run()
		
			// data := workerManager.Read()
			// if err := workerManager.Error(); err != nil{
			// 	return
			// }
		})
	})
}