.PHONY: build clean

BINARY_NAME = wxccserver
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/$(BINARY_NAME) ./cmd
	cp ./cmd/config.toml ./build/config.toml
	cp ./cmd/accounts.toml ./build/accounts.toml

build-mac:
	cd cmd/
	go build -o ./build/$(BINARY_NAME) ./cmd
	cp ./cmd/config.toml ./build/config.toml
	cp ./cmd/accounts.toml ./build/accounts.toml

clean: 
	rm -rf build/

help:
	@echo "make build: 编译程序"
	@echo "make clean: 清空编译文件"