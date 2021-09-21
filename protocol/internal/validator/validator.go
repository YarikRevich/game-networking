package validator

import (
	// "encoding/json"
	"reflect"
	// "github.com/YarikRevich/game-networking/common"
)

var (
	v Validator
)

// type Sync interface{
// 	Sync(common.Conn)
// }

type Validator interface {
	IsProtocolMsg(interface{}) bool
	IsProtocolSet() bool
	SetProtocolTypes(interface{})
}

type V struct {
	protocolTypes map[string]string
}

func (v *V) IsProtocolMsg(m interface{}) bool {
	var val reflect.Type

	if reflect.ValueOf(m).Kind() == reflect.Ptr{
		val = reflect.TypeOf(m).Elem()
	}else{
		val = reflect.TypeOf(m)
	}

	if val.NumField() != len(v.protocolTypes){
		return false
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(0)
		if t, ok := v.protocolTypes[f.Name]; !ok || f.Type.Name() != t {
			return false
		}
	}

	return true
}

func (p *V) IsProtocolSet() bool {
	return len(p.protocolTypes) != 0
}

func (p *V)SetProtocolTypes(t interface{}) {
	val := reflect.TypeOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		p.protocolTypes[f.Name] = f.Type.Name()
	}
}

// func (p *V) Marshal()[]byte{
	
//  }

// func (p *V) AddField(f reflect.Kind){
	
// }

//It's a struct which will be given through
//the dialer and then should be processed by
//server
// type Msg struct {
// 	ID        int         `json:"id"`
// 	Procedure string      `json:"procedure"`
// 	Data      interface{} `json:"data,omitempty"`
// }

// func UnmarshalProtocol(msg []byte) (Msg, error) {
// 	var umsg Msg
// 	return umsg, json.Unmarshal(msg, &umsg)
// }

func UseValidator() Validator {
	if v == nil {
		v = &V{
			protocolTypes: make(map[string]string),
		}
	}
	return v
}
