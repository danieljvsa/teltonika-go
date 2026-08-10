// Command encode demonstrates encoding Teltonika frames with the public
// package: a Codec 08 data packet and a login packet.
package main

import (
	"encoding/hex"
	"fmt"
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func main() {
	dataPacket := &teltonika.Packet{
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

	frame, err := teltonika.Encode(dataPacket)
	if err != nil {
		panic(err)
	}
	fmt.Printf("data frame (%d bytes): %s\n", len(frame), hex.EncodeToString(frame))

	loginPacket := &teltonika.Packet{Kind: teltonika.KindLogin, IMEI: "356307042441013"}
	login, err := teltonika.Encode(loginPacket)
	if err != nil {
		panic(err)
	}
	fmt.Printf("login frame (%d bytes): %s\n", len(login), hex.EncodeToString(login))
}
