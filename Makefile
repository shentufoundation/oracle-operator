PKG_LIST := $(shell go list ./...)
GOBIN ?= $(GOPATH)/bin
VERSION := $(shell echo $(shell git describe --tags --always) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=shentud \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=oracle-operator \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \

ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'

install: go.sum
	go install $(BUILD_FLAGS) .

build: go.sum
	go build $(BUILD_FLAGS) .

release: go.sum
	GOOS=linux go build $(BUILD_FLAGS) -o build/oracle-operator .
	GOOS=windows go build $(BUILD_FLAGS) -o build/oracle-operator.exe .
	GOOS=darwin go build $(BUILD_FLAGS) -o build/oracle-operator-macos .

tidy:
	@gofmt -s -w .
	@go mod tidy

lint: tidy
	@GO111MODULE=on go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout=10m

test: tidy
	@GO111MODULE=on go test ${PKG_LIST}

test-unit-cover: tidy
	@GO111MODULE=on go test ${PKG_LIST} -timeout=5m -tags='norace' -coverprofile=coverage.txt -covermode=atomic

all: install release lint
