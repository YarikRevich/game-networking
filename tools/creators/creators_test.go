package creators

import (
	"testing"

	"github.com/franela/goblin"
)

func TestCreateAddr(t *testing.T){
	g := goblin.Goblin(t)
	g.Describe("Testing 'CreateAddr' function", func() {
		g.It("Test 'CreateAddr'(success)", func(){
			addr, err := CreateAddr("127.0.0.1", "9999")
			g.Assert(err).IsNil()
			g.Assert(addr).Equal("127.0.0.1:9999")
		})
	})
}