package teltonika

import (
	"encoding/binary"
	"fmt"

	tools "github.com/danieljvsa/teltonika-go/tools"
)

// crcValid reports whether the CRC appended to a TCP payload is valid.
// The payload must include the trailing 4-byte CRC, preceded by the data
// over which the CRC is computed (codec id and payload, excluding the header).
func crcValid(codec byte, payload []byte) bool {
	if len(payload) < 4 {
		return false
	}
	tram := make([]byte, 0, len(payload)+1)
	tram = append(tram, codec)
	tram = append(tram, payload...)
	return tools.IsValidTram(tram)
}

// appendCRC appends the 4-byte CRC-16 IBM checksum of the given data.
func appendCRC(data []byte) []byte {
	return tools.AppendCRC16IBM(data)
}

// decodeGPSData parses the fixed 15-byte GPS block of an AVL record.
func decodeGPSData(data []byte) (GPSData, error) {
	if len(data) < 15 {
		return GPSData{}, fmt.Errorf("invalid GPS data length %d", len(data))
	}
	longitude := int32(binary.BigEndian.Uint32(data[0:4]))
	latitude := int32(binary.BigEndian.Uint32(data[4:8]))
	return GPSData{
		Latitude:   float64(latitude) / 10000000.0,
		Longitude:  float64(longitude) / 10000000.0,
		Altitude:   int64(int16(binary.BigEndian.Uint16(data[8:10]))),
		Angle:      int64(int16(binary.BigEndian.Uint16(data[10:12]))),
		Satellites: int64(data[12]),
		Speed:      int64(binary.BigEndian.Uint16(data[13:15])),
	}, nil
}

// encodeGPSData encodes GPSData into the fixed 15-byte Teltonika GPS block.
func encodeGPSData(gps GPSData) ([]byte, error) {
	if gps.Altitude < -32768 || gps.Altitude > 32767 {
		return nil, fmt.Errorf("altitude out of int16 range")
	}
	if gps.Angle < 0 || gps.Angle > 65535 {
		return nil, fmt.Errorf("angle out of uint16 range")
	}
	if gps.Satellites < 0 || gps.Satellites > 255 {
		return nil, fmt.Errorf("satellites out of uint8 range")
	}
	if gps.Speed < 0 || gps.Speed > 65535 {
		return nil, fmt.Errorf("speed out of uint16 range")
	}

	latitude := int32(gps.Latitude * 10000000.0)
	longitude := int32(gps.Longitude * 10000000.0)

	data := make([]byte, 15)
	binary.BigEndian.PutUint32(data[0:4], uint32(longitude))
	binary.BigEndian.PutUint32(data[4:8], uint32(latitude))
	binary.BigEndian.PutUint16(data[8:10], uint16(gps.Altitude))
	binary.BigEndian.PutUint16(data[10:12], uint16(gps.Angle))
	data[12] = byte(gps.Satellites)
	binary.BigEndian.PutUint16(data[13:15], uint16(gps.Speed))
	return data, nil
}
