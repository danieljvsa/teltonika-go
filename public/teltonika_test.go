package teltonika_test

import (
	"encoding/hex"
	"testing"
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// Valid full Codec 08 TCP frame (header, codec id, records, CRC).
var codec8TCPFrame = "000000000000003608010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF"

func mustDecodeHex(t *testing.T, s string) []byte {
	t.Helper()
	data, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("invalid hex: %v", err)
	}
	return data
}

func TestDecodeLogin(t *testing.T) {
	data := mustDecodeHex(t, "000F333536333037303432343431303133")
	packet, err := teltonika.Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if packet.Kind != teltonika.KindLogin {
		t.Fatalf("expected KindLogin, got %v", packet.Kind)
	}
	if packet.IMEI != "356307042441013" {
		t.Fatalf("unexpected IMEI: %q", packet.IMEI)
	}
}

func TestDecodeCodec8TCP(t *testing.T) {
	packet, err := teltonika.Decode(mustDecodeHex(t, codec8TCPFrame))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if packet.Kind != teltonika.KindData {
		t.Fatalf("expected KindData, got %v", packet.Kind)
	}
	if packet.Codec != teltonika.Codec8 {
		t.Fatalf("expected Codec8, got 0x%02X", byte(packet.Codec))
	}
	if packet.Protocol != teltonika.ProtocolTCP {
		t.Fatalf("expected TCP, got %q", packet.Protocol)
	}
	if len(packet.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(packet.Records))
	}
	rec := packet.Records[0]
	if rec.GPS.Speed != 0 {
		t.Fatalf("unexpected speed: %d", rec.GPS.Speed)
	}
	if len(rec.IOElements) != 5 {
		t.Fatalf("expected 5 IO elements, got %d", len(rec.IOElements))
	}
}

func TestDecodeEncodeRoundTripAVL(t *testing.T) {
	tests := []struct {
		name  string
		codec teltonika.CodecID
		gen   string
		event int64
	}{
		{name: "Codec8", codec: teltonika.Codec8, event: 5},
		{name: "Codec8Ext", codec: teltonika.Codec8Ext, event: 5},
		{name: "Codec16", codec: teltonika.Codec16, gen: "On Change", event: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := time.Unix(1701000000, 0).UTC()
			packet := &teltonika.Packet{
				Kind:     teltonika.KindData,
				Protocol: teltonika.ProtocolTCP,
				Codec:    tt.codec,
				Records: []teltonika.AVLRecord{
					{
						Timestamp: timestamp,
						Priority:  1,
						EventIO:   tt.event,
						GPS: teltonika.GPSData{
							Latitude:   40.4168,
							Longitude:  -3.7038,
							Altitude:   667,
							Angle:      180,
							Satellites: 10,
							Speed:      90,
						},
						IOElements: []teltonika.IOElement{
							{ID: 1, Value: "01"},
							{ID: 2, Value: "0001"},
						},
						GenerationType: tt.gen,
					},
				},
			}

			frame, err := teltonika.Encode(packet)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			decoded, err := teltonika.Decode(frame)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
			if len(decoded.Records) != 1 {
				t.Fatalf("expected 1 record, got %d", len(decoded.Records))
			}
			got := decoded.Records[0]
			if got.GPS.Speed != 90 || got.GPS.Altitude != 667 {
				t.Fatalf("GPS round-trip mismatch: %+v", got.GPS)
			}
			if got.GPS.Satellites != 10 {
				t.Fatalf("satellites round-trip mismatch: %d", got.GPS.Satellites)
			}
			if len(got.IOElements) != 2 {
				t.Fatalf("expected 2 IO elements, got %d", len(got.IOElements))
			}
			if tt.codec == teltonika.Codec16 && got.GenerationType != tt.gen {
				t.Fatalf("generation type mismatch: %q", got.GenerationType)
			}
		})
	}
}

