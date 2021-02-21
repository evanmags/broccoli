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

func (p *Proxy) AddMiddleware(mw Middleware) {
	p.Middleware = append(p.Middleware, mw)
}

type Listener struct {
	Port    int
	URI     string
	Methods []string // http.Method*
}

type Middleware struct {
	Name           string
	InboundAccess  func(*http.Request)
	Enforcer       func(http.ResponseWriter, *http.Request) bool
	OutboundAccess func(*http.Response) error
}

func NewMiddleware(name string) Middleware {
	newMW := Middleware{
		Name:           name,
		InboundAccess:  func(r *http.Request) {},
		Enforcer:       func(http.ResponseWriter, *http.Request) bool { return true },
		OutboundAccess: func(r *http.Response) error { return nil },
	}
	return newMW
}

type Upstream struct {
	Protocol string
	Host     string
	Port     int
	Uri      string
}
