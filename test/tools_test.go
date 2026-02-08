package teltonika_go_test

import (
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func TestDecodeGPSData(t *testing.T) {
	// Example valid data (hex representation of values)
	data, _ := hex.DecodeString("F0A9AFC0209CCA800123456789ab")

	gps, err := tools.DecodeGPSData(data)
	if err != nil {
		t.Fatal("Expected valid GPSData, got nil")
	}

	expectedLong, err := strconv.ParseUint("F0A9AFC0", 16, 32)
	if err != nil {
		t.Fatal("Error parsing longitude:", err)
	}

	expectedLat, err := strconv.ParseUint("209CCA80", 16, 32)
	if err != nil {
		t.Fatal("Error parsing longitude:", err)
	}

	expectedLatitude := float64(int32(expectedLat)) / 10000000.0
	expectedLongitude := float64(int32(expectedLong)) / 10000000.0
	expectedAltitude := int64(0x0123)
	expectedAngle := int64(0x4567)
	expectedSatelites := int64(0x89)
	expectedSpeed := int64(0xab)

	if gps.Latitude != expectedLatitude {
		t.Errorf("Latitude: expected %v, got %v", expectedLatitude, gps.Latitude)
	}
	if gps.Longitude != expectedLongitude {
		t.Errorf("Longitude: expected %v, got %v", expectedLongitude, gps.Longitude)
	}
	if gps.Altitude != expectedAltitude {
		t.Errorf("Altitude: expected %d, got %d", expectedAltitude, gps.Altitude)
	}
	if gps.Angle != expectedAngle {
		t.Errorf("Angle: expected %d, got %d", expectedAngle, gps.Angle)
	}
	if gps.Satelites != expectedSatelites {
		t.Errorf("Satelites: expected %d, got %d", expectedSatelites, gps.Satelites)
	}
	if gps.Speed != expectedSpeed {
		t.Errorf("Speed: expected %d, got %d", expectedSpeed, gps.Speed)
	}
}

func TestEncodeGPSDataRoundTrip(t *testing.T) {
	gps := &tool_domain.GPSData{
		Latitude:  37.7749,
		Longitude: -122.4194,
		Altitude:  15,
		Angle:     90,
		Satelites: 8,
		Speed:     55,
	}

	encoded, err := tools.EncodeGPSData(gps)
	if err != nil {
		t.Fatalf("EncodeGPSData failed: %v", err)
	}

	decoded, err := tools.DecodeGPSData(encoded)
	if err != nil {
		t.Fatalf("DecodeGPSData failed: %v", err)
	}

	if decoded.Speed != gps.Speed {
		t.Errorf("expected speed %d, got %d", gps.Speed, decoded.Speed)
	}
}

func TestCalcTimestamp(t *testing.T) {
	// Example timestamp: 1711765200000 (Unix timestamp in milliseconds)
	// This corresponds to 2024-04-01 12:00:00 UTC.
	// Hex representation: 00018F0E3F380000
	data, _ := hex.DecodeString("0000016B40D8EA30")

	// Expected time
	expectedTime := time.UnixMilli(1560161086000).UTC()

	// Call function
	result, err := tools.CalcTimestamp(data)
	if err != nil {
		t.Fatal("Expected valid time, got nil")
	}

	// Compare result
	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestEncodeTimestampSeconds(t *testing.T) {
	ts := time.Unix(1457006240, 0).UTC()
	encoded, err := tools.EncodeTimestampSeconds(&ts)
	if err != nil {
		t.Fatalf("EncodeTimestampSeconds failed: %v", err)
	}
	decoded, err := tools.CalcTimestampSecondsBigEndian(encoded)
	if err != nil {
		t.Fatalf("CalcTimestampSecondsBigEndian failed: %v", err)
	}
	if !decoded.Equal(ts) {
		t.Errorf("expected %v, got %v", ts, decoded)
	}
}

func TestCalcTimestampSeconds(t *testing.T) {
	// Test with seconds-based timestamp
	// 0x56D826A0 = 1457006240 seconds (2016-03-03 11:57:20 UTC)
	data, _ := hex.DecodeString("56D826A0")
	result, err := tools.CalcTimestampSeconds(data)
	if err != nil {
		t.Fatal("CalcTimestampSeconds failed:", err)
	}

	expectedTime := time.Unix(1457006240, 0).UTC()
	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestCalcTimestampSecondsBigEndian(t *testing.T) {
	tests := []struct {
		name      string
		hexInput  string
		wantTime  time.Time
		wantError bool
	}{
		{
			name:     "Valid big-endian timestamp",
			hexInput: "56D826A0",
			wantTime: time.Unix(1457006240, 0).UTC(),
		},
		{
			name:      "Empty data",
			hexInput:  "",
			wantError: true,
		},
		{
			name:      "Too short data",
			hexInput:  "5656",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tt.hexInput)
			result, err := tools.CalcTimestampSecondsBigEndian(data)

			if (err != nil) != tt.wantError {
				t.Errorf("wantError %v, got %v", tt.wantError, err)
			}
			if !tt.wantError && !result.Equal(tt.wantTime) {
				t.Errorf("Expected %v, got %v", tt.wantTime, result)
			}
		})
	}
}

func TestDecodeIMEI(t *testing.T) {
	tests := []struct {
		name      string
		hexInput  string
		wantIMEI  string
		wantError bool
	}{
		{
			name:     "Valid IMEI",
			hexInput: "333536333037303432343431303133",
			wantIMEI: "333536333037303432343431303133",
		},
		{
			name:      "Empty data",
			hexInput:  "",
			wantError: true,
		},
		{
			name:      "Too short IMEI",
			hexInput:  "3536",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tt.hexInput)
			result, err := tools.DecodeIMEI(data)

			if (err != nil) != tt.wantError {
				t.Errorf("wantError %v, got %v", tt.wantError, err)
			}
			if !tt.wantError && result != tt.wantIMEI {
				t.Errorf("Expected %s, got %s", tt.wantIMEI, result)
			}
		})
	}
}

