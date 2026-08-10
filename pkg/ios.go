package teltonika_go

import (
	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// DecodeIos8 decodes a Codec 8 I/O block and returns the legacy internal
// response model.
func DecodeIos8(data []byte, startByte int64) (*io_domain.ResponseDecode, error) {
	return decodeIos(data, startByte, 1, false, false)
}

// DecodeIos8Extended decodes a Codec 8E I/O block and returns the legacy
// internal response model.
func DecodeIos8Extended(data []byte, startByte int64) (*io_domain.ResponseDecode, error) {
	return decodeIos(data, startByte, 2, false, true)
}

// DecodeIos16 decodes a Codec 16 I/O block and returns the legacy internal
// response model, including the generation type.
func DecodeIos16(data []byte, startByte int64) (*io_domain.ResponseDecode, error) {
	return decodeIos(data, startByte, 2, true, false)
}

func decodeIos(data []byte, startByte int64, idSize int, withGen bool, hasXGroup bool) (*io_domain.ResponseDecode, error) {
	// Preserve legacy behavior: too-short input decodes to an empty result.
	if int64(len(data))-startByte < 4 {
		return &io_domain.ResponseDecode{IOs: []io_domain.IOData{}, NumberOfIOs: 0, LastByte: 0}, nil
	}

	decoded, err := teltonika.DecodeCodecIO(data, int(startByte), idSize, withGen, hasXGroup)
	if err != nil {
		return &io_domain.ResponseDecode{IOs: []io_domain.IOData{}, NumberOfIOs: 0, LastByte: int64(decoded.NextOffset), GenerationType: decoded.GenerationType}, err
	}

	legacy := make([]io_domain.IOData, 0, len(decoded.Elements))
	for _, io := range decoded.Elements {
		legacy = append(legacy, io_domain.IOData{IO: io.ID, Value: io.Value})
	}

	return &io_domain.ResponseDecode{
		IOs:            legacy,
		NumberOfIOs:    int64(len(decoded.Elements)),
		LastByte:       int64(decoded.NextOffset),
		GenerationType: decoded.GenerationType,
	}, nil
}
