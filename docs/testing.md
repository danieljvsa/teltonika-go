# Testing

This section describes how to verify the library.

## Make targets

The Makefile provides these commands.

| Command | What it runs |
| --- | --- |
| make build | go build ./... |
| make test | go test ./... |
| make test-all | go test -v ./... |
| make vet | go vet ./... |
| make fmt | gofmt -w . |
| make fmt-check | fail when gofmt -l . reports a file |

Run `make fmt-check` before you commit. Keep the output empty.

## Race detection

Run the tests with the race detector:

```
go test -race ./...
```

## Fuzz the decoder

The conformance package has three fuzz targets. Run each one for at
least 30 seconds:

```
go test -run=^$ -fuzz=^FuzzDecodeNeverPanics$ -fuzztime=30s ./conformance
go test -run=^$ -fuzz=^FuzzDecodeLoginNeverPanics$ -fuzztime=30s ./conformance
go test -run=^$ -fuzz=^FuzzDecodeCodecDataNeverPanics$ -fuzztime=30s ./conformance
```

The fuzzers feed random bytes to the decoder. The decoder must never
panic.

## Conformance suite

The `conformance` package decodes and re-encodes official Teltonika
wire fixtures. The fixtures live in `conformance/testdata`. They come
from the official Teltonika documentation.

Most fixtures round-trip byte-for-byte. The Codec 16 UDP fixture is
truncated. It decodes only, and only with `WithLenientUDPLength()`.

See `conformance/testdata/README.md` for the provenance of every
fixture.

## Add a fixture

To add a conformance fixture:

1. Copy the hexadecimal frame into `conformance/testdata`.
2. Name the file with its codec and transport, for example
   `codec8e-tcp.hex`.
3. Add a row to the table in `conformance/testdata/README.md`.
4. Record where the frame came from and any transformation.
5. Write a test that decodes the frame and verifies the result.
