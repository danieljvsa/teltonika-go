package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type GPSData struct {
	Latitude  float64
	Longitude float64
	Altitude  int64
	Angle     int64
	Satelites int64
	Speed     int64
}

func CalcTimestamp(data []byte) *time.Time {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return nil
	}
	date := time.UnixMilli(timestamp).UTC()
	return &date
}

func DecodeGPSData(data []byte) *GPSData {
	if len(data) < 14 {
		fmt.Println("Invalid data length", len(data))
		return nil
	}

	lat, err := strconv.ParseUint(hex.EncodeToString(data[4:8]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
		return nil
	}

	long, err := strconv.ParseUint(hex.EncodeToString(data[:4]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
		return nil
	}
	altitude, err := strconv.ParseInt(hex.EncodeToString(data[8:10]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing altitude:", err)
		return nil
	}
	angle, err := strconv.ParseInt(hex.EncodeToString(data[10:12]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing angle:", err)
		return nil
	}
	satelites, err := strconv.ParseInt(hex.EncodeToString(data[12:13]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing satelites:", err)
		return nil
	}

	speed, err := strconv.ParseInt(hex.EncodeToString(data[13:14]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing speed:", err)
		return nil
	}

	longitude := float64(int32(long)) / 10000000.0
	latitude := float64(int32(lat)) / 10000000.0

	fmt.Println("Latitude:", latitude)
	fmt.Println("Longitude:", longitude)
	fmt.Println("Altitude:", altitude)
	fmt.Println("Angle:", angle)
	fmt.Println("Satelites:", satelites)
	fmt.Println("Speed:", speed)

	return &GPSData{
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
		Angle:     angle,
		Satelites: satelites,
		Speed:     speed,
	}
}

// CRC16 IBM/ARC implementation (poly = 0x8005, reflected = true, init = 0x0000)
func crc16IBM(data []byte) uint16 {
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

// Validate tram with last 2 bytes as CRC (little-endian)
func isValidTram(tram []byte) bool {
	if len(tram) < 2 {
		return false
	}
	length := len(tram)
	data := tram[:len(tram)-4]
	receivedCRC := uint16(binary.BigEndian.Uint32(tram[length-4:]))
	calculatedCRC := crc16IBM(data)
	fmt.Println(receivedCRC, calculatedCRC)
	fmt.Printf("Calculated CRC-16: 0x%04X (little-endian: [%02X %02X])\n", calculatedCRC, byte(calculatedCRC&0xFF), byte(calculatedCRC>>8))
	fmt.Printf("Received CRC-16: 0x%04X (little-endian: [%02X %02X])\n", receivedCRC, byte(receivedCRC&0xFF), byte(receivedCRC>>8))
	return receivedCRC == calculatedCRC
}
