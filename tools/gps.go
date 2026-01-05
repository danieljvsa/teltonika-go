package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"

	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
)

// DecodeGPSData parses a 14-byte binary GPS data block from Teltonika protocol
// and extracts latitude, longitude, altitude, angle, satellite count, and speed.
//
// Data format (14 bytes):
//   - Bytes 0-3: Longitude (int32, degrees * 10^7)
//   - Bytes 4-7: Latitude (int32, degrees * 10^7)
//   - Bytes 8-9: Altitude (int16, meters)
//   - Bytes 10-11: Angle (int16, degrees 0-359)
//   - Byte 12: Number of satellites
//   - Byte 13: Speed (1 byte)
//
// Parameters:
//   - data: exactly 14 bytes of GPS data
//
// Returns:
//   - *tool_domain.GPSData: pointer to decoded GPS information
//   - error: if data length is invalid or parsing fails
//
// Example:
//
//	gpsData, err := DecodeGPSData(rawGPSBytes)
//	if err == nil {
//		fmt.Printf("Location: %f, %f\n", gpsData.Latitude, gpsData.Longitude)
//	}
func DecodeGPSData(data []byte) (*tool_domain.GPSData, error) {
	if len(data) < 14 {
		return nil, fmt.Errorf("invalid data length %d", len(data))
	}

	lat, err := strconv.ParseUint(hex.EncodeToString(data[4:8]), 16, 32)
	if err != nil {
		return nil, err
	}

	long, err := strconv.ParseUint(hex.EncodeToString(data[:4]), 16, 32)
	if err != nil {
		return nil, err
	}
	altitude, err := strconv.ParseInt(hex.EncodeToString(data[8:10]), 16, 64)
	if err != nil {
		return nil, err
	}
	angle, err := strconv.ParseInt(hex.EncodeToString(data[10:12]), 16, 64)
	if err != nil {
		return nil, err
	}
	satelites, err := strconv.ParseInt(hex.EncodeToString(data[12:13]), 16, 64)
	if err != nil {
		return nil, err
	}

	speed, err := strconv.ParseInt(hex.EncodeToString(data[13:14]), 16, 64)
	if err != nil {
		return nil, err
	}

	longitude := float64(int32(long)) / 10000000.0
	latitude := float64(int32(lat)) / 10000000.0

	return &tool_domain.GPSData{
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
		Angle:     angle,
		Satelites: satelites,
		Speed:     speed,
	}, nil
}
