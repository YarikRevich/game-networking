package client

import (
	"syscall"
	"testing"
	"time"

	"github.com/YarikRevich/game-networking/pkg/config"
	"github.com/YarikRevich/game-networking/pkg/server"
	"github.com/franela/goblin"
	"github.com/google/uuid"
)

type ResultStub struct {
	Result string
	Stack  struct {
		Key   string
		Value string
		ID    uuid.UUID
	}
}

func TestDialer(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("TestDialer", func() {
		go func() {
			c := server.Listen(config.Config{IP: "127.0.0.1", Port: "8090"})

			c.AddHandler("one_level_slice_1", func(m interface{}) (interface{}, error) {
				return []ResultStub{{Result: "1", Stack: struct{Key string; Value string; ID uuid.UUID}{Key: "Yarik", Value: "Svitlitsky", ID: uuid.New()}}}, nil
			})
			c.AddHandler("one_level_slice_2", func(m interface{}) (interface{}, error) {
				return []ResultStub{{Result: "2", Stack: struct{Key string; Value string; ID uuid.UUID}{Key: "Yarik", Value: "Svitlitskyi", ID: uuid.New()}}}, nil
			})
			c.AddHandler("one_level_slice_3", func(m interface{}) (interface{}, error) {
				return []ResultStub{{Result: "3", Stack: struct{Key string; Value string; ID uuid.UUID}{Key: "Yarik", Value: "Svitlitskyi", ID: uuid.New()}}}, nil
			})

			c.AddHandler("mult_level_slice_1", func(m interface{}) (interface{}, error) {
				return [][]ResultStub{{{Result: "1", Stack: struct{Key string; Value string; ID uuid.UUID}{Key: "Yarik", Value: "Svitlitskyi", ID: uuid.New()}}}}, nil
			})
			c.AddHandler("mult_level_slice_2", func(m interface{}) (interface{}, error) {
				return [][]ResultStub{{{Result: "2", Stack: struct{Key string; Value string; ID uuid.UUID}{Key: "Yarik", Value: "Svitlitskyi", ID: uuid.New()}}}}, nil
			})
			c.AddHandler("mult_level_slice_3", func(m interface{}) (interface{}, error) {
				return [][]ResultStub{{{Result: "3", Stack: struct{Key string; Value string; ID uuid.UUID}{Key: "Yarik", Value: "Svitlitskyi", ID: uuid.New()}}}}, nil
			})

			go func() {
				time.Sleep(15 * time.Second)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()
			g.Assert(c.WaitForInterrupt()).IsNil()
		}()

		time.Sleep(2 * time.Second)

		clientConfig := config.Config{
			IP:   "127.0.0.1",
			Port: "8090",
		}

		d, err := Dial(clientConfig)

		g.It("Check if connection is ok", func() {
			g.Assert(err).IsNil(err)
		})

		g.After(func() {
			time.Sleep(time.Second * 5)
			g.Assert(d.Close()).IsNil()
		})

		g.It("dial call, test one level slice dst", func() {
			var first []ResultStub
			d.Call("one_level_slice_1", nil, &first)

			var second []ResultStub
			d.Call("one_level_slice_2", nil, &second)

			var stub []ResultStub
			d.Call("one_level_slice_3", nil, &stub)
			d.Call("one_level_slice_3", nil, &stub)

			var third []ResultStub
			d.Call("one_level_slice_3", nil, &third)

			g.Assert(len(first)).Eql(1)
			g.Assert(len(second)).Eql(1)
			g.Assert(len(third)).Eql(1)
		})

		g.It("dial call, test multiple level slice dst", func() {
			var first [][]ResultStub
			d.Call("mult_level_slice_1", nil, &first)

			var second [][]ResultStub
			d.Call("mult_level_slice_2", nil, &second)

			var stub [][]ResultStub
			d.Call("mult_level_slice_3", nil, &stub)
			d.Call("mult_level_slice_3", nil, &stub)

			var third [][]ResultStub
			d.Call("mult_level_slice_3", nil, &third)

			g.Assert(len(first[0])).Eql(1)
			g.Assert(len(second[0])).Eql(1)
			g.Assert(len(third[0])).Eql(1)
		})
	})
}
