DIST_DIR=dist

.PHONY: build
build:
	mkdir -p "${DIST_DIR}"
	cd app && go build -ldflags="-s -w" -pgo=auto -o "../${DIST_DIR}/app"

.PHONY: run
run:
	cd app && go run main.go

.PHONY: clean
clean:
	go clean
	rm -rf "${DIST_DIR}"
