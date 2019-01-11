build: deps
	env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/main main.go

test:
	go test ./...

lint: ## Lint the files
	${GOPATH}/bin/golint ./... | grep -v "should have comment" || true

fmt: ## Format
	go fmt ./...

deps: ## Get the dependencies
	@go get -v -d ./...
	@go get -u github.com/golang/lint/golint

all: build test lint fmt
