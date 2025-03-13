package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func CalcTimestamp(data []byte) *int64 {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return nil
	}

	return &timestamp
}
