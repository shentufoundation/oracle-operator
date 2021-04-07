PKG_LIST := $(shell go list ./...)

export GO111MODULE = on

install: go.sum
	go install .

release: go.sum
	GOOS=linux go build -o build/oracle-operator .
	GOOS=windows go build -o build/oracle-operator.exe .
	GOOS=darwin go build -o build/oracle-operator-macos .

tidy:
	@gofmt -s -w .
	@go mod tidy

lint: tidy
	@GO111MODULE=on golangci-lint run --config .golangci.yml

test: tidy
	@GO111MODULE=on go test ${PKG_LIST}

all: install release lint
