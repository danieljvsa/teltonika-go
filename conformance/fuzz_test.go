package conformance_test

import (
	"encoding/hex"
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

func FuzzDecodeNeverPanics(f *testing.F) {
	seed, _ := hex.DecodeString("000000000000003608010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF")
	f.Add(seed)
	f.Add([]byte("000F333536333037303432343431303133"))
	f.Add([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = teltonika.Decode(data)
	})
}

func FuzzDecodeLoginNeverPanics(f *testing.F) {
	f.Add([]byte("000F333536333037303432343431303133"))
	f.Add([]byte{0x00, 0x01, 0x31})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = teltonika.DecodeLogin(data)
	})
}

func FuzzDecodeCodecDataNeverPanics(f *testing.F) {
	// Codec 14 over UDP exercises the previously panic-prone IMEI slice path.
	f.Add([]byte{0x01, 0x06, 0x00, 0x00, 0x00, 0x48, 0x01})
	f.Add([]byte{0x01, 0x06, 0x00, 0x00, 0x00, 0x08})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = teltonika.DecodeCodecData(data, teltonika.Codec14, teltonika.ProtocolUDP)
		_, _ = teltonika.DecodeCodecData(data, teltonika.Codec8, teltonika.ProtocolUDP)
	})
}