package conformance_test

import (
	"testing"
	"time"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// TestCodec12EncodesResponseCount verifies that the count fields written at
// the beginning and end of a Codec 12 payload equal the number of encoded
// responses, not the number of command groups.
func TestCodec12EncodesResponseCount(t *testing.T) {
	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolTCP,
		Codec:    teltonika.Codec12,
		Commands: []teltonika.Command{
			{
				Type: "Response",
				Responses: []teltonika.CommandResponse{
					{Response: "FIRST"},
					{Response: "SECOND"},
				},
			},
		},
	}

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if frame[8] != byte(teltonika.Codec12) {
		t.Fatalf("unexpected codec: 0x%02X", frame[8])
	}
	if frame[9] != 2 {
		t.Fatalf("encoded count=%d, expected=2", frame[9])
	}

	// Four final bytes are CRC; the preceding byte is the trailing count.
	if trailingCount := frame[len(frame)-5]; trailingCount != 2 {
		t.Fatalf("trailing count=%d, expected=2", trailingCount)
	}

	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded.Commands) != 1 {
		t.Fatalf("expected one command group, got %d", len(decoded.Commands))
	}
	if got := len(decoded.Commands[0].Responses); got != 2 {
		t.Fatalf("expected two responses, got %d", got)
	}
}

func TestEncodeRejectsMultipleCommandGroups(t *testing.T) {
	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolTCP,
		Codec:    teltonika.Codec12,
		Commands: []teltonika.Command{
			{Type: "Response", Responses: []teltonika.CommandResponse{{Response: "ONE"}}},
			{Type: "Response", Responses: []teltonika.CommandResponse{{Response: "TWO"}}},
		},
	}

	if _, err := teltonika.Encode(packet); err == nil {
		t.Fatal("expected multiple command groups to be rejected")
	}
}

// TestDecodeRejectsCommandBytesAfterTrailingCount verifies that a command
// payload is fully consumed: only the trailing count (plus TCP CRC) may follow
// the responses. UDP has no CRC guard, so stray trailing bytes must be caught
// by the decoder itself.
func TestDecodeRejectsCommandBytesAfterTrailingCount(t *testing.T) {
	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolUDP,
		Codec:    teltonika.Codec12,
		Header: teltonika.Header{
			UDP: &teltonika.HeaderUDP{},
		},
		Commands: []teltonika.Command{
			{Type: "Response", Responses: []teltonika.CommandResponse{{Response: "12"}}},
		},
	}

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// UDP frame: length(2) packetId(2) version(1) avlPacketId(1) imeiLen(2)
	// imei(0) codec(1) codec-data. The empty IMEI keeps the codec at offset 8.
	if len(frame) < 10 || frame[8] != byte(teltonika.Codec12) {
		t.Fatalf("unexpected UDP frame layout: % X", frame)
	}
	payload := append([]byte(nil), frame[9:]...)

	if _, err := teltonika.DecodeCodecData(payload, teltonika.Codec12, teltonika.ProtocolUDP); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}

	extra := append(payload, 0x00)
	if _, err := teltonika.DecodeCodecData(extra, teltonika.Codec12, teltonika.ProtocolUDP); err == nil {
		t.Fatal("expected bytes after trailing count to be rejected")
	}
}

func TestCodec13ExactRoundTrip(t *testing.T) {
	ts := time.Unix(1701000000, 123000000).UTC()
	packet := commandPacket(teltonika.Codec13, nil, &ts, nil, "getinfo")

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	resp := decoded.Commands[0].Responses[0]
	if !resp.Timestamp.Equal(ts) {
		t.Fatalf("timestamp mismatch: %v != %v", resp.Timestamp, ts)
	}
	if resp.Response != "getinfo" {
		t.Fatalf("response mismatch: %q", resp.Response)
	}
}

func TestCodec14ExactRoundTrip(t *testing.T) {
	imei := "0123456789abcdef" // lowercase: hex.EncodeToString output
	const message = "getinfo"

	packet := commandPacket(teltonika.Codec14, &imei, nil, nil, message)
	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded.Commands) != 1 || len(decoded.Commands[0].Responses) != 1 {
		t.Fatal("unexpected decoded command structure")
	}

	response := decoded.Commands[0].Responses[0]
	if response.IMEI != imei {
		t.Errorf("IMEI=%q, expected=%q", response.IMEI, imei)
	}
	if response.Response != message {
		t.Errorf("response=%q, expected=%q", response.Response, message)
	}
}

func TestCodec15ExactRoundTrip(t *testing.T) {
	imei := "0123456789abcdef"
	const message = "PING"
	ts := time.Unix(1701000000, 0).UTC() // codec 15 carries second precision

	packet := commandPacket(teltonika.Codec15, &imei, &ts, nil, message)

	frame, err := teltonika.Encode(packet)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := teltonika.Decode(frame)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	resp := decoded.Commands[0].Responses[0]
	if resp.IMEI != imei {
		t.Errorf("IMEI=%q, expected=%q", resp.IMEI, imei)
	}
	if resp.Response != message {
		t.Errorf("response=%q, expected=%q", resp.Response, message)
	}
	if !resp.Timestamp.Equal(ts) {
		t.Errorf("timestamp=%v, expected=%v", resp.Timestamp, ts)
	}
}

func commandPacket(codec teltonika.CodecID, imei *string, ts *time.Time, commandType *string, response string) *teltonika.Packet {
	typ := "Response"
	if commandType != nil {
		typ = *commandType
	}
	resp := teltonika.CommandResponse{Response: response, Timestamp: ts}
	if imei != nil {
		resp.IMEI = *imei
	}
	return &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolTCP,
		Codec:    codec,
		Commands: []teltonika.Command{
			{Type: typ, Responses: []teltonika.CommandResponse{resp}},
		},
	}
}
