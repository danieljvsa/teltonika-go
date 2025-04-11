package main

import (
	"encoding/hex"
	"strconv"
	"testing"
	"time"
)

func TestDecodeGPSData(t *testing.T) {
	// Example valid data (hex representation of values)
	data, _ := hex.DecodeString("F0A9AFC0209CCA800123456789ab")

	gps := DecodeGPSData(data)
	if gps == nil {
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

func TestCalcTimestamp(t *testing.T) {
	// Example timestamp: 1711765200000 (Unix timestamp in milliseconds)
	// This corresponds to 2024-04-01 12:00:00 UTC.
	// Hex representation: 00018F0E3F380000
	data, _ := hex.DecodeString("0000016B40D8EA30")

	// Expected time
	expectedTime := time.UnixMilli(1560161086000).UTC()

	// Call function
	result := CalcTimestamp(data)
	if result == nil {
		t.Fatal("Expected valid time, got nil")
	}

	// Compare result
	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}
