package conformance_test

import (
	"testing"
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// TestDecodeRejectsMismatchedAVLRecordCounts verifies that the initial and
// trailing AVL record counts must match. UDP is used so that the trailing
// count can be mutated without recomputing a CRC.
func TestDecodeRejectsMismatchedAVLRecordCounts(t *testing.T) {
	frame, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// The final byte is the trailing record count in this UDP frame.
	frame[len(frame)-1] = 2

	if _, err := teltonika.Decode(frame); err == nil {
		t.Fatal("expected record-count mismatch error")
	}
}

func TestDecodeRejectsMissingTrailingRecordCount(t *testing.T) {
	frame, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	frame = frame[:len(frame)-1]

	if _, err := teltonika.Decode(frame); err == nil {
		t.Fatal("expected missing trailing record count error")
	}
}

func TestDecodeAcceptsMatchingMultipleRecordCounts(t *testing.T) {
	packet := udpCodec8Packet()
	second := packet.Records[0]
	second.Timestamp = time.Unix(1702000000, 0).UTC()
	packet.Records = append(packet.Records, second)

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(decoded.Records))
	}
}

func TestDecodeRejectsBytesAfterTrailingRecordCount(t *testing.T) {
	frame, err := teltonika.Encode(udpCodec8Packet())
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	frame = append(frame, 0x00)

	if _, err := teltonika.Decode(frame); err == nil {
		t.Fatal("expected bytes-after-trailing-count error")
	}
}