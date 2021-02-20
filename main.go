package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"plugin"

	"github.com/evanmags/broccoli/internals/api"

	"github.com/evanmags/broccoli/internals"
)

// The 'main' package of this program is just used to initiate the server
// this is what runs on 'startup' once the server is installed and running
// the daemon can be interfaced with using the `broc` command
func main() {
	log.Println("Loading Config")
	config, err := internals.LoadConfig("./broccoli.config.yaml")
	if err != nil {
		panic(err)
	}

	log.Println("Loading Plugins")
	mapping := loadPluginsToRoutes(config.PluginsDir)

	// This is the defualt handler, It allows us to handle relative proxied paths that
	// are absolute in the response... hacky and need to do research on how to improve.
	// TODO: handle recursive searches of referrer paths, this is bound to be a buggy
	//       piece of code to handle.
	log.Println("Loading default route '/'")
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		u, _ := url.Parse(req.Referer())

		if route, ok := mapping[u.Path]; ok {
			req.URL.Path = path.Join(u.Path, req.URL.Path)
			route.ServeHTTP(res, req)
		} else {
			http.NotFound(res, req)
		}
	})

	go api.InitApi(config.AdminPort, mapping)

	log.Fatalf("Server Error: %s", http.ListenAndServe(fmt.Sprintf(":%d", config.GatewayPort), nil))
}

// consumes a config pointer and returns a map of urls that point to the route
// object. This function also attaches handles to the server.
func loadPluginsToRoutes(pluginsDir string) map[string]*internals.Route {
	matches, err := filepath.Glob(pluginsDir + "*/*.so")
	if err != nil {
		panic(err)
	}

	serviceMap := map[string]*internals.Route{}

	for _, file := range matches {
		log.Printf("Found Plugin %s", file)
		plg, err := plugin.Open(file)
		if err != nil {
			panic(err)
		}
		route := internals.NewRouteFromPlugin(plg)

		log.Printf("Setting Route '%s' Proxy To '%s'", route.Uri, route.Upstream.Host)
		http.Handle(route.Uri, http.StripPrefix(route.Uri, route))
		serviceMap[route.Uri] = route
	}

	return serviceMap
}
