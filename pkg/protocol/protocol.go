package protocol

import "crypto/sha256"

type Protocol struct {
	Msg       interface{}       `json:"msg"`
	HashSum   [sha256.Size]byte `json:"hash_sum"`
	Procedure string            `json:"procedure"`
}
