
run:
	go run ./cmd/web_console

generate: install_tools
	go generate ./...

install_tools:
	go get github.com/wzshiming/gen/cmd/...

fmt:
	go fmt ./...
