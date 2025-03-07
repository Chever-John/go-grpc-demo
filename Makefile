# 定义项目根目录变量
ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	@bash $(ROOT_DIR)/scripts/gen-proto.sh
	@echo "Done."

.PHONY: certs
certs:
	@echo "Generating certificates..."
	@bash $(ROOT_DIR)/scripts/gen-certs.sh
	@echo "Certificates generated."

.PHONY: build
build:
	@echo "Building server and client..."
	@go build -o bin/server ./cmd/server
	@go build -o bin/client ./cmd/client
	@echo "Build completed: binaries in ./bin/"

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -rf bin/
	@echo "Clean completed."

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make proto   - Generate Protocol Buffers code"
	@echo "  make certs   - Generate TLS certificates"
	@echo "  make build   - Build server and client binaries"
	@echo "  make clean   - Remove generated binaries"