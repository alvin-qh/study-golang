TARGET_DIR=dist

.PHONY: build-plugin1
build-plugin1:
	cd ./plugins/plugin1 && go build -buildmode=plugin -o ../../${TARGET_DIR}/p1.so

.PHONY: build-plugin2
build-plugin2:
	cd ./plugins/plugin2 && go build -buildmode=plugin -o ../../${TARGET_DIR}/p2.so

.PHONY: build-app
build-app:
	cd ./app && go build -o ../dist/app

.PHONY: build
build: build-plugin1 build-plugin2 build-app

.PHONY: clean
clean:
	go clean
	rm -rf "${TARGET_DIR}"