func TestDecodeToHexThenASCII(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		message   string
		size      int
		wantHex   string
		wantASCII string
		wantError bool
	}{
		{
			name:      "Valid hex decode",
			data:      "48656C6C6F", // "Hello" in hex
			message:   "",
			size:      5,
			wantHex:   "48656c6c6f",
			wantASCII: "Hello",
		},
		{
			name:      "Invalid size",
			data:      "48656C",
			message:   "",
			size:      10,
			wantError: true,
		},
		{
			name:      "Zero size",
			data:      "48656C",
			message:   "",
			size:      0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tt.data)
			msg, _ := hex.DecodeString(tt.message)

			hex, ascii, err := tools.DecodeToHexThenASCII(data, msg, tt.size)

			if (err != nil) != tt.wantError {
				t.Errorf("wantError %v, got %v", tt.wantError, err)
			}
			if !tt.wantError {
				if hex != tt.wantHex {
					t.Errorf("Hex: expected %s, got %s", tt.wantHex, hex)
				}
				if ascii != tt.wantASCII {
					t.Errorf("ASCII: expected %s, got %s", tt.wantASCII, ascii)
				}
			}
		})
	}
}

func TestCrc16IBM(t *testing.T) {
	tests := []struct {
		name     string
		hexInput string
		wantCRC  uint16
	}{
		{
			name:     "Valid tram CRC",
			hexInput: "0d01060000000f0000016c0a81c320676574696e666f01",
			wantCRC:  0x5b66,
		},
		{
			name:     "Empty data",
			hexInput: "",
			wantCRC:  0x0000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tt.hexInput)
			result := tools.Crc16IBM(data)
			if result != tt.wantCRC {
				t.Errorf("Expected %04x, got %04x", tt.wantCRC, result)
			}
		})
	}
}

func TestIsLogin(t *testing.T) {
	tests := []struct {
		name      string
		hexInput  string
		wantLogin bool
		wantError bool
	}{
		{
			name:      "Valid login",
			hexInput:  "000F333536333037303432343431303133",
			wantLogin: true,
		},
		{
			name:      "Empty data",
			hexInput:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tt.hexInput)
			result, err := tools.IsLogin(data)

			if (err != nil) != tt.wantError {
				t.Errorf("wantError %v, got %v", tt.wantError, err)
			}
			if !tt.wantError && result != tt.wantLogin {
				t.Errorf("Expected %v, got %v", tt.wantLogin, result)
			}
		})
	}
}
