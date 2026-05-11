package teltonika_go

import (
	"encoding/hex"
	"fmt"
	"strconv"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	tools_domain "github.com/danieljvsa/teltonika-go/internal/tool"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func ensureRead(data []byte, read int, need int) error {
	if need < 0 {
		return fmt.Errorf("invalid length requested")
	}
	if read+need > len(data) {
		return fmt.Errorf("data length too short")
	}
	return nil
}

// DecodeCodec8 decodes Teltonika Codec 08 (AVL data) frames containing GPS records with I/O data.
// Codec 08 is the most common protocol for transmitting vehicle location and telemetry.
//
// Frame format per record:
//   - Timestamp: 8 bytes (milliseconds since Unix epoch)
//   - Priority: 1 byte
//   - GPS Data: 15 bytes (latitude, longitude, altitude, angle, satellites, speed)
//   - Event IO: 1 byte
//   - IO Data: Variable length (depends on IO element count)
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing all decoded records
//   - error: if frame is malformed or CRC check fails
//
// Example:
//
//	codecData, err := DecodeCodec8(frameData, "TCP")
//	if err == nil {
//		for _, record := range codecData.Records {
//			fmt.Printf("Lat: %f, Lon: %f\n", record.GPSData.Latitude, record.GPSData.Longitude)
//		}
//	}
func DecodeCodec8(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []decoder_domain.Record
	for range int(numberOfRecords) {
		if err := ensureRead(data, read, 8); err != nil {
			return nil, err
		}
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp")
		}
		read += 8
		if err := ensureRead(data, read, 1); err != nil {
			return nil, err
		}
		priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing priority: %w", err)
		}
		read += 1
		if err := ensureRead(data, read, 15); err != nil {
			return nil, err
		}
		gpsData, err := tools.DecodeGPSData(data[read : read+15])
		if err != nil {
			return nil, fmt.Errorf("error parsing GPS data")
		}
		read += 15
		if err := ensureRead(data, read, 1); err != nil {
			return nil, err
		}
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

// DecodeCodec8Ext decodes Teltonika Codec 8E (Extended AVL data) frames.
// Codec 8E is an extended version of Codec 8 with additional field support and
// enhanced structure for more detailed telemetry.
//
// Similar to Codec 8 but with extended structure and additional data fields.
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing all decoded records
//   - error: if frame is malformed or CRC check fails
//
// Example:
//
//	codecData, err := DecodeCodec8Ext(frameData, "TCP")
func DecodeCodec8Ext(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []decoder_domain.Record
	for range int(numberOfRecords) {
		if err := ensureRead(data, read, 8); err != nil {
			return nil, err
		}
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp")
		}
		read += 8
		if err := ensureRead(data, read, 1); err != nil {
			return nil, err
		}
		priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing priority: %w", err)
		}
		read += 1
		if err := ensureRead(data, read, 15); err != nil {
			return nil, err
		}
		gpsData, err := tools.DecodeGPSData(data[read : read+15])
		if err != nil {
			return nil, fmt.Errorf("error parsing GPS data")
		}
		read += 15
		if err := ensureRead(data, read, 2); err != nil {
			return nil, err
		}
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

// DecodeCodec16 decodes Teltonika Codec 16 (GPRS/IQ frames with device type).
// Codec 16 includes device type information and specialized handling for
// GPRS-based communication frames.
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing all decoded records
//   - error: if frame is malformed or CRC check fails
//
// Example:
//
//	codecData, err := DecodeCodec16(frameData, "TCP")
func DecodeCodec16(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}

	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	var records []decoder_domain.Record
	for range int(numberOfRecords) {
		if err := ensureRead(data, read, 8); err != nil {
			return nil, err
		}
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp")
		}
		read += 8
		if err := ensureRead(data, read, 1); err != nil {
			return nil, err
		}
		priority, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing priority: %w", err)
		}
		read += 1
		if err := ensureRead(data, read, 15); err != nil {
			return nil, err
		}
		gpsData, err := tools.DecodeGPSData(data[read : read+15])
		if err != nil {
			return nil, fmt.Errorf("error parsing GPS data")
		}
		read += 15
		if err := ensureRead(data, read, 2); err != nil {
			return nil, err
		}
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

