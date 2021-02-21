package internals

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"plugin"
	"strings"

	"github.com/evanmags/broccoli/pdk"
)

type Route struct {
	Uri            string
	Upstream       *url.URL
	Proxy          *httputil.ReverseProxy
	AllowedMethods []string
	Middleware     []pdk.Middleware
}

func NewRoute(
	listener pdk.Listener,
	dest pdk.Upstream,
	middleware []pdk.Middleware,
) *Route {
	upstreamURL := &url.URL{
		Scheme: dest.Protocol,
		Host:   fmt.Sprintf("%s:%d", dest.Host, dest.Port),
	}

	r := &Route{
		Uri:            listener.URI,
		Upstream:       upstreamURL,
		Proxy:          httputil.NewSingleHostReverseProxy(upstreamURL),
		AllowedMethods: listener.Methods,
		Middleware:     middleware,
	}

	origDirector := r.Proxy.Director
	r.Proxy.Director = func(req *http.Request) {
		origDirector(req)
		r.ensureProxyPathStripped(req)
		for _, mw := range middleware {
			mw.InboundAccess(req)
		}
	}

	r.Proxy.ModifyResponse = func(resp *http.Response) error {
		for _, mw := range middleware {
			if err := mw.OutboundAccess(resp); err != nil {
				return err
			}
		}
		return nil
	}

	return r
}

func NewRouteFromPlugin(p *plugin.Plugin) *Route {
	proxy := loadPluginMember(p)
	return NewRoute(proxy.Listener, proxy.Upstream, proxy.Middleware)
}

func loadPluginMember(p *plugin.Plugin) pdk.Proxy {
	member, err := p.Lookup("BuildProxy")
	if err != nil {
		panic(err)
	}
	buildProxy := member.(func() pdk.Proxy)

	return buildProxy()
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("Forwarding '%s' on '%s'", req.URL.Path, r.Uri)

	// method verifications
	if !r.MethodIsAllowed(req.Method) {
		log.Printf("Method '%s' Not Allowed On Route '%s'", req.Method, req.URL.Path)
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	for _, mw := range r.Middleware {
		if !mw.Enforcer(res, req) {
			log.Printf("Failed Enforcer Stage in Middleware '%s'", mw.Name)
			return
		}
	}

	r.Proxy.ServeHTTP(res, req)
}

func (r *Route) MethodIsAllowed(method string) bool {
	for _, allowed := range r.AllowedMethods {
		if allowed == method {
			log.Printf("Method '%s' Allowed On Route '%s'", method, r.Uri)
			return true
		}
	}

	return false
}

func (r *Route) ensureProxyPathStripped(req *http.Request) {
	for strings.HasPrefix(req.URL.Path, r.Uri) {
		log.Printf("Stripping Path '%s' from '%s'", r.Uri, req.URL.Path)
		req.URL.Path = strings.Replace(req.URL.Path, r.Uri, "/", 1)
	}
}
