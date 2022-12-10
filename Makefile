.phony:
	default
	fmt
	build

default:
	make build

fmt:
	go fmt ./...

build:
	go build ./cmd/gvm
