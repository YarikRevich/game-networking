package models

import (
	"encoding/json"
)

var (
	msgFields = []string{
		"id",
		"procedure",
		"data",
	}
)

//It's a struct which will be given through
//the dialer and then should be processed by
//server
type Msg struct {
	ID        int      `json:"id"`
	Procedure string      `json:"procedure"`
	Data      interface{} `json:"data,omitempty"`
}

func IsProtocolMsg(msg []byte) bool {
	var stub map[string]json.RawMessage
	if err := json.Unmarshal(msg, &stub); err != nil {
		return false
	}

	for _, value := range msgFields {
		if _, ok := stub[value]; !ok {
			return false
		}
	}

	return true
}

func UnmarshalProtocol(msg []byte) (Msg, error) {
	var umsg Msg
	return umsg, json.Unmarshal(msg, &umsg)
}
