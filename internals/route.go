package internals

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"plugin"
	"strings"

	"github.com/evanmags/broccoli/pdk"
)

type Route struct {
	Uri           string
	Upstream      *url.URL
	Proxy         *httputil.ReverseProxy
	InboundPolicy func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)
	// OutboundPolicy func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)
}

func NewRoute(uri string,
	dest pdk.Upstream,
	inboundPolicy func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request),
) *Route {
	upstreamURL := &url.URL{
		Scheme: dest.Protocol,
		Host:   fmt.Sprintf("%s:%d", dest.Host, dest.Port),
	}

	return &Route{
		Uri:           uri,
		Upstream:      upstreamURL,
		Proxy:         httputil.NewSingleHostReverseProxy(upstreamURL),
		InboundPolicy: inboundPolicy,
	}
}

func NewRouteFromPlugin(p *plugin.Plugin) *Route {
	r := loadPluginMember(p, "Route").(*pdk.Route)
	s := loadPluginMember(p, "Service").(*pdk.Upstream)
	ibp := loadPluginMember(p, "ApplyInboundPolicy").(func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request))

	return NewRoute(r.Uri, *s, ibp)
}

func loadPluginMember(p *plugin.Plugin, name string) plugin.Symbol {
	member, err := p.Lookup(name)
	if err != nil {
		panic(err)
	}

	return member
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("forwarding on", r.Uri)

	res, req = r.InboundPolicy(res, req)

	req.URL.Path = strings.TrimPrefix(req.RequestURI, r.Uri)
	// req.Host = r.Upstream.Host
	r.Proxy.ServeHTTP(res, req)
}
