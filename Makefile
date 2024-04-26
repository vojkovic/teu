.PHONY: clean checks test build image
$(info $(SHELL))
export GO111MODULE=on
export CGO_ENABLED=0

MAIN_DIRECTORY := ./cmd/teu/

BIN_OUTPUT := ./teu

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

default: clean build

clean:
	@echo BIN_OUTPUT: ${BIN_OUTPUT}
	rm -rf dist/ builds/ ./teu

build: clean
	@echo Version: $(VERSION)
	CGO_CFLAGS="-Wno-deprecated-declarations" GOARCH=arm64 CGO_ENABLED=1 go build -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIRECTORY}

build-linux-amd64: clean
	@echo Version: $(VERSION)
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIRECTORY}

test:
	CGO_CFLAGS="-Wno-deprecated-declarations" GOARCH=arm64 CGO_ENABLED=1 go test -v -cover ./cmd/teu

checks:
	golangci-lint run

