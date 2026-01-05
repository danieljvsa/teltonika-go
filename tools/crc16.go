package tools

import (
	"encoding/binary"
)

// Crc16IBM computes the CRC-16 IBM (also known as CRC-16 ModBus) checksum
// for the given byte slice. This is commonly used in Teltonika protocol for
// frame validation.
//
// The CRC-16 IBM algorithm uses polynomial 0xA001 and processes each bit
// of every byte in the input.
//
// Parameters:
//   - data: byte slice to compute checksum for
//
// Returns:
//   - uint16: 16-bit CRC value
//
// Example:
//
//	data := []byte{0x01, 0x02, 0x03}
//	crc := Crc16IBM(data)
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

// IsValidTram validates a Teltonika protocol frame (tram) by checking
// its CRC-16 checksum. The tram is expected to have the CRC value
// as the last 4 bytes in big-endian format.
//
// Frame structure: [data...][CRC(4 bytes big-endian)]
//
// Parameters:
//   - tram: complete frame including CRC at the end
//
// Returns:
//   - bool: true if CRC is valid, false otherwise
//   - Does not return error; returns false for invalid frames
//
// Example:
//
//	isValid := IsValidTram(frameData)
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
