package tools

import (
	"encoding/hex"
	"fmt"
)

// DecodeIMEI extracts and converts IMEI (International Mobile Equipment Identity)
// data from a byte slice to its hexadecimal string representation.
//
// IMEI is a unique identifier for mobile devices, typically 15 digits represented
// as 8 bytes in Teltonika protocol (after decoding the BCD or hex format).
//
// Parameters:
//   - data: byte slice containing IMEI data (minimum 8 bytes)
//
// Returns:
//   - string: hexadecimal string representation of IMEI bytes
//   - error: if data is empty or too short (less than 8 bytes)
//
// Example:
//
//	imeiBytes := []byte{0x35, 0x64, 0x23, 0x40, 0x12, 0x34, 0x56, 0x78}
//	imei, err := DecodeIMEI(imeiBytes)
//	if err == nil {
//		fmt.Println(imei) // Output: 3564234012345678
//	}
func DecodeIMEI(data []byte) (string, error) {
	// Decode hex-encoded bytes to IMEI string
	// IMEI is typically 15 bytes (30 hex characters)
	if len(data) == 0 {
		return "", fmt.Errorf("empty IMEI data")
	}
	if len(data) < 8 {
		return "", fmt.Errorf("IMEI data too short: %d bytes", len(data))
	}

	hexString := hex.EncodeToString(data)
	return hexString, nil
}
