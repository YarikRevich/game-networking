package tests

import (
	"testing"

	"github.com/franela/goblin"
)

func TestNetwork(t *testing.T){
	g := goblin.Goblin(t)
	g.Describe("Testing network module", func(){
		g.It("Test connector", func(){
			client.Connect(client.Config{
				IP: "127.0.0.1",
				Port: "8080",
				PingerAddr: "http://google.com",
			})
		})
	})
}