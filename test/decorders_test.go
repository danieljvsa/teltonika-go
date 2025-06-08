package teltonika_go_test

import (
	"encoding/hex"
	"testing"

	pkg "github.com/danieljvsa/teltonika-go/pkg"
)

func TestDecodeCodec8(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		protocol string
		wantErr  bool
	}{
		{
			name: "Valid Codec8 data (1 byte, 2 bytes, 4 bytes, 8 bytes)",
			data: func() []byte {
				hexStr := "010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "TCP",
			wantErr:  false,
		},
		{
			name: "Valid Codec8 data (1 byte, 2 bytes)",
			data: func() []byte {
				hexStr := "010000016B40D9AD80010000000000000000000000000000000103021503010101425E100000010000F22A" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "TCP",
			wantErr:  false,
		},
		{
			name: "Valid Codec8 data (2 records)",
			data: func() []byte {
				hexStr := "020000016B40D57B480100000000000000000000000000000001010101000000000000016B40D5C198010000000000000000000000000000000101010101000000020000252C" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "TCP",
			wantErr:  false,
		},
		{
			name: "Valid Codec8 data (UDP protocol)",
			data: func() []byte {
				hexStr := "010000016B4F815B30010000000000000000000000000000000103021503010101425DBC000001" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "UDP",
			wantErr:  false,
		},
		{
			name: "Short Data Error",
			data: func() []byte {
				hexStr := "01000000" // Too short to be valid
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "TCP",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pkg.DecodeCodec8(tt.data, tt.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCodec8() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCodec8Ext(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		protocol string
		wantErr  bool
	}{
		{
			name: "Valid Codec8 Extended data",
			data: func() []byte {
				hexStr := "010000016B412CEE000100000000000000000000000000000000010005000100010100010011001D00010010015E2C880002000B000000003544C87A000E000000001DD7E06A00000100002994" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "TCP",
			wantErr:  false,
		},
		{
			name: "Valid Codec8 Extended data UDP",
			data: func() []byte {
				hexStr := "010000016B4F831C680100000000000000000000000000000000010005000100010100010011009D00010010015E2C880002000B000000003544C87A000E000000001DD7E06A000001" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "UDP",
			wantErr:  false,
		},
		{
			name: "Short Data Error",
			data: func() []byte {
				hexStr := "01000000" // Too short to be valid
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			protocol: "TCP",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pkg.DecodeCodec8Ext(tt.data, tt.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCodec8Ext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
