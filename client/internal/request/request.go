package request

import (
	// "github.com/YarikRevich/game-networking/client/tools"
	"encoding/json"
	"log"

	"github.com/YarikRevich/game-networking/protocol/pkg/models"
)

func NewRequest(procudure string, data interface{}) models.Msg {
	return models.Msg{Procedure: procudure, Data: data}
}

func CompleteRequestWithID(id int, msg *models.Msg) {
	msg.ID = id
}

func FormatRequestToJSON(msg models.Msg) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

// func CreateReq(c interface{}) []byte {

// 	// return ProcessReq(c, ping)
// }

// func ProcessReq(c interface{}, ping bool) []byte {

// 	id, err := tools.CreateUUID()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	b, err := json.Marshal({ID: id, Ping: ping, Data: c})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return b
// }
