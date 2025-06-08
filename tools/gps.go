package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"

	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
)

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
