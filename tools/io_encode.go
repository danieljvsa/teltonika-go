package tools

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"

	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
)

func encodeIOValue(value string) ([]byte, error) {
	if value == "" {
		return nil, fmt.Errorf("IO value is empty")
	}
	if len(value)%2 != 0 {
		return nil, fmt.Errorf("IO value must have even-length hex")
	}
	data, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid IO value hex: %w", err)
	}
	return data, nil
}

// EncodeIOData8 encodes IO elements for Codec 8.
func EncodeIOData8(ios []io_domain.IOData) ([]byte, int64, error) {
	var oneByte, twoByte, fourByte, eightByte []io_domain.IOData

	for _, io := range ios {
		if io.IO < 0 || io.IO > 255 {
			return nil, 0, fmt.Errorf("IO ID out of byte range: %d", io.IO)
		}
		valueBytes, err := encodeIOValue(io.Value)
		if err != nil {
			return nil, 0, err
		}
		switch len(valueBytes) {
		case 1:
			oneByte = append(oneByte, io)
		case 2:
			twoByte = append(twoByte, io)
		case 4:
			fourByte = append(fourByte, io)
		case 8:
			eightByte = append(eightByte, io)
		default:
			return nil, 0, fmt.Errorf("invalid IO value length: %d bytes", len(valueBytes))
		}
	}

	if len(ios) > 255 {
		return nil, 0, fmt.Errorf("IO count exceeds uint8 range")
	}

	buffer := &bytes.Buffer{}
	buffer.WriteByte(byte(len(ios)))

	writeGroup := func(group []io_domain.IOData) error {
		buffer.WriteByte(byte(len(group)))
		sort.Slice(group, func(i, j int) bool {
			return group[i].IO < group[j].IO
		})
		for _, io := range group {
			buffer.WriteByte(byte(io.IO))
			valueBytes, err := encodeIOValue(io.Value)
			if err != nil {
				return err
			}
			buffer.Write(valueBytes)
		}
		return nil
	}

	if err := writeGroup(oneByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(twoByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(fourByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(eightByte); err != nil {
		return nil, 0, err
	}

	return buffer.Bytes(), int64(len(ios)), nil
}

// EncodeIOData8Extended encodes IO elements for Codec 8E (extended IO format).
func EncodeIOData8Extended(ios []io_domain.IOData) ([]byte, int64, error) {
	var oneByte, twoByte, fourByte, eightByte, xByte []io_domain.IOData

	for _, io := range ios {
		if io.IO < 0 || io.IO > 65535 {
			return nil, 0, fmt.Errorf("IO ID out of uint16 range: %d", io.IO)
		}
		valueBytes, err := encodeIOValue(io.Value)
		if err != nil {
			return nil, 0, err
		}
		switch len(valueBytes) {
		case 1:
			oneByte = append(oneByte, io)
		case 2:
			twoByte = append(twoByte, io)
		case 4:
			fourByte = append(fourByte, io)
		case 8:
			eightByte = append(eightByte, io)
		default:
			xByte = append(xByte, io)
		}
	}

	if len(ios) > 65535 {
		return nil, 0, fmt.Errorf("IO count exceeds uint16 range")
	}

	buffer := &bytes.Buffer{}
	countBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(countBytes, uint16(len(ios)))
	buffer.Write(countBytes)

	writeGroup := func(group []io_domain.IOData) error {
		groupCount := make([]byte, 2)
		binary.BigEndian.PutUint16(groupCount, uint16(len(group)))
		buffer.Write(groupCount)
		sort.Slice(group, func(i, j int) bool {
			return group[i].IO < group[j].IO
		})
		for _, io := range group {
			idBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(idBytes, uint16(io.IO))
			buffer.Write(idBytes)
			valueBytes, err := encodeIOValue(io.Value)
			if err != nil {
				return err
			}
			buffer.Write(valueBytes)
		}
		return nil
	}

	if err := writeGroup(oneByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(twoByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(fourByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(eightByte); err != nil {
		return nil, 0, err
	}

	xCount := make([]byte, 2)
	binary.BigEndian.PutUint16(xCount, uint16(len(xByte)))
	buffer.Write(xCount)
	sort.Slice(xByte, func(i, j int) bool {
		return xByte[i].IO < xByte[j].IO
	})
	for _, io := range xByte {
		valueBytes, err := encodeIOValue(io.Value)
		if err != nil {
			return nil, 0, err
		}
		idBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(idBytes, uint16(io.IO))
		buffer.Write(idBytes)
		if len(valueBytes) > 65535 {
			return nil, 0, fmt.Errorf("IO value too large: %d bytes", len(valueBytes))
		}
		lengthBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(lengthBytes, uint16(len(valueBytes)))
		buffer.Write(lengthBytes)
		buffer.Write(valueBytes)
	}

	return buffer.Bytes(), int64(len(ios)), nil
}

// EncodeIOData16 encodes IO elements for Codec 16, including generation type.
func EncodeIOData16(ios []io_domain.IOData, generationType string) ([]byte, int64, error) {
	var oneByte, twoByte, fourByte, eightByte []io_domain.IOData

	genValue, err := EncodeGenerationType(generationType)
	if err != nil {
		return nil, 0, err
	}

	for _, io := range ios {
		if io.IO < 0 || io.IO > 65535 {
			return nil, 0, fmt.Errorf("IO ID out of uint16 range: %d", io.IO)
		}
		valueBytes, err := encodeIOValue(io.Value)
		if err != nil {
			return nil, 0, err
		}
		switch len(valueBytes) {
		case 1:
			oneByte = append(oneByte, io)
		case 2:
			twoByte = append(twoByte, io)
		case 4:
			fourByte = append(fourByte, io)
		case 8:
			eightByte = append(eightByte, io)
		default:
			return nil, 0, fmt.Errorf("invalid IO value length: %d bytes", len(valueBytes))
		}
	}

	if len(ios) > 255 {
		return nil, 0, fmt.Errorf("IO count exceeds uint8 range")
	}

	buffer := &bytes.Buffer{}
	buffer.WriteByte(genValue)
	buffer.WriteByte(byte(len(ios)))

	writeGroup := func(group []io_domain.IOData) error {
		buffer.WriteByte(byte(len(group)))
		sort.Slice(group, func(i, j int) bool {
			return group[i].IO < group[j].IO
		})
		for _, io := range group {
			idBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(idBytes, uint16(io.IO))
			buffer.Write(idBytes)
			valueBytes, err := encodeIOValue(io.Value)
			if err != nil {
				return err
			}
			buffer.Write(valueBytes)
		}
		return nil
	}

	if err := writeGroup(oneByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(twoByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(fourByte); err != nil {
		return nil, 0, err
	}
	if err := writeGroup(eightByte); err != nil {
		return nil, 0, err
	}

	return buffer.Bytes(), int64(len(ios)), nil
}
