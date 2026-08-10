// Command decode demonstrates decoding a Teltonika frame with the public
// package: a full Codec 08 TCP data frame and a login frame.
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

	fmt.Printf("protocol=%s codec=0x%02X records=%d\n", packet.Protocol, byte(packet.Codec), len(packet.Records))
	for _, rec := range packet.Records {
		fmt.Printf("  lat=%.6f lon=%.6f speed=%d io=%d\n", rec.GPS.Latitude, rec.GPS.Longitude, rec.GPS.Speed, len(rec.IOElements))
	}
}
