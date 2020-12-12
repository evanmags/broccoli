package main

import (
	"net/http"

	"github.com/evanmags/broccoli/pdk"
)

var Route = pdk.Route{
	Uri:  "/test/",
	Port: 8080,
}

func ApplyInboundPolicy(res http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
	return res, req
}

var Service = pdk.Upstream{
	Protocol: "https",
	Host:     "httpbin.org",
	Port:     443,
	Uri:      "/",
	Methods:  []string{"GET"},
}

func ApplyOutboundPolicy(res http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
	return res, req
}
