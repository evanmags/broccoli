package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evanmags/broccoli/internals"
)

type JsonRoute struct {
	Upstream       string
	Middleware     []string
	AllowedMethods []string
}

func InitApi(port int, serviceMap map[string]*internals.Route) {
	serviceJson := map[string]interface{}{}

	for k, v := range serviceMap {
		mwares := []string{}
		println(mwares)
		for _, mw := range v.Middleware {
			mwares = append(mwares, mw.Name)
			println(mwares)
		}
		serviceJson[k] = JsonRoute{
			Upstream:       fmt.Sprintf("%s://%s/%s", v.Upstream.Scheme, v.Upstream.Host, v.Upstream.Path),
			Middleware:     mwares,
			AllowedMethods: v.AllowedMethods,
		}

	}

	jsonMap, err := json.Marshal(serviceJson)
	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Write(jsonMap)
		}
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
