package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/evanmags/broccoli/internals"
)

type JsonRoute struct {
	Upstream       *url.URL
	AllowedMethods []string
}

func InitApi(port int, serviceMap map[string]*internals.Route) {
	serviceJson := map[string]interface{}{}

	for k, v := range serviceMap {
		serviceJson[k] = JsonRoute{
			Upstream:       v.Upstream,
			AllowedMethods: v.AllowedMethods,
		}
	}

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		jsonMap, err := json.Marshal(serviceJson)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Write(jsonMap)
		}
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
