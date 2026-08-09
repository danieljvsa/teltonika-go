build:
	go build ./...

test-all:
	go test -v ./test

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l .