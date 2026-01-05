package tools

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// CalcTimestamp converts an 8-byte big-endian byte slice representing
// a Unix timestamp in milliseconds to a UTC time.Time pointer.
// It uses hex string encoding to parse the bytes as a 64-bit integer.
//
// Parameters:
//   - data: 8-byte slice containing millisecond timestamp
//
// Returns:
//   - *time.Time: pointer to the decoded UTC time
//   - error: parsing error if the byte slice cannot be converted to integer
//
// Example:
//
//	data := []byte{0x00, 0x00, 0x01, 0x6A, 0x57, 0x1F, 0x00, 0x00}
//	timestamp, err := CalcTimestamp(data)
func CalcTimestamp(data []byte) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		return nil, err
	}
	date := time.UnixMilli(timestamp).UTC()
	return &date, nil
}

// CalcTimestampSeconds converts a byte slice (hex-encoded) representing
// a Unix timestamp in seconds to a UTC time.Time pointer.
// Uses hex string encoding to parse bytes as a 64-bit integer.
//
// Parameters:
//   - data: byte slice containing second timestamp
//
// Returns:
//   - *time.Time: pointer to the decoded UTC time
//   - error: parsing error if conversion fails
//
// Example:
//
//	data := []byte{0x56, 0xD8, 0x26, 0xA0}
//	timestamp, err := CalcTimestampSeconds(data)
func CalcTimestampSeconds(data []byte) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		return nil, err
	}
	date := time.Unix(timestamp, 0).UTC()
	return &date, nil
}

// CalcTimestampSecondsLittleEndian converts a little-endian byte slice
// representing a Unix timestamp in seconds to a UTC time.Time pointer.
// Uses encoding/binary package for proper little-endian byte order handling.
//
// Parameters:
//   - data: 4-byte little-endian slice containing second timestamp
//
// Returns:
//   - *time.Time: pointer to the decoded UTC time
//   - error: if data length is insufficient
//
// Example:
//
//	data := []byte{0xA0, 0x26, 0xD8, 0x56}
//	timestamp, err := CalcTimestampSecondsLittleEndian(data)
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
// representing a Unix timestamp in seconds to a UTC time.Time pointer.
// Uses encoding/binary package for proper big-endian byte order handling.
//
// Parameters:
//   - data: 4-byte big-endian slice containing second timestamp
//
// Returns:
//   - *time.Time: pointer to the decoded UTC time
//   - error: if data length is insufficient
//
// Example:
//
//	data := []byte{0x56, 0xD8, 0x26, 0xA0}
//	timestamp, err := CalcTimestampSecondsBigEndian(data)
func CalcTimestampSecondsBigEndian(data []byte) (*time.Time, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for timestamp")
	}
	// Read 4 bytes as big-endian uint32
	timestamp := int64(binary.BigEndian.Uint32(data[:4]))
	date := time.Unix(timestamp, 0).UTC()
	return &date, nil
}
