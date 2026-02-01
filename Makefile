.PHONY: all tidy fa fmt lint

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
