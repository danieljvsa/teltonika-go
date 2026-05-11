package teltonika_go

import (
	"bytes"
	"encoding/binary"
	"fmt"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
	tools_domain "github.com/danieljvsa/teltonika-go/internal/tool"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func EncodeCodec8(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeAVLCodec(codecData, 0x08)
}

func EncodeCodec8Ext(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeAVLCodec(codecData, 0x8E)
}

func EncodeCodec12(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCommandCodec(codecData, 0x0C)
}

func EncodeCodec13(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCommandCodec(codecData, 0x0D)
}

func EncodeCodec14(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCommandCodec(codecData, 0x0E)
}

func EncodeCodec15(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCommandCodec(codecData, 0x0F)
}

func EncodeCodec16(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeAVLCodec(codecData, 0x10)
}

func encodeAVLCodec(codecData *decoder_domain.CodecData, codecID byte) ([]byte, error) {
	if codecData == nil {
		return nil, fmt.Errorf("codec data is nil")
	}
	if len(codecData.Records) == 0 {
		return nil, fmt.Errorf("no records to encode")
	}
	if codecData.NumberOfRecords > 0 && codecData.NumberOfRecords != int64(len(codecData.Records)) {
		return nil, fmt.Errorf("number of records mismatch")
	}

	buffer := &bytes.Buffer{}
	buffer.WriteByte(byte(len(codecData.Records)))

	for _, record := range codecData.Records {
		if record.Timestamp == nil {
			return nil, fmt.Errorf("record timestamp is required")
		}
		if record.Priority == nil {
			return nil, fmt.Errorf("record priority is required")
		}
		if record.GPSData == nil {
			return nil, fmt.Errorf("record GPS data is required")
		}
		if record.EventIO == nil {
			return nil, fmt.Errorf("record event IO is required")
		}

		timestampBytes, err := tools.EncodeTimestampMillis(record.Timestamp)
		if err != nil {
			return nil, err
		}
		buffer.Write(timestampBytes)
		if *record.Priority < 0 || *record.Priority > 255 {
			return nil, fmt.Errorf("record priority out of range: %d", *record.Priority)
		}
		buffer.WriteByte(byte(*record.Priority))

		gpsBytes, err := tools.EncodeGPSData(record.GPSData)
		if err != nil {
			return nil, err
		}
		buffer.Write(gpsBytes)
		buffer.WriteByte(0x00)

		switch codecID {
		case 0x08:
			if *record.EventIO < 0 || *record.EventIO > 255 {
				return nil, fmt.Errorf("event IO out of range for codec 8: %d", *record.EventIO)
			}
			buffer.WriteByte(byte(*record.EventIO))
			ioBytes, _, err := tools.EncodeIOData8(resolveRecordIOs(record))
			if err != nil {
				return nil, err
			}
			buffer.Write(ioBytes)
		case 0x8E:
			if *record.EventIO < 0 || *record.EventIO > 65535 {
				return nil, fmt.Errorf("event IO out of range for codec 8E: %d", *record.EventIO)
			}
			eventBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(eventBytes, uint16(*record.EventIO))
			buffer.Write(eventBytes)
			ioBytes, _, err := tools.EncodeIOData8Extended(resolveRecordIOs(record))
			if err != nil {
				return nil, err
			}
			buffer.Write(ioBytes)
		case 0x10:
			if *record.EventIO < 0 || *record.EventIO > 65535 {
				return nil, fmt.Errorf("event IO out of range for codec 16: %d", *record.EventIO)
			}
			eventBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(eventBytes, uint16(*record.EventIO))
			buffer.Write(eventBytes)
			generationType, err := resolveGenerationType(record)
			if err != nil {
				return nil, err
			}
			ioBytes, _, err := tools.EncodeIOData16(resolveRecordIOs(record), generationType)
			if err != nil {
				return nil, err
			}
			buffer.Write(ioBytes)
		default:
			return nil, fmt.Errorf("unsupported codec: 0x%X", codecID)
		}
	}

	buffer.WriteByte(byte(len(codecData.Records)))

	crcData := append([]byte{codecID}, buffer.Bytes()...)
	return tools.AppendCRC16IBM(crcData)[1:], nil
}

func encodeCommandCodec(codecData *decoder_domain.CodecData, codecID byte) ([]byte, error) {
	if codecData == nil {
		return nil, fmt.Errorf("codec data is nil")
	}

	commandResponses, commandType, err := resolveCommandResponses(codecData)
	if err != nil {
		return nil, err
	}
	if len(commandResponses) == 0 {
		return nil, fmt.Errorf("no command responses to encode")
	}

	if codecData.NumberOfRecords > 0 && codecData.NumberOfRecords != int64(len(commandResponses)) {
		return nil, fmt.Errorf("number of commands mismatch")
	}

	buffer := &bytes.Buffer{}
	buffer.WriteByte(byte(len(commandResponses)))

	responseType, err := resolveResponseType(codecID, commandType)
	if err != nil {
		return nil, err
	}
	buffer.WriteByte(responseType)

	for _, response := range commandResponses {
		commandBytes, err := tools.EncodeHexMessage(response.Response, response.HexMessage)
		if err != nil {
			return nil, err
		}

		var payload []byte
		switch codecID {
		case 0x0C:
			payload = commandBytes
		case 0x0D:
			if response.Timestamp == nil {
				return nil, fmt.Errorf("codec 13 requires timestamp in command response")
			}
			timestampBytes, err := tools.EncodeTimestampMillis(response.Timestamp)
			if err != nil {
				return nil, err
			}
			payload = append(timestampBytes, commandBytes...)
		case 0x0E:
			imeiBytes, err := tools.EncodeIMEI(response.IMEI)
			if err != nil {
				return nil, err
			}
			payload = append(imeiBytes, commandBytes...)
		case 0x0F:
			if response.Timestamp == nil {
				return nil, fmt.Errorf("codec 15 requires timestamp in command response")
			}
			timestampBytes, err := tools.EncodeTimestampSeconds(response.Timestamp)
			if err != nil {
				return nil, err
			}
			imeiBytes, err := tools.EncodeIMEI(response.IMEI)
			if err != nil {
				return nil, err
			}
			payload = append(timestampBytes, imeiBytes...)
			payload = append(payload, commandBytes...)
		default:
			return nil, fmt.Errorf("unsupported command codec: 0x%X", codecID)
		}

		if len(payload) > int(^uint32(0)) {
			return nil, fmt.Errorf("payload too large for codec %X", codecID)
		}
		sizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sizeBytes, uint32(len(payload)))
		buffer.Write(sizeBytes)
		buffer.Write(payload)
	}

	buffer.WriteByte(byte(len(commandResponses)))

	crcData := append([]byte{codecID}, buffer.Bytes()...)
	return tools.AppendCRC16IBM(crcData)[1:], nil
}

