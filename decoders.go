package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Codec8Data struct {
	NumberOfRecords int64
	Timestamp       time.Time // Change to actual timestamp type
	Priority        int64
	GPSData         GPSData // Change to actual GPSData type
	EventIO         int64
	NumberOfIOs     int64
	IOs             []IOData
}

type Codec8EData struct {
	NumberOfRecords int64
	Timestamp       time.Time // Change to actual timestamp type
	Priority        int64
	GPSData         GPSData // Change to actual GPSData type
	EventIO         int64
	NumberOfIOs     int64
	IOs             []IOData
}

func decodeCodec8(data []byte, dataLength int64) (*Codec8Data, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	timestamp, err := CalcTimestamp(data[read : read+8])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp")
	}
	read += 8
	priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing priority: %w", err)
	}
	read += 1
	gpsData, err := DecodeGPSData(data[read : read+15])
	if err != nil {
		return nil, fmt.Errorf("error parsing GPS data")
	}
	read += 15
	eventIO, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing event IO: %w", err)
	}
	read += 1
	ioData, err := decodeIos(data[read:dataLength], int64(len(data[read:dataLength])), 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing IO data: %w", err)
	}

	tram := append([]byte{0x08}, data...)
	checkCRC := isValidTram(tram)
	if !checkCRC {
		return nil, fmt.Errorf("CRC is not valid")
	}

	return &Codec8Data{
		NumberOfRecords: numberOfRecords,
		Timestamp:       *timestamp,
		Priority:        priority,
		GPSData:         *gpsData,
		EventIO:         eventIO,
		NumberOfIOs:     ioData.NumberOfIOs,
		IOs:             ioData.IOs,
	}, nil
}

func decodeCodec8Ext(data []byte, dataLength int64) (*Codec8EData, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}
	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	timestamp, err := CalcTimestamp(data[read : read+8])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp")
	}
	read += 8
	priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing priority: %w", err)
	}
	read += 1
	gpsData, err := DecodeGPSData(data[read : read+15])
	if err != nil {
		return nil, fmt.Errorf("error parsing GPS data")
	}
	read += 15
	eventIO, err := strconv.ParseInt(hex.EncodeToString(data[read:read+2]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing event IO: %w", err)
	}
	read += 2
	ioData, err := decodeIos8Extended(data[read:dataLength], int64(len(data[read:dataLength])), 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing IO data: %w", err)
	}
	tram := append([]byte{0x8E}, data...)
	checkCRC := isValidTram(tram)
	if !checkCRC {
		return nil, fmt.Errorf("CRC is not valid")
	}

	return &Codec8EData{
		NumberOfRecords: numberOfRecords,
		Timestamp:       *timestamp,
		Priority:        priority,
		GPSData:         *gpsData,
		EventIO:         eventIO,
		NumberOfIOs:     ioData.NumberOfIOs,
		IOs:             ioData.IOs,
	}, nil
}

func decodeCodec12(data []byte, dataLength int64) {}
func decodeCodec13(data []byte, dataLength int64) {}
func decodeCodec14(data []byte, dataLength int64) {}
func decodeCodec15(data []byte, dataLength int64) {}
func decodeCodec16(data []byte, dataLength int64) {}
