package timeout

import (
	"reflect"
	"github.com/go-ping/ping"
)

type Timeout struct {
	addr string //Addr for pinger
	read int
	write int
}

func (t *Timeout) getPingStatistics()(*ping.Statistics, error){
	p, err := ping.NewPinger(t.addr)
	if err != nil{
		return nil, err
	}
	p.Count = 5
	if err := p.Run(); err != nil{
		return nil, err
	}
	return p.Statistics(), nil
}

func (t *Timeout)EstimateProperTimout()error{
	s, err := t.getPingStatistics()
	if err != nil{
		return err
	}

	t.read = int(s.AvgRtt * 20)
	t.write = int(s.AvgRtt * 30)
	return nil
}

func (t *Timeout) GetReadTimeout()int{
	return t.read
}

func (t *Timeout) GetWriteTimeout()int{
	return t.write
}

func NewTimeout(addr string)*Timeout{
	var a string = addr
	if reflect.ValueOf(addr).IsZero(){
		a = "www.google.com"
	}
	return &Timeout{addr: a}
}