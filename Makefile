build:
	go build

build_plugins:
	go build -buildmode plugin -o ./plugins/plugin.httpbin.so ./plugins/httpbin/plugin.go

run:
	go build -o ./bin && ./bin/broccoli