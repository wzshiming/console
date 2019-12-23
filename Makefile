
run:
	go run ./cmd/web_console

generate: install_tools
	go generate ./...

fmt:
	go fmt ./...
