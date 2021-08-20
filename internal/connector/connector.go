package connector

import (
	"net"
	"github.com/YarikRevich/game-networking/internal/timeout"
)

type Connector struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	timeout *timeout.Timeout
}

func (c *Connector) EstablishConnection() error {
	conn, err := net.DialUDP("udp", nil, c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func NewConnector(addr *net.UDPAddr) *Connector {
	return &Connector{
		addr: addr,
		timeout: timeout.NewTimeout("www.google.com"),
	}
}