// DecodeCodec12 decodes Teltonika Codec 12 (Command response codec).
// Codec 12 handles command responses from devices with command handling and
// response data including timestamps and IMEI information.
//
// Frame structure:
//   - Number of commands: 1 byte
//   - Response type: 1 byte (5=Command, 6=Response)
//   - For each command:
//   - Response size: 4 bytes
//   - Timestamp: 4 bytes (seconds)
//   - IMEI: 8 bytes
//   - Command response: variable length
//   - Response type count: 1 byte
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing command responses
//   - error: if frame is malformed or validation fails
//
// Example:
//
//	codecData, err := DecodeCodec12(frameData, "TCP")
//	if err == nil {
//		for _, record := range codecData.Records {
//			if record.CommandResponses != nil {
//				for _, cmd := range *record.CommandResponses {
//					fmt.Println("Response:", cmd.Response)
//				}
//			}
//		}
//	}
func DecodeCodec12(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 2); err != nil {
		return nil, err
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
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
	read += 1
	var command_responses []tools_domain.CommandResponse
	for range int(numberOfCommands) {
		if err := ensureRead(data, read, 4); err != nil {
			return nil, err
		}
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		read += 4
		if responseSize < 0 {
			return nil, fmt.Errorf("invalid response size")
		}
		if err := ensureRead(data, read, int(responseSize)); err != nil {
			return nil, err
		}
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(responseSize)], []byte{}, int(responseSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		read += int(responseSize)
		commandResponse := &tools_domain.CommandResponse{
			Response:   message,
			HexMessage: hexString,
		}
		command_responses = append(command_responses, *commandResponse)
	}
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
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
				CommandType:      &responseType,
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}

