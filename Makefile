build:
	go build -o bin/main ./cmd/teltonika_go/main.go

test-all:
	go test -v ./test