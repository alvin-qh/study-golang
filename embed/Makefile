BINARY_NAME=out/embed-demo

.PHONY: all
all: build test

.PHONY: build
build:
	go build -ldflags="-s -w" -pgo=auto -o ${BINARY_NAME} main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	go clean
	rm -rf "$(dir ${BINARY_NAME})"
	rm -rf "./dist"
