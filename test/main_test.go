package teltonika_go_test

import (
	"encoding/hex"
	"testing"

	tg "github.com/danieljvsa/teltonika-go/cmd/teltonika-go"
)

func TestLoginDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string // hex string
		wantErr bool
	}{
		{
			name:  "Login",
			input: "000F333536333037303432343431303133",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.input)
			if err != nil {
				t.Fatalf("invalid test input: %v", err)
			}
			res := tg.LoginDecoder(data)
			if res.Error != nil {
				t.Fatalf("invalid internal test input: %v", res.Error)
			}
		})
	}
}

func TestTramDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string // hex string
		wantErr bool
	}{
		{
			name:  "Valid Codec 08 TCP",
			input: "000000000000003608010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF", // header, length 3+1=4, codec 08
		},
		{
			name:  "Valid Codec 8E TCP",
			input: "000000000000004A8E010000016B412CEE000100000000000000000000000000000000010005000100010100010011001D00010010015E2C880002000B000000003544C87A000E000000001DD7E06A00000100002994",
		},
		{
			name:  "Valid Codec 08 UDP",
			input: "003DCAFE0105000F33353230393330383634303336353508010000016B4F815B30010000000000000000000000000000000103021503010101425DBC000001",
		},
		{
			name:  "Valid Codec 8E UDP",
			input: "005FCAFE0107000F3335323039333038363430333635358E010000016B4F831C680100000000000000000000000000000000010005000100010100010011009D00010010015E2C880002000B000000003544C87A000E000000001DD7E06A000001",
		},
		{
			name:  "Valid Codec 16 TCP",
			input: "000000000000005F10020000016BDBC7833000000000000000000000000000000000000B05040200010000030002000B00270042563A00000000016BDBC7871800000000000000000000000000000000000B05040200010000030002000B00260042563A00000200005FB3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.input)
			if err != nil {
				t.Fatalf("invalid test input: %v", err)
			}
			res := tg.TramDecoder(data)
			if res.Error != nil {
				t.Fatalf("invalid internal test input: %v", res.Error)
			}
		})
	}
}
