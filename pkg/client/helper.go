package client

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"reflect"
)

func GenerateHashSum(src interface{}) ([sha256.Size]byte, error) {
	var (
		b   []byte
		err error
	)
	if !reflect.ValueOf(src).IsValid() {
		rint, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return sha256.Sum256([]byte("")), err
		}
		b, err = json.Marshal(rint.String())
		if err != nil {
			return sha256.Sum256([]byte("")), err
		}

	} else {
		b, err = json.Marshal(src)
		if err != nil {
			return sha256.Sum256([]byte("")), err
		}
	}
	return sha256.Sum256(b), err
}
