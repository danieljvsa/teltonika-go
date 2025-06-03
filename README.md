# Teltonika Go Parser

A lightweight Go library to decode and work with binary data from **Teltonika GPS devices**, including login and AVL data packets (Codecs 08, 8E, etc.).

---

## 📦 Version

**v0.1.0**

---

## ✨ Features

- Decode login packets  
- Parse AVL records using Codecs 08 and 8E  
- Validate and interpret Teltonika TCP/UDP headers  
- Graceful error handling with structured responses  
- Minimal dependencies, pure Go

---

## 🏗️ Project Structure

```
├── main.go              # Core decoder logic (entry point)
├── decoders.go          # Codec 08/8E parsing
├── encoders.go          # Placeholder for encoding support
├── ios.go               # I/O element parsing
├── router.go            # TCP connection routing
├── tools.go             # Utility functions
├── *_test.go            # Unit tests
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
	tg "github.com/danieljvsa/teltonika-go"
)

func main() {
	// Replace with actual Teltonika login and AVL packet bytes
	rawLogin := []byte{ /* login packet */ }
	rawTram := []byte{ /* AVL packet */ }

	// Decode login packet
	login := tg.LoginDecoder(rawLogin)
	if login.Error != nil {
		fmt.Println("Login decode error:", login.Error)
	} else {
		fmt.Printf("Login decoded: %+v\n", login.Response)
	}

	// Decode AVL/tram packet
	tram := tg.TramDecoder(rawTram)
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
