package protocol

import "crypto/sha256"

type Protocol struct {
	HashSum [sha256.Size]byte `json:"hash_sum"`
	Procedure string `json:"procedure"`
	Msg interface{} `json:"msg"`
	Error error `json:"error"`
}

// type ProtocolResponse struct {
// 	Service
// 	Procedure string `json:"procedure"`
// 	Msg interface{} `json:"msg"`
// }

// import (
// 	"github.com/YarikRevich/game-networking/protocol/internal/fields"
// 	"github.com/YarikRevich/game-networking/protocol/internal/validator"
// )

// var (
// 	p Protocol
// )

// type Protocol interface {
// 	FieldManager() fields.FieldManager
// 	Validator() validator.Validator
// 	SetProtocol(interface{})
// }

// type P struct {
// 	fieldManager fields.FieldManager
// 	val          validator.Validator
// }

// func (q *P) FieldManager() fields.FieldManager {
// 	return q.fieldManager
// }

// func (q *P) Validator() validator.Validator {
// 	return q.val
// }

// func (q *P) SetProtocol(t interface{}){
// 	q.fieldManager.SetProtocolValues(t)
// 	q.val.SetProtocolTypes(t)
// }

// func UseProtocol(t interface{}) Protocol {
// 	if p == nil {
// 		p = &P{
// 			fieldManager: fields.UseFieldManager(),
// 			val: validator.UseValidator(),
// 		}
// 		p.SetProtocol(t)
// 	}
// 	return p
// }
