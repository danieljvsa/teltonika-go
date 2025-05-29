package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Record struct {
	Timestamp   time.Time // Change to actual timestamp type
	Priority    int64
	GPSData     GPSData // Change to actual GPSData type
	EventIO     int64
	NumberOfIOs int64
	IOs         []IOData
}

type CodecData struct {
	NumberOfRecords int64
	Records         []Record
}

type Codec12Data struct{}
type Codec13Data struct{}
type Codec14Data struct{}
type Codec15Data struct{}
type Codec16Data struct{}

func decodeCodec8(data []byte, protocol string) (*CodecData, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	fmt.Println("Records:", numberOfRecords)
	var records []Record
	for range int(numberOfRecords) {
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
		ioData, err := decodeIos(data[read:], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing IO data: %w", err)
		}
		read += int(ioData.LastByte)
		record := Record{
			Timestamp:   *timestamp,
			Priority:    priority,
			GPSData:     *gpsData,
			EventIO:     eventIO,
			NumberOfIOs: ioData.NumberOfIOs,
			IOs:         ioData.IOs,
		}
		records = append(records, record)
	}
	fmt.Println("Protocol:", protocol)
	if protocol == "TCP" {
		tram := append([]byte{0x08}, data...)
		checkCRC := isValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}

	decodedData := &CodecData{
		NumberOfRecords: numberOfRecords,
		Records:         records,
	}
	fmt.Println("Decoded:", decodedData)
	return decodedData, nil
}

func decodeCodec8Ext(data []byte, protocol string) (*CodecData, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}
	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []Record
	for range int(numberOfRecords) {
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
		ioData, err := decodeIos8Extended(data[read:], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing IO data: %w", err)
		}
		read += int(ioData.LastByte)
		record := Record{
			Timestamp:   *timestamp,
			Priority:    priority,
			GPSData:     *gpsData,
			EventIO:     eventIO,
			NumberOfIOs: ioData.NumberOfIOs,
			IOs:         ioData.IOs,
		}
		records = append(records, record)
	}
	if protocol == "TCP" {
		tram := append([]byte{0x8E}, data...)
		checkCRC := isValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &CodecData{
		NumberOfRecords: numberOfRecords,
		Records:         records,
	}
	fmt.Println("Decoded:", decodedData)
	return decodedData, nil
}

func decodeCodec12(data []byte, dataLength int64) (*Codec12Data, error) { return nil, nil }
func decodeCodec13(data []byte, dataLength int64) (*Codec13Data, error) { return nil, nil }
func decodeCodec14(data []byte, dataLength int64) (*Codec14Data, error) { return nil, nil }
func decodeCodec15(data []byte, dataLength int64) (*Codec15Data, error) { return nil, nil }
func decodeCodec16(data []byte, dataLength int64) (*Codec16Data, error) { return nil, nil }
