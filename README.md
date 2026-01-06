# Teltonika Go Parser

A lightweight Go library to decode and work with binary data from **Teltonika GPS devices**, including login and AVL data packets (Codecs 08, 8E, etc.).
This version uses a clean, idiomatic Go project layout to separate concerns between command-line usage, internal logic, and reusable packages.

---

## 📦 Version

**v0.4.1**

---

## ✨ Features

- Decode login packets  
- Parse AVL records using Codecs 08, 8E, 16, 12, 13, 14, and 15
- Support for command response codecs with command handling
- Validate and interpret Teltonika TCP/UDP headers  
- Graceful error handling with structured responses  
- Minimal dependencies, pure Go
- Comprehensive test coverage with 30+ unit tests

---

## 🆕 Changes Introduced

### v0.4.0
- 🆕 **Added Codec 12, 13, 14, 15 support** - Full support for command response codecs with command handling
- 🧹 **Production code cleanup** - Removed all debug print statements from decoder functions
- ✅ **Comprehensive test coverage** - Added extensive unit tests for all codec types and tool functions
- 🛡️ **Improved error handling** - Added bounds checking in header decoder to prevent panics on invalid data
- ⏰ **Enhanced timestamp support** - Added CalcTimestampSeconds and CalcTimestampSecondsBigEndian functions for 4-byte second timestamps
- 📦 **Better data structures** - Improved Record model with pointer fields for optional data support

### Previous Changes
- 🎧 Added support for decoding with Codec 16  
- 🧬 Updated internal types to support `generation_type` type workflows

---

## 🏗️ Project Structure

```
├── go.mod              # Go module file
├── LICENSE             # License (MIT)
├── Makefile            # Automation tasks
├── README.md           # Project documentation
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
├── pkg/                # Public API surface
│   ├── decoders.go
│   ├── encoders.go
│   ├── headers.go
│   └── ios.go
├── test/               # Test suite
│   ├── decorders_test.go
│   ├── ios_test.go
│   ├── main_test.go
│   └── tools_test.go
└── tools/              # Teltonika protocol utilities
    ├── crc16.go
    ├── gps.go
    ├── login.go
    ├── protocol.go
    └── timestamp.go
```

---

## 🚀 Getting Started

### Requirements

- Go 1.20+
- Teltonika GPS device (e.g., FMB920, FMM125)

### Installation

```bash
go get github.com/danieljvsa/teltonika-go
```

---

## 📄 Example Usage

```go
package main

import (
	"fmt"
	pkg "github.com/danieljvsa/teltonika-go/pkg" // For general functions
	tools "github.com/danieljvsa/teltonika-go/tools" // For general functions
)

func main() {
	// Replace with actual Teltonika login and AVL packet bytes
	rawLogin := []byte{ /* login packet */ }
	rawTram := []byte{ /* AVL packet */ }

	// Decode login packet
	login := pkg.LoginDecoder(rawLogin)
	if login.Error != nil {
		fmt.Println("Login decode error:", login.Error)
	} else {
		fmt.Printf("Login decoded: %+v\n", login.Response)
	}

	// Decode AVL/tram packet
	tram := pkg.TramDecoder(rawTram)
	if tram.Error != nil {
		fmt.Println("Tram decode error:", tram.Error)
	} else {
		fmt.Printf("Tram decoded: %+v\n", tram.Response)
	}
}
```

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

