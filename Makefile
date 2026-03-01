.PHONY: all tidy fa fmt lint selfcrt test up down build

all: tidy fa fmt lint

tidy:
	go mod tidy

fa:
	@fieldalignment -fix ./...

fmt:
	@goimports -w -local github.com/pixel365/zoner .
	@gofmt -w .
	@golines -w .

lint:
	@golangci-lint run

test:
	@go test ./internal/... ./epp/... ./cmd/...

selfcrt:
	openssl req -x509 -nodes -days 365 \
      -newkey rsa:2048 \
      -keyout server.key \
      -out server.crt \
      -config self-signed.dev.cnf

integration:
	go test -v ./tests

up:
	@docker-compose -p zoner -f docker-compose.dev.yaml up -d

down:
	@docker-compose -p zoner -f docker-compose.dev.yaml down

build:
	@go build -ldflags "-s -w" -o ./build/migrate ./cmd/migrate

