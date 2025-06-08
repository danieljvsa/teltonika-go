package tools

import (
	"encoding/hex"
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
