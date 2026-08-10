# Teltonika Go Parser

A lightweight Go library to decode and work with binary data from **Teltonika GPS devices**, including login and AVL data frames (Codecs 08, 8E, 12, 13, 14, 15 and 16) over TCP and UDP.
This version uses a clean, idiomatic Go project layout to separate concerns between command-line usage, internal logic, and reusable packages.

---

## 📦 Version

**v0.6.0**

---

## ✨ Features

- Decode and encode from the **public package** `github.com/danieljvsa/teltonika-go/public` (no `internal/` imports needed)
- Decode login packets
- Parse AVL records using Codecs 08, 8E, 16, 12, 13, 14, and 15 (AVL and command codecs)
- Encode AVL records and command responses for Codecs 08, 8E, 16, 12, 13, 14, and 15
- Strict Teltonika TCP/UDP framing validation (declared lengths, record counts, CRC)
- Decoding and encoding options (`WithStrictIMEI`, `WithLenientUDPLength`)
- Graceful error handling with structured responses
- Minimal dependencies, pure Go
- Comprehensive test coverage including a protocol-conformance suite and decoder fuzz tests

---

## 🆕 Changes Introduced

### v0.6.0
- 🏗️ **Public package API** - The `public` package exposes `Decode`/`Encode`, `Packet`, `AVLRecord`, `GPSData`, `IOElement`, `Command`, and `CommandResponse` types directly
- 🧬 **Backward compatibility preserved** - `pkg/` and `tools/` keep their existing APIs as thin shims over the public package, so existing code keeps compiling
- ✅ **External-consumer proof** - Public-package tests are written as an external test package and only import the public package
- 📝 **Updated README** with public-package examples and runnable `examples/` programs

### v0.5.0
- 🐛 **Fixed GPS speed decoding** - Corrected GPS block size from 14 to 15 bytes per Teltonika spec; speed is now parsed as uint16 instead of uint8
- 🧹 **Removed stray padding byte** - Eliminated erroneous `0x00` byte after GPS data in encoder that caused field misalignment
- 🆕 **Added Codec 12, 13, 14, 15 support** - Full support for command response codecs with command handling
- 🧹 **Production code cleanup** - Removed all debug print statements from decoder functions
- ✅ **Comprehensive test coverage** - Added extensive unit tests for all codec types and tool functions
- 🛡️ **Improved error handling** - Added bounds checking in header decoder to prevent panics on invalid data
- ⏰ **Enhanced timestamp support** - Added CalcTimestampSeconds and CalcTimestampSecondsBigEndian functions for 4-byte second timestamps
- 📦 **Better data structures** - Improved Record model with pointer fields for optional data support
- 🎧 **Added Codec 16 decoder support** - Support for GPRS/IQ frames with generation type
- 🧬 **Updated internal types** - Support for `generation_type` type workflows
- 📝 **Updated documentation** - Enhanced README with usage examples and project structure

---

## 🏗️ Project Structure

```
├── go.mod              # Go module file
├── LICENSE             # License (MIT)
├── Makefile            # Automation tasks
├── README.md           # Project documentation
├── public/             # Public package: the modern public API
│   ├── decode.go       # Decode, decode helpers
│   ├── doc.go          # Package documentation
│   ├── encode.go       # Encode, encode helpers
│   ├── gps.go          # GPS block encode/decode + CRC
│   ├── header.go       # header helpers
│   ├── models.go       # Public Packet/Record/Command types
│   └── teltonika_test.go
├── examples/           # Runnable example programs
│   ├── decode/
│   │   └── main.go     # Decode a Codec 08 frame
│   └── encode/
│       └── main.go     # Encode Codec 08 + login frames
├── cmd/
│   └── teltonika_go/
│       └── main.go     # CLI entry point
├── internal/           # Internal logic (not imported externally)
│   ├── decoder/
│   │   └── models.go   # Decoding-related structs
│   ├── encoder/
│   │   └── models.go   # Encoding logic structs
│   ├── header/
│   │   └── models.go   # AVL header model
│   ├── io/
│   │   └── models.go   # I/O element models
│   └── tool/
│       └── models.go   # Utility data types
├── pkg/                # Legacy API surface (shims over the public package)
│   ├── decoders.go
│   ├── encoders.go
│   ├── headers.go
│   └── ios.go
├── conformance/        # Protocol-conformance tests with official fixtures
│   ├── testdata/       # Golden wire frames (see testdata/README.md)
│   └── ...
├── test/               # Test suite
├── tools/              # Teltonika protocol utilities
│   ├── crc16.go
│   ├── gps.go
│   ├── login.go
│   ├── protocol.go
│   └── timestamp.go
```

---

## 🚀 Getting Started

### Requirements

- Go 1.22+ (see `go.mod`)
- Teltonika GPS device (e.g., FMB920, FMM125)

### Installation

```bash
go get github.com/danieljvsa/teltonika-go
```

---

## 📄 Example Usage

### Decode a full frame with the public package

