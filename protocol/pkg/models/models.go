package models

import (
	"encoding/json"
)

var (
	msgFields = []string{
		"id",
		"ping",
		"data",
	}
)

type ID struct {
	ID string `json:"id"`
}

type Pinger struct {
	Ping bool `json:"ping"`
}

type Data struct {
	Data interface{} `json:"data"`
}

//It's a struct which will be given through
//the dialer and then should be processed by
//server
type Msg struct {
	ID
	Pinger
	Data
}

func IsProtocolMsg(msg []byte)bool {
	var stub map[string]json.RawMessage
	if err := json.Unmarshal(msg, &stub); err != nil{
		return false
	}

	for _, value := range msgFields{
		if _, ok := stub[value]; !ok{
			return false
		}
	}

	return true
}
