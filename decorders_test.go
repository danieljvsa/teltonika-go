package main

import (
	"encoding/hex"
	"testing"
)

func TestDecodeCodec8(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		dataLength int64
		wantErr    bool
	}{
		{
			name: "Valid Codec8 data",
			data: func() []byte {
				hexStr := "010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 53,
			wantErr:    false,
		},
		{
			name: "Short Data Error",
			data: func() []byte {
				hexStr := "01000000" // Too short to be valid
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 4,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeCodec8(tt.data, tt.dataLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCodec8() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCodec8Ext(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		dataLength int64
		wantErr    bool
	}{
		{
			name: "Valid Codec8 Extended data",
			data: func() []byte {
				hexStr := "010000016B412CEE000100000000000000000000000000000000010005000100010100010011001D00010010015E2C880002000B000000003544C87A000E000000001DD7E06A00000100002994" // Mock valid data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 73,
			wantErr:    false,
		},
		{
			name: "Short Data Error",
			data: func() []byte {
				hexStr := "01000000" // Too short to be valid
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 4,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeCodec8Ext(tt.data, tt.dataLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCodec8Ext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
