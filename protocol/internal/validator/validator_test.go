package validator

import (
	"testing"
	"github.com/franela/goblin"
)

func TestValidator(t *testing.T){
	g := goblin.Goblin(t)

	g.Describe("TestProtocol", func ()  {
		type Stub struct {
			Name string
			Surname string
		}

		g.It("TestValidatorSet", func(){
			val := UseValidator()
			val.SetProtocol(new(Stub))

			g.Assert(val.IsProtocolSet()).IsTrue()
		})

		g.It("TestValidatorChecker(Success)", func(){
			val := UseValidator()
			val.SetProtocol(new(Stub))

			g.Assert(val.IsProtocolMsg(Stub{Name: "yarik", Surname: "Svit"})).IsTrue()
		})


		g.It("TestValidatorChecker(False)", func(){
			val := UseValidator()
			val.SetProtocol(new(Stub))

			g.Assert(val.IsProtocolMsg(struct{Name string; Surname string; Age int}{Name: "yarik", Surname: "Svit", Age: 10})).IsFalse()
		})
	})
}