DIST_DIR=dist

.PHONY: all
build: build-plugin1 build-plugin2 build-app

.PHONY: build-plugin1
build-plugin1:
	cd ./plugins/plugin1 && go build -ldflags="-s -w" -pgo=auto -buildmode=plugin -o "../../${DIST_DIR}/p1.so"

.PHONY: build-plugin2
build-plugin2:
	cd ./plugins/plugin2 && go build -ldflags="-s -w" -pgo=auto -buildmode=plugin -o "../../${DIST_DIR}/p2.so"

.PHONY: build-app
build-app: test
	cd ./app && go build -ldflags="-s -w" -pgo=auto -o "../${DIST_DIR}/app"

.PHONY: test
test: build-plugin1 build-plugin2
	cd app && go test -v ./...

.PHONY: run
run: build-plugin1 build-plugin2
	cd app && go run main.go

.PHONY: clean
clean:
	go clean
	rm -rf "${DIST_DIR}"
