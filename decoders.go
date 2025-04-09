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

	timestamp := CalcTimestamp(data[2:9])
	if timestamp == nil {
		return nil, fmt.Errorf("error parsing timestamp")
	}

	fmt.Println("Timestamp:", timestamp)

	priority, err := strconv.ParseInt(hex.EncodeToString(data[9:10]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing priority: %w", err)
	}

	gpsData := DecodeGPSData(data[10:25])
	if gpsData == nil {
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

	fmt.Println("Number Of Records:", numberOfRecords)
	fmt.Println("Timestamp:", *timestamp)
	fmt.Println("Priority:", priority)
	fmt.Println("GPSData:", *gpsData)
	fmt.Println("EventIO:", eventIO)
	fmt.Println("IOs:", ioData)
	fmt.Println("CRC16:", data[dataLength:])
	crc := CRC16IBM(data[dataLength:])
	tram := append(data[dataLength:], byte(crc>>8), byte(crc&0xFF))
	checkCRC := VerifyTramCRC(tram)
	if !checkCRC {
		fmt.Println("CRC is not valid")
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
