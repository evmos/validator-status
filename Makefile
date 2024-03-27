.PHONY: build

run:
	@go run cmd/status/main.go

generate:
	@sqlc generate

install-deps:
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

lint:
	golangci-lint run --fix --out-format=line-number --issues-exit-code=0 --config .golangci.yml --color always ./...
