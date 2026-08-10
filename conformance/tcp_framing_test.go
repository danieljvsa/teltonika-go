package conformance_test

import (
	"encoding/binary"
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// TestEncodeTCPDataLengthExcludesCRC verifies that the TCP Data Field Length
// covers the codec id, records and trailing record count but not the
// eight-byte header or four-byte CRC.
func TestEncodeTCPDataLengthExcludesCRC(t *testing.T) {
	packet := validCodec8Packet()

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(frame) < 12 {
		t.Fatalf("frame too short: %d", len(frame))
	}

	declared := binary.BigEndian.Uint32(frame[4:8])
	expected := uint32(len(frame) - 8 - 4)

	if declared != expected {
		t.Fatalf("incorrect DataLength: declared=%d expected=%d frameSize=%d", declared, expected, len(frame))
	}
}

func TestDecodeRejectsTCPDataLengthMismatch(t *testing.T) {
	valid, err := teltonika.Encode(validCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	tests := []struct {
		name   string
		change int32
	}{
		{name: "declared length too small", change: -1},
		{name: "declared length too large", change: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := append([]byte(nil), valid...)
			length := binary.BigEndian.Uint32(frame[4:8])
			binary.BigEndian.PutUint32(frame[4:8], uint32(int32(length)+tt.change))

			if _, err := teltonika.Decode(frame); err == nil {
				t.Fatal("expected data-length mismatch error")
			}
		})
	}
}

func TestDecodeRejectsTruncatedTCPFrame(t *testing.T) {
	valid, err := teltonika.Encode(validCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	for _, n := range []int{4, 8, 12, len(valid) - 4, len(valid) - 1} {
		if _, err := teltonika.Decode(valid[:n]); err == nil {
			t.Fatalf("expected error for truncated frame of %d bytes", n)
		}
	}
}

func TestDecodeRejectsTCPFrameWithTrailingBytes(t *testing.T) {
	valid, err := teltonika.Encode(validCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	frame := append(append([]byte(nil), valid...), 0x00, 0x00)
	if _, err := teltonika.Decode(frame); err == nil {
		t.Fatal("expected error for frame with trailing bytes")
	}
}
