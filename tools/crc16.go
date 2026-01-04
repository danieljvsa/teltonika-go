package tools

import (
	"encoding/binary"
)

func Crc16IBM(data []byte) uint16 {
	var crc uint16 = 0x0000
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

func IsValidTram(tram []byte) bool {
	if len(tram) < 2 {
		return false
	}
	length := len(tram)
	data := tram[:len(tram)-4]
	receivedCRC := uint16(binary.BigEndian.Uint32(tram[length-4:]))
	calculatedCRC := Crc16IBM(data)
	return receivedCRC == calculatedCRC
}
