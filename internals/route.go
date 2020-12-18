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
	Uri            string
	Upstream       *url.URL
	Proxy          *httputil.ReverseProxy
	AllowedMethods []string
	InboundPolicy  func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)
	// OutboundPolicy func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)
}

func NewRoute(
	listener *pdk.Route,
	dest *pdk.Upstream,
	inboundPolicy func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request),
) *Route {
	upstreamURL := &url.URL{
		Scheme: dest.Protocol,
		Host:   fmt.Sprintf("%s:%d", dest.Host, dest.Port),
	}

	return &Route{
		Uri:            listener.Uri,
		Upstream:       upstreamURL,
		Proxy:          httputil.NewSingleHostReverseProxy(upstreamURL),
		AllowedMethods: listener.Methods,
		InboundPolicy:  inboundPolicy,
	}
}

func NewRouteFromPlugin(p *plugin.Plugin) *Route {
	r := loadPluginMember(p, "Route").(*pdk.Route)
	s := loadPluginMember(p, "Service").(*pdk.Upstream)
	ibp := loadPluginMember(p, "ApplyInboundPolicy").(func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request))

	return NewRoute(r, s, ibp)
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

	// method verifications
	if !r.MethodIsAllowed(req.Method) {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	res, req = r.InboundPolicy(res, req)
	print(req.RequestURI)
	clean := strings.Split(req.RequestURI, "?")
	print("clean", clean)
	req.URL.Path = strings.TrimPrefix(clean[0], r.Uri)
	println(req.URL.Query().Encode())
	req.RequestURI = req.URL.Path
	req.Host = r.Upstream.Host
	r.Proxy.ServeHTTP(res, req)
}

func (r *Route) MethodIsAllowed(method string) bool {
	for _, allowed := range r.AllowedMethods {
		if allowed == method {
			return true
		}
	}

	return false
}
