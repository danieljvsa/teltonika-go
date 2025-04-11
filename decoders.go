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

func decodeCodec8(data []byte, dataLength int64) (*Codec8Data, error) {
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}

	timestamp, err := CalcTimestamp(data[2:9])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp")
	}

	priority, err := strconv.ParseInt(hex.EncodeToString(data[9:10]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing priority: %w", err)
	}

	gpsData, err := DecodeGPSData(data[10:25])
	if err != nil {
		return nil, fmt.Errorf("error parsing GPS data")
	}

	eventIO, err := strconv.ParseInt(hex.EncodeToString(data[25:26]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing event IO: %w", err)
	}

	ioData, err := decodeIos(data[26:dataLength], int64(len(data[26:dataLength])), 0)
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

func decodeCodec8Ext(data []byte, dataLength int64) {}
func decodeCodec12(data []byte, dataLength int64)   {}
func decodeCodec13(data []byte, dataLength int64)   {}
func decodeCodec14(data []byte, dataLength int64)   {}
func decodeCodec15(data []byte, dataLength int64)   {}
func decodeCodec16(data []byte, dataLength int64)   {}
