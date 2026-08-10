package conformance_test

import (
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func TestDecodeCodecDataRejectsUnknownProtocol(t *testing.T) {
	_, err := teltonika.DecodeCodecData(
		[]byte{0x01},
		teltonika.Codec8,
		teltonika.Protocol("HTTP"),
	)
	if err == nil {
		t.Fatal("expected unsupported protocol error")
	}
}

// TestDecodeCodecDataRejectsEmptyProtocol ensures an empty protocol is not
// silently treated as TCP or UDP: the presence of a CRC depends on the
// transport, so the caller must state it explicitly.
func TestDecodeCodecDataRejectsEmptyProtocol(t *testing.T) {
	_, err := teltonika.DecodeCodecData(
		[]byte{0x01},
		teltonika.Codec8,
		teltonika.Protocol(""),
	)
	if err == nil {
		t.Fatal("expected unsupported protocol error for empty protocol")
	}
}
