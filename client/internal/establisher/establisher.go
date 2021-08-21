package establisher

import (
	"net"
	"time"

	"github.com/YarikRevich/game-networking/client/internal/timeout"
	"github.com/YarikRevich/game-networking/client/internal/workers"
)

type Connector struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	timeout *timeout.Timeout
	wmanager *workers.WorkerManager
}

func (c *Connector) EstablishConnection() error {
	conn, err := net.DialUDP("udp", nil, c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Connector) InitTimeouts() error{
	return c.timeout.EstimateProperTimout()
} 

func (c *Connector) InitWorkers(count int){
	c.wmanager = workers.New(count, c.GetConn())
}

func (c *Connector) SetReadDeadLine(){
	rt := c.timeout.GetReadTimeout()
	c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(rt)))
}

func (c *Connector) SetWriteDeadLine(){
	wt := c.timeout.GetWriteTimeout()
	c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(wt)))
}

func (c *Connector) Close()error{
	return c.conn.Close()
}

func (c *Connector) GetConn() *net.UDPConn{
	return c.conn
}

func (c *Connector) Ping() error{
	return c.wmanager.Ping()
}

func NewConnector(addr *net.UDPAddr, timeout *timeout.Timeout) *Connector {
	return &Connector{
		addr: addr,
		timeout: timeout,
	}
}
