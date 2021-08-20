package network

import (
	"net"
	"github.com/YarikRevich/game-client-networking/internal/connector"
)

func CreateConnection(addr *net.UDPAddr)*connector.Connector{
	return connector.NewConnector(addr);
}