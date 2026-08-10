package conformance_test

import (
	"bytes"
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func TestDecodeOfficialCodec8Fixture(t *testing.T) {
	packet, err := teltonika.Decode(readFixture(t, "codec8-tcp.hex"))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if packet.Codec != teltonika.Codec8 || len(packet.Records) != 1 {
		t.Fatalf("unexpected packet: codec=%02X records=%d", byte(packet.Codec), len(packet.Records))
	}
}

// TestRoundTripOfficialFixtures decodes each known-good fixture, re-encodes it
// and compares the complete frame byte for byte. The codec 16 UDP fixture is
// covered decode-only: it is the truncated official wiki packet whose declared
// length (347) points beyond the delivered data, so a byte-identical
// re-encode is impossible.
func TestRoundTripOfficialFixtures(t *testing.T) {
	fixtures := []string{
		"codec8-tcp.hex",
		"codec8e-tcp.hex",
		"codec16-tcp.hex",
		"codec12-command.hex",
		"codec13-response.hex",
		"codec14-response.hex",
		"codec15-response.hex",
		"codec8-udp.hex",
		"codec8e-udp.hex",
	}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			fixture := readFixture(t, name)
			packet, err := teltonika.Decode(fixture)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
			encoded, err := teltonika.Encode(packet)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}
			if !bytes.Equal(encoded, fixture) {
				t.Fatalf("encoded packet mismatch\nexpected: %X\nactual:   %X", fixture, encoded)
			}
		})
	}
}

func TestDecodeOfficialCodec15Fixture(t *testing.T) {
	packet, err := teltonika.Decode(readFixture(t, "codec15-response.hex"))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	resp := packet.Commands[0].Responses[0]
	if resp.Response != "Hello!\n" {
		t.Fatalf("unexpected response: %q", resp.Response)
	}
	if resp.IMEI != "0123456789123456" {
		t.Fatalf("unexpected IMEI: %q", resp.IMEI)
	}
}

func TestDecodeOfficialCodec16UDPFixture(t *testing.T) {
	packet, err := teltonika.Decode(readFixture(t, "codec16-udp.hex"), teltonika.WithLenientUDPLength())
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if packet.Codec != teltonika.Codec16 || len(packet.Records) != 1 {
		t.Fatalf("unexpected packet: codec=%02X records=%d", byte(packet.Codec), len(packet.Records))
	}
	if packet.Header.UDP == nil || packet.Header.UDP.Length != 0x015B {
		t.Fatalf("expected declared UDP length 0x015B preserved, got %v", packet.Header.UDP)
	}
}