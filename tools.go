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

func CalcTimestamp(data []byte) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return nil, err
	}
	date := time.UnixMilli(timestamp).UTC()
	return &date, nil
}

func DecodeGPSData(data []byte) (*GPSData, error) {
	if len(data) < 14 {
		fmt.Println("Invalid data length", len(data))
		return nil, fmt.Errorf("Invalid data length %d", len(data))
	}

	lat, err := strconv.ParseUint(hex.EncodeToString(data[4:8]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
		return nil, err
	}

	long, err := strconv.ParseUint(hex.EncodeToString(data[:4]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
		return nil, err
	}
	altitude, err := strconv.ParseInt(hex.EncodeToString(data[8:10]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing altitude:", err)
		return nil, err
	}
	angle, err := strconv.ParseInt(hex.EncodeToString(data[10:12]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing angle:", err)
		return nil, err
	}
	satelites, err := strconv.ParseInt(hex.EncodeToString(data[12:13]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing satelites:", err)
		return nil, err
	}

	speed, err := strconv.ParseInt(hex.EncodeToString(data[13:14]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing speed:", err)
		return nil, err
	}

	longitude := float64(int32(long)) / 10000000.0
	latitude := float64(int32(lat)) / 10000000.0

	return &GPSData{
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
		Angle:     angle,
		Satelites: satelites,
		Speed:     speed,
	}, nil
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
	return receivedCRC == calculatedCRC
}
