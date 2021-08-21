package establisher

import (
	"net"
)

type Establisher struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func (e *Establisher) EstablishListening()error{
	conn, err := net.ListenUDP("udp", e.addr)
	if err != nil{
		return err
	}
	e.conn= conn
	return nil
}

func (e *Establisher) Close()error{
	return e.conn.Close()
}

func New()*Establisher{
	return new(Establisher)
}