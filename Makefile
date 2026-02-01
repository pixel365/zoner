.PHONY: all tidy fa fmt lint selfcrt

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

selfcrt:
	openssl req -x509 -nodes -days 365 \
      -newkey rsa:2048 \
      -keyout server.key \
      -out server.crt \
      -config self-signed.dev.cnf

