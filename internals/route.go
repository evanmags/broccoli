package internals

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/evanmags/broccoli/pdk"
)

type Route struct {
	Uri      string
	Upstream *url.URL
	Proxy    *httputil.ReverseProxy
}

func NewRoute(uri string, dest pdk.Upstream) *Route {
	upstreamURL := &url.URL{
		Scheme: dest.Protocol,
		Host:   fmt.Sprintf("%s:%d", dest.Host, dest.Port),
	}

	return &Route{
		Uri:      uri,
		Upstream: upstreamURL,
		Proxy:    httputil.NewSingleHostReverseProxy(upstreamURL),
	}
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("forwarding on", r.Uri)
	req.URL.Path = strings.TrimPrefix(req.RequestURI, r.Uri)
	req.Host = r.Upstream.Host
	r.Proxy.ServeHTTP(res, req)
}
