package conformance_test

import (
	"encoding/binary"
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func TestEncodeRejectsUnknownProtocol(t *testing.T) {
	packet := validCodec8Packet()
	packet.Protocol = teltonika.Protocol("HTTP")

	if _, err := teltonika.Encode(packet); err == nil {
		t.Fatal("expected unsupported protocol error")
	}
}

func TestEncodeDefaultsEmptyProtocolToTCP(t *testing.T) {
	packet := validCodec8Packet()
	packet.Protocol = ""

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(frame) < 8 || frame[0] != 0 || frame[1] != 0 || frame[2] != 0 || frame[3] != 0 {
		t.Fatalf("expected TCP preamble in frame: % X", frame)
	}
	declared := binary.BigEndian.Uint32(frame[4:8])
	if int(declared) != len(frame)-12 {
		t.Fatalf("DataLength expected to exclude CRC: declared=%d", declared)
	}
}

func TestEncodeUDPRequiresUDPHeader(t *testing.T) {
	packet := validCodec8Packet()
	packet.Protocol = teltonika.ProtocolUDP

	if _, err := teltonika.Encode(packet); err == nil {
		t.Fatal("expected UDP header required error")
	}
}

func TestDecodeRejectsUDPLengthMismatch(t *testing.T) {
	valid, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	tests := []struct {
		name   string
		change int32
	}{
		{name: "declared length too small", change: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := append([]byte(nil), valid...)
			length := binary.BigEndian.Uint16(frame[0:2])
			binary.BigEndian.PutUint16(frame[0:2], uint16(int32(length)+tt.change))

			if _, err := teltonika.Decode(frame); err == nil {
				t.Fatal("expected UDP length mismatch error")
			}
		})
	}
}

// TestDecodeUDPAcceptsLargerDeclaredLength locks in the lenient behavior for
// declared lengths that exceed the delivered frame. Real devices and truncated
// captures (such as the official wiki codec 16 example) declare more bytes
// than are present; only declared lengths too small for the data are invalid.
func TestDecodeUDPAcceptsLargerDeclaredLength(t *testing.T) {
	valid, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	frame := append([]byte(nil), valid...)
	length := binary.BigEndian.Uint16(frame[0:2])
	binary.BigEndian.PutUint16(frame[0:2], length+1)

	if _, err := teltonika.Decode(frame); err != nil {
		t.Fatalf("expected decode with larger declared length, got error: %v", err)
	}
}

func TestDecodeUDPAcceptsValidLength(t *testing.T) {
	valid, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if _, err := teltonika.Decode(valid); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
}