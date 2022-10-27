export GOBIN := $(PWD)/bin
GO_PATH:=$(shell go env GOPATH)
export PATH:=$(GOBIN):$(GO_PATH):$(PATH)
GOMOD:=GO111MODULE=on
GO_TOOLS_VERSION=v0.1.10

install-go:
	goenv install -s $$(cat .go-version)

install-modules:
	go mod tidy

install-tools:
	mkdir -p bin; \
	go install golang.org/x/tools/cmd/goimports@$(GO_TOOLS_VERSION);
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.0;
	go install github.com/golang/mock/mockgen@v1.6.0;

install: install-go install-modules install-tools

lint:
	$(GOMOD) $(GOBIN)/golangci-lint run

gen-go:
	@$(GOMOD) go generate ./...

run-prod:
	go run app/cmd/main.go -debug=false

run:
	go run app/cmd/main.go -debug=true

test:
	go test ./...
build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/cmd app/cmd/main.go

deploy: build
	sls deploy --verbose

sls-offline: build
	sls offline
