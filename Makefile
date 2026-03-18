SHELL := /usr/bin/env bash

BINARY ?= keepass
PKG ?= ./...
GO ?= go
GOFLAGS ?= -trimpath
GOCACHE ?= $(CURDIR)/.cache/go-build
GOMODCACHE ?= $(CURDIR)/.cache/go-mod
GOPROXY ?= https://proxy.golang.org,direct
GOOS ?=
GOARCH ?=
CGO_ENABLED ?= 0
GOLANGCI_LINT ?= golangci-lint
VERSION ?= dev
COMMIT ?= unknown
BUILD_TIME ?= unknown
GO_LDFLAGS ?= -X github.com/photowey/keepass/internal/version.version=$(VERSION) -X github.com/photowey/keepass/internal/version.commit=$(COMMIT) -X github.com/photowey/keepass/internal/version.buildTime=$(BUILD_TIME)

.DEFAULT_GOAL := help

.PHONY: help clean deps download tidy verify-mod fmt fmt-check vet test test-race lint check build build-linux build-windows build-macos linux windows macos

GO_RUN = env GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)" GOPROXY="$(GOPROXY)" GOFLAGS="$(GOFLAGS)" $(GO)

## Show this help.
help:
	@{ \
		if [ -t 1 ] && command -v tput >/dev/null 2>&1 && [ "$$(tput colors 2>/dev/null || echo 0)" -ge 8 ]; then \
			C_CYAN="$$(tput setaf 6)"; C_BOLD="$$(tput bold)"; C_RESET="$$(tput sgr0)"; \
		else \
			C_CYAN=""; C_BOLD=""; C_RESET=""; \
		fi; \
		echo; \
		echo "Usage:"; \
		echo "  make <target> [VAR=value]"; \
		echo; \
		echo "Common targets:"; \
		awk -v cyan="$$C_CYAN" -v reset="$$C_RESET" 'BEGIN {FS=":.*##"} /^[a-zA-Z0-9_.-]+:.*##/ {printf "  %s%-18s%s %s\n", cyan, $$1, reset, $$2}' $(MAKEFILE_LIST) | sort; \
		echo; \
		echo "Variables (override with VAR=value):"; \
		printf "  %-18s %s\n" "GO" "$(GO)"; \
		printf "  %-18s %s\n" "PKG" "$(PKG)"; \
		printf "  %-18s %s\n" "BINARY" "$(BINARY)"; \
		printf "  %-18s %s\n" "GOCACHE" "$(GOCACHE)"; \
		printf "  %-18s %s\n" "GOMODCACHE" "$(GOMODCACHE)"; \
		printf "  %-18s %s\n" "VERSION" "$(VERSION)"; \
		printf "  %-18s %s\n" "COMMIT" "$(COMMIT)"; \
		echo; \
		echo "Examples:"; \
		echo "  make deps"; \
		echo "  make build"; \
		echo "  make check"; \
		echo "  make build GOOS=linux GOARCH=amd64"; \
		echo; \
	}

# ----------------------------------------------------------------

## Remove local build output.
clean:
	rm -rf "$(BINARY)"

## Download Go module dependencies.
deps: ## Alias: download
	$(GO_RUN) mod download

## Download Go module dependencies (compat alias).
download: deps ## Backward-compatible alias

## Tidy go.mod/go.sum.
tidy: ## go mod tidy
	$(GO_RUN) mod tidy

## Verify go.mod/go.sum are tidy and unchanged.
verify-mod: ## ensure go mod tidy makes no effective changes
	@tmpdir="$$(mktemp -d)"; \
	status=0; \
	cp go.mod "$$tmpdir/go.mod"; \
	cp go.sum "$$tmpdir/go.sum"; \
	if ! $(GO_RUN) mod tidy; then \
		status=$$?; \
	elif ! cmp -s "$$tmpdir/go.mod" go.mod || ! cmp -s "$$tmpdir/go.sum" go.sum; then \
		echo "go.mod/go.sum are not tidy. Run 'make tidy'."; \
		diff -u "$$tmpdir/go.mod" go.mod || true; \
		diff -u "$$tmpdir/go.sum" go.sum || true; \
		status=1; \
	fi; \
	cp "$$tmpdir/go.mod" go.mod; \
	cp "$$tmpdir/go.sum" go.sum; \
	rm -rf "$$tmpdir"; \
	exit $$status

## Format all Go packages.
fmt: ## go fmt ./...
	$(GO_RUN) fmt $(PKG)

## Check formatting without rewriting files.
fmt-check: ## gofmt -l on tracked Go files
	@files="$$(git ls-files '*.go')"; \
	if [ -z "$$files" ]; then \
		exit 0; \
	fi; \
	unformatted="$$(gofmt -l $$files)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

## Run go vet.
vet: ## go vet ./...
	$(GO_RUN) vet $(PKG)

## Run tests with coverage.
test: ## go test ./... -cover
	$(GO_RUN) test $(PKG) -cover -p 1

## Run race-enabled tests with shuffled order.
test-race: ## go test -race -shuffle=on ./...
	$(GO_RUN) test $(PKG) -race -shuffle=on -covermode=atomic -coverprofile=coverage.out

## Run the full local verification suite.
check: verify-mod fmt-check vet test-race ## local quality gate

## Run golangci-lint (install if needed).
lint: ## golangci-lint run
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found. Install: https://golangci-lint.run/docs/welcome/install/"; exit 1; }
	env GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)" GOPROXY="$(GOPROXY)" GOFLAGS="$(GOFLAGS)" $(GOLANGCI_LINT) run

## Build local binary for current OS/arch.
build: clean ## go build -o keepass
	env GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)" GOPROXY="$(GOPROXY)" GOFLAGS="$(GOFLAGS)" GOOS="$(GOOS)" GOARCH="$(GOARCH)" CGO_ENABLED="$(CGO_ENABLED)" $(GO) build -ldflags "$(GO_LDFLAGS)" -o "$(BINARY)" .

## Build Linux amd64 binary.
build-linux: clean ## GOOS=linux GOARCH=amd64
	env GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)" GOPROXY="$(GOPROXY)" GOFLAGS="$(GOFLAGS)" GOOS=linux GOARCH=amd64 CGO_ENABLED="$(CGO_ENABLED)" $(GO) build -ldflags "$(GO_LDFLAGS)" -o "$(BINARY)" .

## Build Windows amd64 binary.
build-windows: clean ## GOOS=windows GOARCH=amd64
	env GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)" GOPROXY="$(GOPROXY)" GOFLAGS="$(GOFLAGS)" GOOS=windows GOARCH=amd64 CGO_ENABLED="$(CGO_ENABLED)" $(GO) build -ldflags "$(GO_LDFLAGS)" -o "$(BINARY).exe" .

## Build macOS amd64 binary.
build-macos: clean ## GOOS=darwin GOARCH=amd64
	env GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)" GOPROXY="$(GOPROXY)" GOFLAGS="$(GOFLAGS)" GOOS=darwin GOARCH=amd64 CGO_ENABLED="$(CGO_ENABLED)" $(GO) build -ldflags "$(GO_LDFLAGS)" -o "$(BINARY)" .

## Backward-compatible aliases.
linux: build-linux ## Alias: build-linux
windows: build-windows ## Alias: build-windows
macos: build-macos ## Alias: build-macos

# ----------------------------------------------------------------

%:
	@echo "Unknown target: $@"
	@echo "Use 'make help' to see available targets."
