DIST_DIR=dist

.PHONY: all
all: test build

.PHONY: build
build:
	go build -ldflags="-s -w" -pgo=auto -o "${DIST_DIR}/viper"

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	go clean
	rm -rf "${DIST_DIR}"
