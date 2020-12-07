BINARY_NAME=cwl
VERSION=dev

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" \
	-o bin/$(BINARY_NAME) -v .