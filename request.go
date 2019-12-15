package coap

import "context"

type Request struct {
	Msg      Message
	Client   *ClientConn
	Ctx      context.Context
	Sequence uint64 // discontinuously growing number for every request from connection starts from 0
}

func (r *Request) Context() context.Context {
	if r.Ctx != nil {
		return r.Ctx
	}
	return context.Background()
}

func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	r2 := Request(*r)
	r2.Ctx = ctx
	return &r2
}