func resolveRecordIOs(record decoder_domain.Record) []io_domain.IOData {
	if record.IOs == nil {
		return []io_domain.IOData{}
	}
	return *record.IOs
}

func resolveGenerationType(record decoder_domain.Record) (string, error) {
	if record.Attributes == nil {
		return "", fmt.Errorf("codec 16 requires generation type in record attributes")
	}
	value, ok := (*record.Attributes)["generation_type"]
	if !ok {
		return "", fmt.Errorf("codec 16 requires generation_type attribute")
	}
	genType, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("generation_type must be a string")
	}
	return genType, nil
}

func resolveCommandResponses(codecData *decoder_domain.CodecData) ([]tools_domain.CommandResponse, string, error) {
	if len(codecData.Records) == 0 {
		return nil, "", fmt.Errorf("no records to encode")
	}
	record := codecData.Records[0]
	if record.CommandResponses == nil {
		return nil, "", fmt.Errorf("record command responses are required")
	}
	commandResponses := *record.CommandResponses
	commandType := ""
	if record.CommandType != nil {
		commandType = *record.CommandType
	} else if len(commandResponses) > 0 {
		commandType = commandResponses[0].CommandType
	}
	return commandResponses, commandType, nil
}

func resolveResponseType(codecID byte, commandType string) (byte, error) {
	switch commandType {
	case "Command", "5":
		if codecID == 0x0D {
			return 0, fmt.Errorf("codec 13 only supports response type")
		}
		return 5, nil
	case "Response", "6":
		return 6, nil
	case "":
		if codecID == 0x0D {
			return 6, nil
		}
		return 6, nil
	default:
		return 0, fmt.Errorf("unknown command type: %s", commandType)
	}
}
