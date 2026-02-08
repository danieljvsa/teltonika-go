package teltonika_go_test

import (
	"encoding/hex"
	"testing"

	pkg "github.com/danieljvsa/teltonika-go/pkg"
)

func TestDecodeCodec8RejectsTruncatedRecord(t *testing.T) {
	data, _ := hex.DecodeString("01")
	_, err := pkg.DecodeCodec8(data, "TCP")
	if err == nil {
		t.Fatalf("expected error for truncated codec 8 data")
	}
}

func TestDecodeCodec12RejectsShortPayload(t *testing.T) {
	data, _ := hex.DecodeString("010600000001")
	_, err := pkg.DecodeCodec12(data, "TCP")
	if err == nil {
		t.Fatalf("expected error for truncated codec 12 data")
	}
}

func TestDecodeCodec15RejectsShortPayload(t *testing.T) {
	data, _ := hex.DecodeString("010600000004FFFFFFFF")
	_, err := pkg.DecodeCodec15(data, "TCP")
	if err == nil {
		t.Fatalf("expected error for truncated codec 15 data")
	}
}
