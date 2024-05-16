.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	go clean
	rm -rf "$(dir ${BINARY_NAME})"
	rm -rf "./dist"
