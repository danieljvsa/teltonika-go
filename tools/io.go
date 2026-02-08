package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

var generationTypeEncoding = map[string]byte{
	"On Exit":     0,
	"On Entrance": 1,
	"On Both":     2,
	"Reserved":    3,
	"Hysteresis":  4,
	"On Change":   5,
	"Eventual":    6,
	"Periodical":  7,
}

func GetGenerationType(data []byte, startByte int64, length int64) (string, error) {
	byte := startByte
	generation_type_translation := "Unknown generation type"

	if len(data) < int(length) || int(length) < 1 {
		return generation_type_translation, fmt.Errorf("data or length is too small")
	}

	generation_type, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+length]), 16, 64)
	if err != nil {
		return generation_type_translation, err
	}

	switch generation_type {
	case 0:
		generation_type_translation = "On Exit"
		return generation_type_translation, nil
	case 1:
		generation_type_translation = "On Entrance"
		return generation_type_translation, nil
	case 2:
		generation_type_translation = "On Both"
		return generation_type_translation, nil
	case 3:
		generation_type_translation = "Reserved"
		return generation_type_translation, nil
	case 4:
		generation_type_translation = "Hysteresis"
		return generation_type_translation, nil
	case 5:
		generation_type_translation = "On Change"
		return generation_type_translation, nil
	case 6:
		generation_type_translation = "Eventual"
		return generation_type_translation, nil
	case 7:
		generation_type_translation = "Periodical"
		return generation_type_translation, nil
	default:
		return generation_type_translation, nil
	}
}

// EncodeGenerationType converts a generation type name to its byte value.
//
// Valid values mirror Teltonika generation types used by Codec 16:
//   - "On Exit"
//   - "On Entrance"
//   - "On Both"
//   - "Reserved"
//   - "Hysteresis"
//   - "On Change"
//   - "Eventual"
//   - "Periodical"
func EncodeGenerationType(generationType string) (byte, error) {
	value, ok := generationTypeEncoding[generationType]
	if !ok {
		return 0, fmt.Errorf("unknown generation type: %s", generationType)
	}
	return value, nil
}
