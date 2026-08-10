# Getting started

This section shows you how to add the library to a project and how to
decode and encode a frame.

## Requirements

- Go 1.22.5 or newer
- A project that uses Go modules

## Install the library

Run this command in your project:

```
go get github.com/danieljvsa/teltonika-go
```

## Decode a frame

`Decode` reads a complete frame and returns a packet. The frame can be
a login frame or a data frame. The function detects the frame type
automatically.

This example decodes a Codec 8 TCP frame:

```go
package main

import (
	"encoding/hex"
	"fmt"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func main() {
	frame, err := hex.DecodeString("000000000000003608010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF")
	if err != nil {
		panic(err)
	}

	packet, err := teltonika.Decode(frame)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}

	for _, rec := range packet.Records {
		fmt.Printf("lat=%.6f lon=%.6f speed=%d\n", rec.GPS.Latitude, rec.GPS.Longitude, rec.GPS.Speed)
	}
}
```

The function returns an error if the frame is corrupt. The decoder
checks the declared length, the record count, and the CRC.

## Encode a frame

`Encode` writes a packet back to a complete frame. A TCP frame gets a
transport header and a CRC. The returned frame is ready to send over
the socket.

This example builds a Codec 8 data packet and writes it to a frame:

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
				Timestamp: time.Unix(1701000000, 0).UTC(),
				Priority:  1,
				EventIO:   5,
				GPS: teltonika.GPSData{
					Latitude:   40.4168,
					Longitude:  -3.7038,
					Altitude:   667,
					Angle:      180,
					Satellites: 10,
					Speed:      90,
				},
				IOElements: []teltonika.IOElement{{ID: 1, Value: "01"}},
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

`Encode` and `Decode` are inverse operations. `Decode(Encode(p))`
equals `p` for every packet `p` that the library accepts.

## Next steps

Read datamodel.md to learn the data types.
