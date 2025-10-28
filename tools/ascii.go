package tools

import (
	"encoding/hex"
	"fmt"
)

func DecodeToHexThenASCII(data []byte, message []byte, size int) (string, string, error) {
	src := data
	if len(message) > 0 {
		src = message
	}
	if size <= 0 {
		return "", "", fmt.Errorf("invalid size: %d", size)
	}
	if len(src) < size {
		return "", "", fmt.Errorf("source length %d is smaller than requested size %d", len(src), size)
	}

	segment := src[:size]
	hexStr := hex.EncodeToString(segment)

	decodedBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return hexStr, "", fmt.Errorf("hex decode error: %w", err)
	}

	return hexStr, string(decodedBytes), nil
}
