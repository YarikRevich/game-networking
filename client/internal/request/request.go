package request

type Request struct {
	req interface{}
}

func (r *Request) CreateReq(c interface{}){
	r.req = c
}

func (r *Request) ProcessReq(c interface{}){}


func New()*Request{
	return new(Request)
}