// DecodeCodec13 decodes Teltonika Codec 13 (Command response codec).
// Codec 13 is similar to Codec 12 but includes timestamps for each command response.
//
// Frame structure:
//   - Number of commands: 1 byte
//   - Response type: 1 byte (6=Response) Other types not supported
//   - For each command:
//   - Response size: 4 bytes
//   - Timestamp: 8 bytes (milliseconds since Unix epoch)
//   - Command response: variable length
//   - Response type count: 1 byte
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing command responses with timestamps
//   - error: if frame is malformed or validation fails
//
// Example:
//
//	codecData, err := DecodeCodec13(frameData, "TCP")
//	if err == nil {
//		for _, record := range codecData.Records {
//			if record.CommandResponses != nil {
//				for _, cmd := range *record.CommandResponses {
//					fmt.Println("Timestamp:", cmd.Timestamp, "Response:", cmd.Response)
//				}
//			}
//		}
//	}
func DecodeCodec13(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 2); err != nil {
		return nil, err
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
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
	read += 1
	var command_responses []tools_domain.CommandResponse
	for range int(numberOfCommands) {
		if err := ensureRead(data, read, 4); err != nil {
			return nil, err
		}
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		read += 4
		if responseSize < 8 {
			return nil, fmt.Errorf("response size too small")
		}
		if err := ensureRead(data, read, 8); err != nil {
			return nil, err
		}
		timestamp, err := tools.CalcTimestamp(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %w", err)
		}
		read += 8
		commandSize := responseSize - 8
		if commandSize < 0 {
			return nil, fmt.Errorf("invalid command size")
		}
		if err := ensureRead(data, read, int(commandSize)); err != nil {
			return nil, err
		}
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(commandSize)], []byte{}, int(commandSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		read += int(commandSize)
		commandResponse := &tools_domain.CommandResponse{
			Timestamp:  timestamp,
			Response:   message,
			HexMessage: hexString,
		}
		command_responses = append(command_responses, *commandResponse)
	}
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
	responseType2, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing size of message: %w", err)
	}
	read += 1
	if numberOfCommands != responseType2 {
		return nil, fmt.Errorf("response type mismatch: %d != %d", numberOfCommands, responseType2)
	}
	if protocol == "TCP" {
		tram := append([]byte{0x0D}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfCommands,
		Records: []decoder_domain.Record{
			{
				CommandType:      &responseType,
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}

// DecodeCodec14 decodes Teltonika Codec 14 (Command response codec).
// Codec 14 handles command responses with IMEI information.
// Frame structure:
//   - Number of commands: 1 byte
//   - Response type: 1 byte (5=Command, 6=Response)
//   - For each command:
//   - Response size: 4 bytes
//   - IMEI: 8 bytes
//   - Command response: variable length
//   - Response type count: 1 byte
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing command responses
//   - error: if frame is malformed or validation fails
//
// Example:
//
//	codecData, err := DecodeCodec14(frameData, "TCP")
//	if err == nil {
//		for _, record := range codecData.Records {
//			if record.CommandResponses != nil {
//				for _, cmd := range *record.CommandResponses {
//					fmt.Println("Response:", cmd.Response)
//				}
//			}
//		}
//	}
func DecodeCodec14(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 2); err != nil {
		return nil, err
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
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
	read += 1
	var command_responses []tools_domain.CommandResponse
	for range int(numberOfCommands) {
		if err := ensureRead(data, read, 4); err != nil {
			return nil, err
		}
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		read += 4
		if responseSize < 8 {
			return nil, fmt.Errorf("response size too small")
		}
		if err := ensureRead(data, read, int(responseSize)); err != nil {
			return nil, err
		}
		imei, err := tools.DecodeIMEI(data[read : read+int(responseSize)])
		if err != nil {
			return nil, fmt.Errorf("error parsing IMEI: %w", err)
		}
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(responseSize)], []byte{}, int(responseSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		read += int(responseSize)
		commandResponse := &tools_domain.CommandResponse{
			Response:    message,
			HexMessage:  hexString,
			IMEI:        imei,
			CommandType: responseType,
		}
		command_responses = append(command_responses, *commandResponse)
	}
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
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
				CommandType:      &responseType,
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}

// DecodeCodec15 decodes Teltonika Codec 15 (Command response codec).
// Codec 15 handles command responses with timestamps and IMEI information.
//
// Frame structure:
//   - Number of commands: 1 byte
//   - Response type: 1 byte (5=Command, 6=Response)
//   - For each command:
//   - Response size: 4 bytes
//   - Timestamp: 4 bytes (seconds since Unix epoch)
//   - IMEI: 8 bytes
//   - Command response: variable length
//   - Response type count: 1 byte
//
// Parameters:
//   - data: decoded frame data without header/CRC
//   - protocol: "TCP" or "UDP" - determines whether to validate CRC
//
// Returns:
//   - *decoder_domain.CodecData: structure containing command responses
//   - error: if frame is malformed or validation fails
//
// Example:
//
//	codecData, err := DecodeCodec15(frameData, "TCP")
//	if err == nil {
//		for _, record := range codecData.Records {
//			if record.CommandResponses != nil {
//				for _, cmd := range *record.CommandResponses {
//					fmt.Println("Response:", cmd.Response)
//				}
//			}
//		}
//	}

func DecodeCodec15(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	read := 0
	if err := ensureRead(data, read, 2); err != nil {
		return nil, err
	}
	numberOfCommands, err := strconv.ParseInt(hex.EncodeToString(data[:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of records: %w", err)
	}
	read += 1
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
	responseTypeNumber, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing response type: %w", err)
	}
	responseType := fmt.Sprintf("%d", responseTypeNumber)
	read += 1
	var command_responses []tools_domain.CommandResponse
	for range int(numberOfCommands) {
		if err := ensureRead(data, read, 4); err != nil {
			return nil, err
		}
		responseSize, err := strconv.ParseInt(hex.EncodeToString(data[read:read+4]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing size of message: %w", err)
		}
		read += 4
		if responseSize < 12 {
			return nil, fmt.Errorf("response size too small")
		}
		if err := ensureRead(data, read, 4); err != nil {
			return nil, err
		}
		timestamp, err := tools.CalcTimestampSeconds(data[read : read+4])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %w", err)
		}
		read += 4
		if err := ensureRead(data, read, 8); err != nil {
			return nil, err
		}
		imei, err := tools.DecodeIMEI(data[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing IMEI: %w", err)
		}
		read += 8
		commandSize := responseSize - 12
		if commandSize < 0 {
			return nil, fmt.Errorf("invalid command size")
		}
		if err := ensureRead(data, read, int(commandSize)); err != nil {
			return nil, err
		}
		hexString, message, err := tools.DecodeToHexThenASCII(data[read:read+int(commandSize)], []byte{}, int(commandSize))
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}
		read += int(commandSize)
		commandResponse := &tools_domain.CommandResponse{
			Timestamp:  timestamp,
			Response:   message,
			HexMessage: hexString,
			IMEI:       imei,
		}
		command_responses = append(command_responses, *commandResponse)
	}
	if err := ensureRead(data, read, 1); err != nil {
		return nil, err
	}
	responseType2, err := strconv.ParseInt(hex.EncodeToString(data[read:read+1]), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing size of message: %w", err)
	}
	read += 1
	if numberOfCommands != responseType2 {
		return nil, fmt.Errorf("response type mismatch: %d != %d", numberOfCommands, responseType2)
	}
	if protocol == "TCP" {
		tram := append([]byte{0x0F}, data...)
		checkCRC := tools.IsValidTram(tram)
		if !checkCRC {
			return nil, fmt.Errorf("CRC is not valid")
		}
	}
	decodedData := &decoder_domain.CodecData{
		NumberOfRecords: numberOfCommands,
		Records: []decoder_domain.Record{
			{
				CommandType:      &responseType,
				CommandResponses: &command_responses,
			},
		},
	}
	return decodedData, nil
}
