BINARY_NAME=dist/app

.PHONY: all
all:
	build test

.PHONY: build
build:
	go build -buildmode=plugin -o dist/p1.so plugins/plugin1/main.go
	go build -buildmode=plugin -o dist/p2.so plugins/plugin2/main.go
	go build -o ${BINARY_NAME} app/main.go

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
