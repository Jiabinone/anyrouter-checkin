.PHONY: dev build fmt lint clean gen-docs

dev:
	@air

build:
	@go build -o bin/server cmd/server/main.go

fmt:
	@go fmt ./...
	@gofumpt -l -w . 2>/dev/null || true

lint:
	@golangci-lint run ./... 2>/dev/null || go vet ./...

gen-docs:
	@swag init -g main.go -o docs --parseInternal --dir ./cmd/server,./internal/handler,./internal/service,./internal/model,./pkg/response

clean:
	@rm -rf bin/ tmp/ data/
