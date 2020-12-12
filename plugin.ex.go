package main

import (
	"github.com/evanmags/broccoli/pdk"
)

var Route = pdk.Route{
	Uri:  "/test/",
	Port: 8080,
}

func ApplyInboundPolicy() {}

var Service = pdk.Upstream{
	Protocol: "https",
	Host:     "httpbin.org",
	Port:     443,
	Uri:      "/",
	Methods:  []string{"GET"},
}

func ApplyOutboundPolicy() {}
