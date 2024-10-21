SERVER_DIR = ./cmd/server
SERVER_BIN_NAME = bamboo-apiserver
CLI_DIR = ./cmd/bamboofwcli
CLI_BIN_NAME = bbfw

.PHONY: all-platform
all-platform:
	build/build.sh $(SERVER_DIR) $(SERVER_BIN_NAME) all
	build/build.sh $(CLI_DIR) $(CLI_BIN_NAME) all

.PHONY: all
all: build-server build-bbfw

.PHONY: build-server
build-server:
	build/build.sh $(SERVER_DIR) $(SERVER_BIN_NAME)

.PHONY: build-bbfw
build-bbfw:
	build/build.sh $(CLI_DIR) $(CLI_BIN_NAME)

.PHONY: clean
clean:
	build/clean.sh