```go
package main

import (
	"fmt"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func main() {
	// Raw Teltonika TCP login or data frame bytes (e.g. from a socket read).
	var frame []byte // = <device bytes>

	packet, err := teltonika.Decode(frame)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}

	switch packet.Kind {
	case teltonika.KindLogin:
		fmt.Printf("Login from IMEI %s\n", packet.IMEI)
	case teltonika.KindData:
		fmt.Printf("Codec 0x%02X, %d records\n", byte(packet.Codec), len(packet.Records))
		for _, rec := range packet.Records {
			fmt.Printf("  lat=%.6f lon=%.6f speed=%d\n", rec.GPS.Latitude, rec.GPS.Longitude, rec.GPS.Speed)
		}
	}
}
```

TCP does not preserve message boundaries: a single socket read may contain
part of a frame, one complete frame, or several frames. Callers must buffer
incoming bytes until a complete frame is available — read the eight-byte
header first, take the Data Field Length (bytes 4-8), and wait for
`8 + DataLength + 4` bytes before calling `Decode`.

### Encode a frame with the public package

```go
package main

import (
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func main() {
	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolTCP,
		Codec:    teltonika.Codec8,
		Records: []teltonika.AVLRecord{
			{
				Timestamp: time.Now().UTC(),
				Priority:  1,
				EventIO:   5,
				GPS: teltonika.GPSData{
					Latitude:   52.520008,
					Longitude:  13.404954,
					Altitude:   120,
					Angle:      25,
					Satellites: 7,
					Speed:      60,
				},
				IOElements: []teltonika.IOElement{{ID: 1, Value: "01"}},
			},
		},
	}

	frame, err := teltonika.Encode(packet)
	if err != nil {
		panic(err)
	}
	_ = frame // ready to send over the socket
}
```

The encoder writes the transport header, codec id, trailing record count and
(the TCP case) CRC for you. UDP frames additionally require a `Header.UDP`
with a valid `PacketID`, `AVLPacketID`, and optionally `IMEI`.

### Decoding and encoding options

`Decode` and `Encode` both accept options:

```go
packet, err := teltonika.Decode(frame, teltonika.WithLenientUDPLength())
```

- `WithLenientUDPLength()` - by default a UDP frame is rejected unless its
  declared length exactly matches the delivered datagram. Some devices and
  truncated official captures over-declare; this option accepts a declared
  length larger than the datagram (the declared length can never be smaller).

```go
frame, err := teltonika.Encode(loginPacket, teltonika.WithStrictIMEI())
```

- `WithStrictIMEI()` - by default login encoding only rejects empty or
  oversized identifiers, matching legacy behavior. With this option the IMEI
  must be exactly 15 decimal digits.

### Command codecs and Codec 15 raw response types

Command codecs 12-15 carry responses. Codec 15 frames may use a raw response
type byte (e.g. `0x0B`) that firmware emits; decoding preserves the byte as
its decimal string (`"11"`) and encoding maps it back, so these frames
round-trip byte-for-byte.

```go
package main

import (
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func main() {
	timestamp := time.Now().UTC()

	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolTCP,
		Codec:    teltonika.Codec15,
		Commands: []teltonika.Command{
			{
				Type: "11", // raw response type byte 0x0B
				Responses: []teltonika.CommandResponse{
					{
						Timestamp: &timestamp,
						IMEI:      "0123456789123456",
						Response:  "Hello!\n",
					},
				},
			},
		},
	}

	frame, err := teltonika.Encode(packet)
	if err != nil {
		panic(err)
	}
	_ = frame
}
```

### Legacy `pkg` API

The `pkg` package keeps its previous API and now delegates to the public package.

```go
package main

import (
	"fmt"
	pkg "github.com/danieljvsa/teltonika-go/pkg"
)

func main() {
	// Replace with actual Teltonika login and AVL packet bytes
	rawLogin := []byte{ /* login packet */ }
	rawFrame := []byte{ /* AVL frame */ }

	// Decode login packet
	login := pkg.LoginDecoder(rawLogin)
	if login.Error != nil {
		fmt.Println("Login decode error:", login.Error)
	} else {
		fmt.Printf("Login decoded: %+v\n", login.Response)
	}

	// Decode AVL/data frame
	frame := pkg.TramDecoder(rawFrame)
	if frame.Error != nil {
		fmt.Println("Frame decode error:", frame.Error)
	} else {
		fmt.Printf("Frame decoded: %+v\n", frame.Response)
	}
}
```

For new code, prefer the `public` package. The legacy `pkg` API decodes with
lenient UDP length handling to remain backward compatible.

---

## 🧪 Testing

```bash
go test ./...
go test -race ./...
go vet ./...
go test -run=^$ -fuzz=^FuzzDecodeNeverPanics$ -fuzztime=30s ./conformance
go test -run=^$ -fuzz=^FuzzDecodeLoginNeverPanics$ -fuzztime=30s ./conformance
go test -run=^$ -fuzz=^FuzzDecodeCodecDataNeverPanics$ -fuzztime=30s ./conformance
```

The `conformance/` suite decodes and re-encodes official Teltonika wire
fixtures stored in `conformance/testdata/` (see `conformance/testdata/README.md`
for provenance).

---

## 📄 License

[MIT License](LICENSE)

---

## 🤝 Contributing

Contributions, issues, and suggestions are welcome.  
Please fork the repo and submit a pull request or open an issue.

---

## 👤 Author

**Daniel Sá**  
[github.com/danieljvsa](https://github.com/danieljvsa)