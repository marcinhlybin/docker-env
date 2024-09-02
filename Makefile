PKG := github.com/marcinhlybin/docker-env
BUILD_DATE := $(shell date "+%Y-%m-%d %H:%M:%S")
COMMIT_HASH := $(shell git rev-parse HEAD)
OUTPUT := docker-env

all: build

test:
	@go test -v -cover ./...

build: clean
	@go build -ldflags "-X '$(PKG)/version.BuildDate=$(BUILD_DATE)' -X '$(PKG)/version.CommitHash=$(COMMIT_HASH)'" -o $(OUTPUT)

install: build
	sudo install -m 0755 $(OUTPUT) /usr/local/bin

run: build
	@./$(OUTPUT)

clean:
	@rm -f $(OUTPUT)

version: build
	@./$(OUTPUT) --version

.PHONY: all build clean version
