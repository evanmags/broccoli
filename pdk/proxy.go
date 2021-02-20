package pdk

import "net/http"

type Proxy struct {
	Listener   Listener
	Middleware []Middleware
	Upstream   Upstream
}

func NewProxy(listener Listener, upstream Upstream) *Proxy {
	return &Proxy{
		Listener: listener,
		Upstream: upstream,
	}
}

func (p *Proxy) AddMiddleware(mw IMiddleware, priorityIn, priorityOut int) {
	newMW := Middleware{
		InboundAccess:  mw.InboundAccess,
		OutboundAccess: mw.OutboundAccess,
	}
	newMW.Priority.In = priorityIn
	newMW.Priority.Out = priorityIn
	p.Middleware = append(p.Middleware, newMW)
}

type Listener struct {
	Port    int
	URI     string
	Methods []string // http.Method*
}

type Middleware struct {
	InboundAccess  func(*http.Request)
	OutboundAccess func(*http.Response) error
	Priority       struct {
		In  int
		Out int
	}
}

type IMiddleware interface {
	InboundAccess(*http.Request)
	OutboundAccess(*http.Response) error
}

func NewMiddleware(priorityIn, priorityOut int) *Middleware {
	newMW := Middleware{}
	newMW.Priority.In = priorityIn
	newMW.Priority.Out = priorityIn
	return &newMW
}

type Upstream struct {
	Protocol string
	Host     string
	Port     int
	Uri      string
}
