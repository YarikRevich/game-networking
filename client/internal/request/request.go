package request

import (
	"encoding/json"
	"log"

	"github.com/YarikRevich/game-networking/client/tools"
)

type PostProcessor struct {
	ID   string
	Ping bool
	Data interface{}
}

type Request struct {}

func (r *Request) CreateReq(c interface{}, ping bool) []byte {
	return r.ProcessReq(c, ping)
}

func (r *Request) ProcessReq(c interface{}, ping bool) []byte {
	id, err := tools.CreateUUID()
	if err != nil {
		log.Fatalln(err)
	}
	b, err := json.Marshal(PostProcessor{ID: id, Ping: ping, Data: c})
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

func New() *Request {
	return new(Request)
}
