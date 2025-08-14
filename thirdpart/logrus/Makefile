DIST_DIR=dist
LOG_DIR=logging/logs

.PHONY: all
all: test build

.PHONY: build
build:
	go build -ldflags="-s -w" -pgo=auto -o "${DIST_DIR}/logrus"

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
	rm -rf "${LOG_DIR}"
