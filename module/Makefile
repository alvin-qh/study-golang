TARGET_DIR=dist

.PHONY: build
build:
	mkdir -p dist
	cd app && go get -u ./...
	cd app && go build -o app
	mv app/app dist

.PHONY: clean
clean:
	go clean
	rm -rf "${TARGET_DIR}"
