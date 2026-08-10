# Using the API

This section describes how to decode and encode with options.

## Decode

Call `Decode` with a complete frame.

```go
packet, err := teltonika.Decode(frame)
```

The function detects the frame type. A login frame becomes a login
packet. A data frame becomes a data packet.

For a payload that you split yourself, use `DecodeCodecData`. It takes
the codec payload, the codec identifier, and the protocol.

Use `DecodeCodecIO` to decode a standalone I/O element block.

## Receive a frame over TCP

TCP does not preserve message boundaries. Buffer the incoming bytes
until a complete frame is available.

1. Read the eight-byte header.
2. Read the Data Field Length from header bytes `[4:8]`.
3. Calculate the complete frame size as `8 + Data Field Length + 4`.
4. Wait until that number of bytes is available.
5. Call `Decode` with exactly one complete frame.

The four bytes after the Data Field are the CRC. The trailing record
count is already included in the Data Field Length.

## Decode options

`Decode` accepts options. Pass none for the default behavior.

- `WithLenientUDPLength()` accepts a UDP frame whose declared length is
  larger than the delivered frame. Some device firmware over-declares
  the length. Truncated official captures do the same. The default
  requires an exact match. A declared length smaller than the delivered
  frame is always corrupt.

```go
packet, err := teltonika.Decode(frame, teltonika.WithLenientUDPLength())
```

## Encode

Call `Encode` with a packet.

```go
frame, err := teltonika.Encode(packet)
```

The encoder writes the transport header, the codec id, the trailing
count, and the CRC for TCP. It validates protocol fields that must fit
fixed wire widths and returns an error instead of silently truncating
out-of-range values.

A UDP packet requires a valid `Header.UDP` with a packet ID and an AVL
packet ID. The IMEI is optional. Set the header on the packet before
you call `Encode`.

## Encode options

`Encode` accepts options.

- `WithStrictIMEI()` requires an IMEI of exactly 15 decimal digits. The
  default rejects empty and oversized identifiers only.

```go
frame, err := teltonika.Encode(login, teltonika.WithStrictIMEI())
```

## Command codecs

Codecs 12, 13, 14, and 15 carry command responses. Build a `Command`
with a type and a list of `CommandResponse` values.

- A Codec 13 response requires a timestamp.
- A Codec 14 response requires an IMEI.
- A Codec 15 response requires a timestamp and an IMEI.
- A Codec 15 frame may carry a raw response type byte. Store the byte
  as its decimal string so the frame round-trips.

This example builds a Codec 15 packet:

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
				Type: "11",
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

The command type `"11"` is the raw response type byte `0x0B` that some
firmware emits. The encoder writes `0x0B` back to the wire.

## Legacy API

The `pkg` package keeps the old API. It delegates to the public
package.

```go
login := pkg.LoginDecoder(rawLogin)
frame := pkg.TramDecoder(rawFrame)
```

Use the public package in new code. The `pkg` package accepts an
over-declared UDP length for backward compatibility.
