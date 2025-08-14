DIST_DIR=dist

.PHONY: all
all: test build

.PHONY: build
build:
	cd app && go build -ldflags="-s -w" -pgo=auto -o "../${DIST_DIR}/app"

.PHONY: test
test:
	cd app && go test -v ./...

.PHONY: run
run:
	cd app && go run main.go

.PHONY: clean
clean:
	go clean
	rm -rf "${DIST_DIR}"
