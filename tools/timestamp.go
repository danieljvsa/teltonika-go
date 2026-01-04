package tools

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func CalcTimestamp(data []byte) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		return nil, err
	}
	date := time.UnixMilli(timestamp).UTC()
	return &date, nil
}

// CalcTimestampSeconds converts a byte slice (hex-encoded) representing
// Unix timestamp in seconds to a UTC time.Time pointer.
func CalcTimestampSeconds(data []byte) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		return nil, err
	}
	date := time.Unix(timestamp, 0).UTC()
	return &date, nil
}

// CalcTimestampSecondsLittleEndian converts a little-endian byte slice
// representing Unix timestamp in seconds to a UTC time.Time pointer.
func CalcTimestampSecondsLittleEndian(data []byte) (*time.Time, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for timestamp")
	}
	// Read 4 bytes as little-endian uint32
	timestamp := int64(binary.LittleEndian.Uint32(data[:4]))
	date := time.Unix(timestamp, 0).UTC()
	return &date, nil
}

// CalcTimestampSecondsBigEndian converts a big-endian byte slice
// representing Unix timestamp in seconds to a UTC time.Time pointer.
func CalcTimestampSecondsBigEndian(data []byte) (*time.Time, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for timestamp")
	}
	// Read 4 bytes as big-endian uint32
	timestamp := int64(binary.BigEndian.Uint32(data[:4]))
	date := time.Unix(timestamp, 0).UTC()
	return &date, nil
}
