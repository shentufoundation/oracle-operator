PKG_LIST := $(shell go list ./...)
GOBIN ?= $(GOPATH)/bin

export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=certik \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=oracle-operator \

ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'

install: go.sum
	go build $(BUILD_FLAGS) -o oracle-operator .
	mv oracle-operator $(GOBIN)

release: go.sum
	GOOS=linux go build $(BUILD_FLAGS) -o build/oracle-operator .
	GOOS=windows go build $(BUILD_FLAGS) -o build/oracle-operator.exe .
	GOOS=darwin go build $(BUILD_FLAGS) -o build/oracle-operator-macos .

tidy:
	@gofmt -s -w .
	@go mod tidy

lint: tidy
	@GO111MODULE=on golangci-lint run --config .golangci.yml

test: tidy
	@GO111MODULE=on go test ${PKG_LIST}

all: install release lint
