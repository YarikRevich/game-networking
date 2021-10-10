package buffer

import (
	"testing"

	"github.com/franela/goblin"
)

func TestBuffer(t *testing.T){
	g := goblin.Goblin(t)
	g.Describe("Testing buffer ", func() {
		g.It("Test creation", func(){
			g.Assert(New()).IsNotNil()
		})

		b := New()
		const stub = "put_stub"

		g.It("Test basic work", func(){
			b.PutToBuffer(stub)
			g.Assert(b.GetFromBuffer().(string)).Equal(stub)
		})

		g.It("Test GetFromBuffer New pool", func(){
			g.Assert(cap(b.GetFromBuffer().([]byte))).Eql(32 * 1024)
		})
	})
}
