package conformance_test

import (
	"strings"
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// TestEncodeRejectsInvalidUDPHeaderFields verifies that UDP header fields are
// validated before they are narrowed to their wire widths; out-of-range or
// inconsistent values must not silently wrap.
func TestEncodeRejectsInvalidUDPHeaderFields(t *testing.T) {
	longIMEI := strings.Repeat("a", 65536)

	tests := []struct {
		name   string
		header teltonika.HeaderUDP
	}{
		{name: "negative packet id", header: teltonika.HeaderUDP{PacketID: -1}},
		{name: "packet id overflow", header: teltonika.HeaderUDP{PacketID: 65536}},
		{name: "negative AVL packet id", header: teltonika.HeaderUDP{AVLPacketID: -1}},
		{name: "AVL packet id overflow", header: teltonika.HeaderUDP{AVLPacketID: 256}},
		{name: "negative IMEI length", header: teltonika.HeaderUDP{IMEILength: -1}},
		{name: "IMEI length overflow", header: teltonika.HeaderUDP{IMEILength: 65536}},
		{
			name: "IMEI too long",
			header: teltonika.HeaderUDP{
				IMEI: longIMEI,
			},
		},
		{
			name: "declared IMEI length smaller than actual",
			header: teltonika.HeaderUDP{
				IMEILength: 3,
				IMEI:       "356307042441013",
			},
		},
		{
			name: "declared IMEI length larger than actual",
			header: teltonika.HeaderUDP{
				IMEILength: 99,
				IMEI:       "356307042441013",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := udpCodec8Packet()
			packet.Header.UDP = &tt.header

			if _, err := teltonika.Encode(packet); err == nil {
				t.Fatal("expected invalid UDP header error")
			}
		})
	}
}

// TestEncodeUDPHeaderZeroIMEILengthDerived confirms that a zero IMEILength is
// derived from the actual IMEI bytes on encode, producing a complete frame.
func TestEncodeUDPHeaderZeroIMEILengthDerived(t *testing.T) {
	frame, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(frame) < 12 {
		t.Fatalf("frame too short: % X", frame)
	}

	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.Header.UDP == nil || decoded.Header.UDP.IMEILength != 15 {
		t.Fatalf("expected derived IMEI length 15, got %v", decoded.Header.UDP)
	}
}