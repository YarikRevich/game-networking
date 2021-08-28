package table

import "sync"

type Table struct {
	sync.Mutex
	table map[string]interface{}
}

func (t *Table) Add(addr string, data interface{}){
	t.Lock()
	defer t.Unlock()
	t.table[addr] = data
}

func (t *Table) Get(addr string)interface{}{
	return t.table[addr]
}

func New()*Table{
	return new(Table)
}