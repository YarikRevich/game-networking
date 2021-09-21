package tests

import (
	"testing"

	"github.com/YarikRevich/game-networking/tools/pkg/creators"
	"github.com/franela/goblin"
)

func TestCreateAddr(t *testing.T){
	g := goblin.Goblin(t)
	g.Describe("Testing 'CreateAddr' function", func() {
		g.It("Test 'CreateAddr'(success)", func(){
			addr, err := creators.CreateAddr("127.0.0.1", "9999")
			g.Assert(err).IsNil()
			g.Assert(addr).Equal("127.0.0.1:9999")
		})
	})
}