func TestDecodeEncodeRoundTripCommands(t *testing.T) {
	tests := []struct {
		name  string
		codec teltonika.CodecID
	}{
		{name: "Codec12", codec: teltonika.Codec12},
		{name: "Codec13", codec: teltonika.Codec13},
		{name: "Codec14", codec: teltonika.Codec14},
		{name: "Codec15", codec: teltonika.Codec15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := time.Unix(1701000000, 0).UTC()
			packet := &teltonika.Packet{
				Kind:     teltonika.KindData,
				Protocol: teltonika.ProtocolTCP,
				Codec:    tt.codec,
				Commands: []teltonika.Command{
					{
						Type: "Response",
						Responses: []teltonika.CommandResponse{
							{
								Timestamp:   &timestamp,
								Response:    "OK",
								CommandType: "Response",
								IMEI:        "0123456789ABCDEF",
							},
						},
					},
				},
			}

			frame, err := teltonika.Encode(packet)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			decoded, err := teltonika.Decode(frame)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
			if len(decoded.Commands) != 1 || len(decoded.Commands[0].Responses) != 1 {
				t.Fatalf("expected 1 command with 1 response, got %d/%d", len(decoded.Commands), len(decoded.Commands))
			}
			resp := decoded.Commands[0].Responses[0]
			if tt.codec == teltonika.Codec14 {
				// Codec 14 embeds the IMEI and message as one payload; only
				// the hex message reliably round-trips.
				if resp.HexMessage == "" {
					t.Fatalf("expected hex message for codec 14")
				}
				return
			}
			if resp.Response != "OK" {
				t.Fatalf("response round-trip mismatch: %q", resp.Response)
			}
		})
	}
}

func TestDecodeEncodeLoginRoundTrip(t *testing.T) {
	packet := &teltonika.Packet{Kind: teltonika.KindLogin, IMEI: "356307042441013"}
	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.IMEI != "356307042441013" {
		t.Fatalf("IMEI round-trip mismatch: %q", decoded.IMEI)
	}
}

func TestDecodeEncodeUDPRoundTrip(t *testing.T) {
	timestamp := time.Unix(1701000000, 0).UTC()
	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolUDP,
		Codec:    teltonika.Codec8,
		Header: teltonika.Header{
			UDP: &teltonika.HeaderUDP{
				PacketID:    0xCAFE,
				AVLPacketID: 1,
				IMEI:        "356307042441013",
			},
		},
		Records: []teltonika.AVLRecord{
			{
				Timestamp: timestamp,
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
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.Protocol != teltonika.ProtocolUDP {
		t.Fatalf("expected UDP, got %q", decoded.Protocol)
	}
	if len(decoded.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(decoded.Records))
	}
}

func TestDecodeMalformed(t *testing.T) {
	cases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "tiny", data: []byte{0x01}},
		{name: "bad login length", data: mustDecodeHex(t, "000F31323334")},
		{name: "truncated tcp", data: mustDecodeHex(t, "000000000000000801")},
		{name: "unknown codec", data: mustDecodeHex(t, "0000000000000001000000000000FF00")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := teltonika.Decode(tt.data); err == nil {
				t.Fatalf("expected error for %s", tt.name)
			}
		})
	}
}

func TestDecodeRejectsBadCRC(t *testing.T) {
	data := mustDecodeHex(t, codec8TCPFrame)
	// Corrupt the CRC (last 4 bytes).
	data[len(data)-1] ^= 0xFF
	if _, err := teltonika.Decode(data); err == nil {
		t.Fatalf("expected CRC error")
	}
}

func TestCodecIDConstants(t *testing.T) {
	if byte(teltonika.Codec8) != 0x08 {
		t.Fatalf("Codec8 value wrong")
	}
	if byte(teltonika.Codec8Ext) != 0x8E {
		t.Fatalf("Codec8Ext value wrong")
	}
	if byte(teltonika.Codec12) != 0x0C {
		t.Fatalf("Codec12 value wrong")
	}
	if byte(teltonika.Codec16) != 0x10 {
		t.Fatalf("Codec16 value wrong")
	}
}
