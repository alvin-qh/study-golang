BINARY_NAME=out/logrus-demo

.PHONY: all
all:
	build test

build:
	go build -o ${BINARY_NAME} main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	go clean
	rm ${BINARY_NAME}
