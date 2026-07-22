.PHONY: run test lint fmt build migrate seed clean

run:
	go run cmd/server/main.go

test:
	go test ./... -v -cover

lint:
	golangci-lint run

fmt:
	gofmt -w .

build:
	go build -o bin/pos cmd/server/main.go

seed:
	go run cmd/seed/main.go

clean:
	rm -rf bin/ data/pos.db
