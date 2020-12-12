package main

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"plugin"

	"github.com/evanmags/broccoli/pdk"

	"github.com/evanmags/broccoli/internals"
)

// The 'main' package of this program is just used to initiate the server
// this is what runs on 'startup' once the server is installed and running
// the daemon can be interfaced with using the `broc` command

func main() {
	config, err := internals.LoadConfig("./broccoli.config.yaml")
	if err != nil {
		panic(err)
	}

	plugins := loadPlugins(config)
	mapping := setPluginsAsHandlers(plugins)

	// This is the defualt handler, It allows us to handle relative proxied paths that
	// are absolute in the response... hacky and need to do research on how to improve.
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		u, _ := url.Parse(req.Referer())

		if _, ok := mapping[u.Path]; ok {
			req.URL.Path = path.Join(u.Path, req.URL.Path)
			mapping[u.Path].ServeHTTP(res, req)
		} else {
			http.NotFound(res, req)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func loadPlugins(config *internals.Config) []*plugin.Plugin {
	matches, err := filepath.Glob(config.PluginsDir + "*/*.so")
	if err != nil {
		panic(err)
	}

	plugins := []*plugin.Plugin{}

	for _, file := range matches {
		plg, err := plugin.Open(file)
		if err != nil {
			panic(err)
		}
		plugins = append(plugins, plg)
	}

	return plugins
}

func setPluginsAsHandlers(plugins []*plugin.Plugin) map[string]*internals.Route {
	serviceMap := map[string]*internals.Route{}

	for _, p := range plugins {
		r, _ := p.Lookup("Route")
		plgRoute := r.(*pdk.Route)

		s, _ := p.Lookup("Service")
		plgService := s.(*pdk.Upstream)

		route := internals.NewRoute(plgRoute.Uri, *plgService)

		http.Handle(route.Uri, route)
		serviceMap[route.Uri] = route
	}

	return serviceMap
}
