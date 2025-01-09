.phony: dist update tidy test-all lint fmt vet gosec

default: dist

dist: tidy test-all lint fmt vet gosec

update:
	go get -u ./...

tidy:
	go mod tidy

test-all:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

gosec:
	gosec ./...
