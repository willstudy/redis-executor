GODIR = $(shell go list ./... | grep -v /vendor/)
PKG := github.com/willstudy/redis-executor
BUILD_IMAGE ?= golang:1.9.0-alpine
GOARCH := amd64
GOOS := linux
BUILD := $(shell git rev-parse HEAD)
LDFLAGS_REDIS_EXECUTOR := -ldflags "-X ${PKG}/cmd/redis-executor/app.Build=${BUILD}"

all: image

binaries: build-redis-executor
.PHONY: binaries

pre-build:
	@echo "pre build"
	@echo "clean all release files"
	@rm -rf ./go && rm -rf release
.PHONY: pre-build

build-dirs: pre-build
	@mkdir -p .go/src/$(PKG) ./go/bin
	@mkdir -p release
.PHONY: build-dirs

build-redis-executor: build-dirs
	@docker run                                                            \
	    --rm                                                               \
	    -ti                                                                \
	    -u $$(id -u):$$(id -g)                                             \
	    -v $$(pwd)/.go:/go                                                 \
	    -v $$(pwd):/go/src/$(PKG)                                          \
	    -v $$(pwd)/release:/go/bin                                         \
	    -e GOOS=$(GOOS)                                                    \
	    -e GOARCH=$(GOARCH)                                                \
	    -e CGO_ENABLED=0                                                   \
	    -w /go/src/$(PKG)                                                  \
	    $(BUILD_IMAGE)                                                     \
	    go install -v -pkgdir /go/pkg $(LDFLAGS_REDIS_EXECUTOR) ./cmd/redis-executor
.PHONY: build-redis-executor

lint:
	@golint $(GODIR)
.PHONY: lint

clean:
	@rm -rf ./release
.PHONY: clean
