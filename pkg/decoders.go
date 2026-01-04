package teltonika_go

import (
	"encoding/hex"
	"fmt"
	"strconv"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	tools_domain "github.com/danieljvsa/teltonika-go/internal/tool"
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
			Timestamp:   timestamp,
			Priority:    &priority,
			GPSData:     gpsData,
			EventIO:     &eventIO,
			NumberOfIOs: &ioData.NumberOfIOs,
			IOs:         &ioData.IOs,
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
			Timestamp:   timestamp,
			Priority:    &priority,
			GPSData:     gpsData,
			EventIO:     &eventIO,
			NumberOfIOs: &ioData.NumberOfIOs,
			IOs:         &ioData.IOs,
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
			Timestamp:   timestamp,
			Priority:    &priority,
			GPSData:     gpsData,
			EventIO:     &eventIO,
			NumberOfIOs: &ioData.NumberOfIOs,
			IOs:         &ioData.IOs,
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

func DecodeCodec12(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if len(data) < 12 {
		return nil, fmt.Errorf("data length too short")
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	fmt.Printf("Number of Commands: %d\n", numberOfCommands)
	read += 1
	responseTypeNumber, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing response type: %w", err)
	}
	var responseType string
	switch responseTypeNumber {
	case 5:
		responseType = "Command"
	case 6:
		responseType = "Response"
	default:
		return nil, fmt.Errorf("unknown response type: %d", responseTypeNumber)
	}
	fmt.Printf("Response Type: %d (%s)\n", responseTypeNumber, responseType)
	read += 1
	var command_responses []tools_domain.CommandResponse
	for range int(numberOfCommands) {
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		fmt.Printf("Response Size: %d | %s\n", responseSize, hex.EncodeToString(data[read:read+4]))
		read += 4
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(responseSize)], []byte{}, int(responseSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		fmt.Printf("Decoded Command Response: %s\n", message)
		fmt.Printf("Hex Command Response: %s\n", hexString)
		read += int(responseSize)
		commandResponse := &tools_domain.CommandResponse{
			Response:   message,
			HexMessage: hexString,
		}
		command_responses = append(command_responses, *commandResponse)
	}
	responseType2, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing size of message: %w", err)
	}
	read += 1
	if numberOfCommands != responseType2 {
		return nil, fmt.Errorf("response type mismatch")
	}
	if protocol == "TCP" {
		tram := append([]byte{0x0C}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfCommands,
		Records: []decoder_domain.Record{
			{
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}

// Note: Codec13 packets are used only when the “Message Timestamp” parameter in RS232 settings is enabled.
func DecodeCodec13(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if len(data) < 12 {
		return nil, fmt.Errorf("data length too short")
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	fmt.Printf("Number of Commands: %d\n", numberOfCommands)
	read += 1
	responseTypeNumber, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing response type: %w", err)
	}
	var responseType string
	switch responseTypeNumber {
	case 5:
		responseType = "Command"
		return nil, fmt.Errorf("codec 13 does not support command type: %s (Only supports Response)", responseType)
	case 6:
		responseType = "Response"
	default:
		return nil, fmt.Errorf("unknown response type: %d", responseTypeNumber)
	}
	fmt.Printf("Response Type: %d (%s)\n", responseTypeNumber, responseType)
	read += 1
	var command_responses []tools_domain.CommandResponse
	fmt.Printf("Read: %d\n", read)
	for range int(numberOfCommands) {
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		fmt.Printf("Response Size: %d | %s\n", responseSize, hex.EncodeToString(data[read:read+4]))
		read += 4
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %w", err)
		}
		fmt.Printf("Timestamp: %s | %s\n", timestamp, hex.EncodeToString(data[read:read+8]))
		read += 8
		commandSize := responseSize - 8
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(commandSize)], []byte{}, int(commandSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		fmt.Printf("Decoded Command Response: %s\n", message)
		fmt.Printf("Hex Command Response: %s\n", hexString)
		read += int(commandSize)
		commandResponse := &tools_domain.CommandResponse{
			Timestamp:  timestamp,
			Response:   message,
			HexMessage: hexString,
		}
		command_responses = append(command_responses, *commandResponse)
		fmt.Printf("Read after command response: %d\n", read)
	}
	fmt.Printf("Read: %d\n", read)
	fmt.Printf("Hex: %s\n", hex.EncodeToString(data[read-15:read]))
	responseType2, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing size of message: %w", err)
	}
	read += 1
	if numberOfCommands != responseType2 {
		fmt.Printf("Hex: %s\n", hex.EncodeToString(data[read:read+1]))
		return nil, fmt.Errorf("response type mismatch: %d != %d", numberOfCommands, responseType2)
	}
	if protocol == "TCP" {
		tram := append([]byte{0x0D}, data...)
		checkCRC := tools.IsValidTram(tram)
		fmt.Printf("CRC Check: %s\n", hex.EncodeToString(tram))
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfCommands,
		Records: []decoder_domain.Record{
			{
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}

func DecodeCodec14(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if len(data) < 12 {
		return nil, fmt.Errorf("data length too short")
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	fmt.Printf("Number of Commands: %d\n", numberOfCommands)
	read += 1
	responseTypeNumber, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing response type: %w", err)
	}
	var responseType string
	switch responseTypeNumber {
	case 5:
		responseType = "Command"
	case 6:
		responseType = "Response"
	default:
		return nil, fmt.Errorf("unknown response type: %d", responseTypeNumber)
	}
	fmt.Printf("Response Type: %d (%s)\n", responseTypeNumber, responseType)
	read += 1
	var command_responses []tools_domain.CommandResponse
	for range int(numberOfCommands) {
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		fmt.Printf("Response Size: %d\n", responseSize)
		read += 4
		imei, err := tools.DecodeIMEI(data[read : read+int(responseSize)])
		if err != nil {
			return nil, fmt.Errorf("error parsing IMEI: %w", err)
		}
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(responseSize)], []byte{}, int(responseSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		fmt.Printf("Decoded Command Response: %s\n", message)
		fmt.Printf("Hex Command Response: %s\n", hexString)
		read += int(responseSize)
		commandResponse := &tools_domain.CommandResponse{
			Response:    message,
			HexMessage:  hexString,
			IMEI:        imei,
			CommandType: responseType,
		}
		command_responses = append(command_responses, *commandResponse)
	}
	responseType2, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing size of message: %w", err)
	}
	read += 1
	if numberOfCommands != responseType2 {
		return nil, fmt.Errorf("response type mismatch")
	}
	if protocol == "TCP" {
		tram := append([]byte{0x0E}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfCommands,
		Records: []decoder_domain.Record{
			{
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}
