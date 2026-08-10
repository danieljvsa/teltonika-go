package conformance_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// validCodec8Packet returns a valid Codec 8 TCP packet used as shared test
// input across the conformance suite.
func validCodec8Packet() *teltonika.Packet {
	return &teltonika.Packet{
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
}

// udpCodec8Packet wraps validCodec8Packet in a UDP transport with the minimal
// required UDP header.
func udpCodec8Packet() *teltonika.Packet {
	packet := validCodec8Packet()
	packet.Protocol = teltonika.ProtocolUDP
	packet.Header.UDP = &teltonika.HeaderUDP{
		PacketID:    1,
		AVLPacketID: 1,
		IMEI:        "356307042441013",
	}
	return packet
}

// readFixture loads a known-good wire frame from the testdata directory.
// The files contain lowercase hexadecimal text.
func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	raw, err := hex.DecodeString(string(data))
	if err != nil {
		t.Fatalf("decode fixture %s: %v", name, err)
	}
	return raw
}
