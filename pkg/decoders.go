package teltonika_go

import (
	"encoding/hex"
	"fmt"
	"strconv"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func DecodeCodec8(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []decoder_domain.Record
	for range int(numberOfRecords) {
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp")
		}
		read += 8
		priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing priority: %w", err)
		}
		read += 1
		gpsData, err := tools.DecodeGPSData(data[read : read+15])
		if err != nil {
			return nil, fmt.Errorf("error parsing GPS data")
		}
		read += 15
		eventIO, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing event IO: %w", err)
		}
		read += 1
		ioData, err := DecodeIos8(data[read:], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing IO data: %w", err)
		}
		read += int(ioData.LastByte)
		record := &decoder_domain.Record{
			Timestamp:   *timestamp,
			Priority:    priority,
			GPSData:     *gpsData,
			EventIO:     eventIO,
			NumberOfIOs: ioData.NumberOfIOs,
			IOs:         ioData.IOs,
		}
		records = append(records, *record)
	}

	if protocol == "TCP" {
		tram := append([]byte{0x08}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}

	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfRecords,
		Records:         records,
	}
	return decodedData, nil
}

func DecodeCodec8Ext(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}
	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []decoder_domain.Record
	for range int(numberOfRecords) {
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp")
		}
		read += 8
		priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing priority: %w", err)
		}
		read += 1
		gpsData, err := tools.DecodeGPSData(data[read : read+15])
		if err != nil {
			return nil, fmt.Errorf("error parsing GPS data")
		}
		read += 15
		eventIO, err := strconv.ParseInt(hex.EncodeToString(data[read:read+2]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing event IO: %w", err)
		}
		read += 2
		ioData, err := DecodeIos8Extended(data[read:], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing IO data: %w", err)
		}
		read += int(ioData.LastByte)
		record := &decoder_domain.Record{
			Timestamp:   *timestamp,
			Priority:    priority,
			GPSData:     *gpsData,
			EventIO:     eventIO,
			NumberOfIOs: ioData.NumberOfIOs,
			IOs:         ioData.IOs,
		}

		records = append(records, *record)
	}
	if protocol == "TCP" {
		tram := append([]byte{0x8E}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfRecords,
		Records:         records,
	}
	return decodedData, nil
}

func DecodeCodec16(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if len(data) < 29 {
		return nil, fmt.Errorf("data length too short")
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []decoder_domain.Record
	for range int(numberOfRecords) {
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp")
		}
		read += 8
		priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing priority: %w", err)
		}
		read += 1
		gpsData, err := tools.DecodeGPSData(data[read : read+15])
		if err != nil {
			return nil, fmt.Errorf("error parsing GPS data")
		}
		read += 15
		eventIO, err := strconv.ParseInt(hex.EncodeToString(data[read:read+2]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing event IO: %w", err)
		}
		read += 2
		ioData, err := DecodeIos16(data[read:], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing IO data: %w", err)
		}
		read += int(ioData.LastByte)
		record := &decoder_domain.Record{
			Timestamp:   *timestamp,
			Priority:    priority,
			GPSData:     *gpsData,
			EventIO:     eventIO,
			NumberOfIOs: ioData.NumberOfIOs,
			IOs:         ioData.IOs,
		}
		records = append(records, *record)
	}

	if protocol == "TCP" {
		tram := append([]byte{0x10}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}

	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfRecords,
		Records:         records,
	}
	return decodedData, nil
}

func DecodeCodec12(data []byte, protocol string) (*decoder_domain.CodecData, error) { return nil, nil }
func DecodeCodec14(data []byte, protocol string) (*decoder_domain.CodecData, error) { return nil, nil }
