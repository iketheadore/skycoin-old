.DEFAULT_GOAL := help
.PHONY: link windows run run-help test lint check format install-linters release clean help

BINARY = skycoin
GOARCH = amd64

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=iketheadore
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT}"

# Build the project
all: link clean test vet linux darwin windows

link:
	BUILD_DIR=${BUILD_DIR}; \
	BUILD_DIR_LINK=${BUILD_DIR_LINK}; \
	CURRENT_DIR=${CURRENT_DIR}; \
	if [ "$${BUILD_DIR_LINK}" != "$${CURRENT_DIR}" ]; then \
	    echo "Fixing symlinks for build"; \
	    rm -f $${BUILD_DIR}; \
	    ln -s $${CURRENT_DIR} $${BUILD_DIR}; \
	fi

windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \
	cd - >/dev/null

# Static files directory
STATIC_DIR = src/gui/static

# Electron files directory
ELECTRON_DIR = electron

run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
	go run cmd/skycoin/skycoin.go --gui-dir="./${STATIC_DIR}" ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/skycoin/skycoin.go --help

test: ## Run tests
	go test ./cmd/... -timeout=1m
	go test ./src/... -timeout=1m

lint: ## Run linters. requires vendorcheck, gometalinter, golint, goimports
	gometalinter --disable-all -E goimports --tests --vendor ./...
	vendorcheck ./...

check: lint test ## Run tests and linters

install-linters: ## Install linters
	go get -u -f github.com/golang/lint/golint
	go get -u -f golang.org/x/tools/cmd/goimports
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/FiloSottile/vendorcheck

format:  # Formats the code. Must have goimports installed (use make install-linters).
	goimports -w ./cmd/...
	goimports -w ./src/...

release: ## Build electron apps, the builds are located in electron/release folder.
	cd $(ELECTRON_DIR) && ./build.sh
	@echo release files are in the folder of electron/release

clean: ## Clean dist files and delete all builds in electron/release
	rm $(ELECTRON_DIR)/release/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
