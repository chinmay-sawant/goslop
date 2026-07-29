.PHONY: build test vet fmt lint

build:
	go build -o bin/codehound ./cmd/codehound

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

lint: vet
	gofmt -l . | tee /tmp/gofmt.out; test ! -s /tmp/gofmt.out
