build:
	go build ./...

test-all:
	go test -v ./...

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

ifeq ($(OS),Windows_NT)
fmt-check:
	@powershell -NoProfile -Command "$$f = gofmt -l .; if ($$f) { Write-Output $$f; exit 1 }"
else
fmt-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
endif
