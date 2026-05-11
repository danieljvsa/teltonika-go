package tools

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
)

// EncodeTimestampMillis converts a time to an 8-byte big-endian millisecond timestamp.
func EncodeTimestampMillis(timestamp *time.Time) ([]byte, error) {
	if timestamp == nil {
		return nil, fmt.Errorf("timestamp is nil")
	}
	if timestamp.Before(time.Unix(0, 0)) {
		return nil, fmt.Errorf("timestamp must be >= Unix epoch")
	}
	value := timestamp.UTC().UnixMilli()
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(value))
	return data, nil
}

// EncodeTimestampSeconds converts a time to a 4-byte big-endian second timestamp.
func EncodeTimestampSeconds(timestamp *time.Time) ([]byte, error) {
	if timestamp == nil {
		return nil, fmt.Errorf("timestamp is nil")
	}
	if timestamp.Before(time.Unix(0, 0)) {
		return nil, fmt.Errorf("timestamp must be >= Unix epoch")
	}
	value := timestamp.UTC().Unix()
	if value > int64(^uint32(0)) {
		return nil, fmt.Errorf("timestamp exceeds uint32 range")
	}
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(value))
	return data, nil
}

// EncodeGPSData converts GPS data into the 14-byte binary Teltonika representation.
func EncodeGPSData(gpsData *tool_domain.GPSData) ([]byte, error) {
	if gpsData == nil {
		return nil, fmt.Errorf("gps data is nil")
	}
	latitudeScaled := math.Round(gpsData.Latitude * 10000000.0)
	longitudeScaled := math.Round(gpsData.Longitude * 10000000.0)
	if latitudeScaled < -2147483648 || latitudeScaled > 2147483647 {
		return nil, fmt.Errorf("latitude out of range")
	}
	if longitudeScaled < -2147483648 || longitudeScaled > 2147483647 {
		return nil, fmt.Errorf("longitude out of range")
	}
	if gpsData.Altitude < -32768 || gpsData.Altitude > 32767 {
		return nil, fmt.Errorf("altitude out of int16 range")
	}
	if gpsData.Angle < 0 || gpsData.Angle > 65535 {
		return nil, fmt.Errorf("angle out of uint16 range")
	}
	if gpsData.Satelites < 0 || gpsData.Satelites > 255 {
		return nil, fmt.Errorf("satellites out of uint8 range")
	}
	if gpsData.Speed < 0 || gpsData.Speed > 255 {
		return nil, fmt.Errorf("speed out of uint8 range")
	}

	data := make([]byte, 14)
	binary.BigEndian.PutUint32(data[0:4], uint32(int32(longitudeScaled)))
	binary.BigEndian.PutUint32(data[4:8], uint32(int32(latitudeScaled)))
	binary.BigEndian.PutUint16(data[8:10], uint16(gpsData.Altitude))
	binary.BigEndian.PutUint16(data[10:12], uint16(gpsData.Angle))
	data[12] = byte(gpsData.Satelites)
	data[13] = byte(gpsData.Speed)
	return data, nil
}

// EncodeIMEI converts a hex-string IMEI into its byte representation.
func EncodeIMEI(imei string) ([]byte, error) {
	if imei == "" {
		return nil, fmt.Errorf("IMEI is empty")
	}
	if len(imei)%2 != 0 {
		return nil, fmt.Errorf("IMEI hex string must have even length")
	}
	data, err := hex.DecodeString(imei)
	if err != nil {
		return nil, fmt.Errorf("invalid IMEI hex: %w", err)
	}
	if len(data) < 8 {
		return nil, fmt.Errorf("IMEI data too short: %d bytes", len(data))
	}
	return data, nil
}

// EncodeHexMessage returns the byte representation of a command response.
// If hexMessage is provided it is decoded, otherwise response is encoded as ASCII.
func EncodeHexMessage(response string, hexMessage string) ([]byte, error) {
	if hexMessage != "" {
		if len(hexMessage)%2 != 0 {
			return nil, fmt.Errorf("hex message must have even length")
		}
		data, err := hex.DecodeString(hexMessage)
		if err != nil {
			return nil, fmt.Errorf("invalid hex message: %w", err)
		}
		return data, nil
	}
	if response == "" {
		return nil, fmt.Errorf("response is empty")
	}
	return []byte(response), nil
}

// AppendCRC16IBM appends CRC-16 IBM checksum (4-byte big endian) to data.
func AppendCRC16IBM(data []byte) []byte {
	crc := Crc16IBM(data)
	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, uint32(crc))
	return append(data, crcBytes...)
}
