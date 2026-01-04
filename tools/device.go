package tools

import (
	"encoding/hex"
	"fmt"
)

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
