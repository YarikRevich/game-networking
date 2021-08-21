package states

const (
	RECEIVE = iota
	SEND
	PING
)

type State struct {
	curr int
}

func (s *State) GetCurrState() int {
	return s.curr
}

func (s *State) SetCurrState(ns int) {
	s.curr = ns
}

func New() *State {
	return new(State)
}
