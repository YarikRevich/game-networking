package config

import "github.com/YarikRevich/wrapper/pkg/wrapper"


type Config struct {
	IP string

	Port string

	Encoder wrapper.WrapperEncoder
	Decoder wrapper.WrapperDecoder
}