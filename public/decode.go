package teltonika

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	tools "github.com/danieljvsa/teltonika-go/tools"
)

// Decode parses a complete Teltonika frame (login or data, TCP or UDP)
// and returns a Packet containing only public types.
//
// This is the main entry point used by external applications.
func Decode(data []byte) (*Packet, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	if isLoginFrame(data) {
		return DecodeLogin(data)
	}
	return DecodeData(data)
}

// DecodeLogin is an alias for Decode that matches the legacy LoginDecoder
// naming. It parses a Teltonika IMEI login packet.
func DecodeLogin(data []byte) (*Packet, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("login frame too small")
	}
	length := binary.BigEndian.Uint16(data[0:2])
	if length == 0 {
		return nil, fmt.Errorf("login length is zero")
	}
	if int(length) != len(data)-2 {
		return nil, fmt.Errorf("login length mismatch: declared %d, actual %d", length, len(data)-2)
	}
	return &Packet{
		Kind:     KindLogin,
		Protocol: ProtocolTCP,
		IMEI:     string(data[2:]),
	}, nil
}

// isLoginFrame reports whether data looks like a Teltonika IMEI login packet:
// a 2-byte big-endian length followed by that many printable ASCII bytes.
func isLoginFrame(data []byte) bool {
	if len(data) < 3 {
		return false
	}
	length := binary.BigEndian.Uint16(data[0:2])
	if length == 0 || int(length) != len(data)-2 {
		return false
	}
	for _, b := range data[2:] {
		if b < 0x20 || b > 0x7E {
			return false
		}
	}
	return true
}

// DecodeData parses a complete AVL data packet (TCP or UDP transport)
// and returns a Packet containing only public types.
func DecodeData(data []byte) (*Packet, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	header, err := decodeHeader(data)
	if err != nil {
		return nil, err
	}
	codecIdx := header.lastByte
	if len(data) <= codecIdx {
		return nil, fmt.Errorf("missing codec id")
	}
	codec := CodecID(data[codecIdx])
	payload := data[codecIdx+1:]

	packet, err := DecodeCodecData(payload, codec, header.protocol)
	if err != nil {
		return nil, err
	}
	packet.Header = header.header
	return packet, nil
}

// DecodeCodecData decodes a codec payload (the bytes following the codec id,
// excluding any transport header) for the given codec and protocol.
//
// This is useful when a frame has already been split into its transport
// header, codec id and payload. For ProtocolTCP the trailing 4-byte CRC is
// validated.
func DecodeCodecData(data []byte, codec CodecID, protocol Protocol) (*Packet, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty payload")
	}

	packet := &Packet{
		Kind:     KindData,
		Protocol: protocol,
		Codec:    codec,
	}

	var err error
	switch codec {
	case Codec8:
		packet.Records, err = decodeAVL(codec, data, protocol, codec8IO())
	case Codec8Ext:
		packet.Records, err = decodeAVL(codec, data, protocol, codec8ExtIO())
	case Codec16:
		packet.Records, err = decodeAVL(codec, data, protocol, codec16IO())
	case Codec12:
		packet.Commands, err = decodeCommandResponses(data, protocol, false, false)
	case Codec13:
		packet.Commands, err = decodeCommandResponses(data, protocol, true, false)
	case Codec14:
		packet.Commands, err = decodeCommandResponses(data, protocol, false, true)
	case Codec15:
		packet.Commands, err = decodeCommandResponses(data, protocol, true, true)
	default:
		return nil, fmt.Errorf("unknown codec: 0x%02X", byte(codec))
	}

	return packet, err
}

// ioFormat describes the width of count/id fields for a codec's I/O block.
type ioFormat struct {
	idSize    int
	countSize int
	genType   bool
	hasXGroup bool
}

func codec8IO() ioFormat { return ioFormat{idSize: 1, countSize: 1} }
func codec8ExtIO() ioFormat {
	return ioFormat{idSize: 2, countSize: 2, hasXGroup: true}
}
func codec16IO() ioFormat { return ioFormat{idSize: 2, countSize: 1, genType: true} }

// decodeAVL decodes AVL records for codecs 08, 8E and 16.
// payload excludes the codec id byte but includes the trailing record count
// and (for TCP) the 4-byte CRC.
func decodeAVL(codec CodecID, payload []byte, protocol Protocol, format ioFormat) ([]AVLRecord, error) {
	if protocol == ProtocolTCP && !crcValid(byte(codec), payload) {
		return nil, fmt.Errorf("CRC is not valid")
	}
	if len(payload) < 1 {
		return nil, fmt.Errorf("data length too short")
	}

	numberOfRecords := int(payload[0])
	read := 1
	var records []AVLRecord

	for range numberOfRecords {
		if len(payload) < read+8 {
			return nil, fmt.Errorf("data length too short")
		}
		timestamp, err := tools.CalcTimestamp(payload[read : read+8])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %w", err)
		}
		read += 8

		if len(payload) < read+1 {
			return nil, fmt.Errorf("data length too short")
		}
		priority := int64(payload[read])
		read += 1

		if len(payload) < read+15 {
			return nil, fmt.Errorf("data length too short")
		}
		gps, err := decodeGPSData(payload[read : read+15])
		if err != nil {
			return nil, err
		}
		read += 15

		eventIOSize := 1
		if format.idSize == 2 {
			eventIOSize = 2
		}
		if len(payload) < read+eventIOSize {
			return nil, fmt.Errorf("data length too short")
		}
		eventIO := int64(payload[read])
		if eventIOSize == 2 {
			eventIO = int64(binary.BigEndian.Uint16(payload[read : read+2]))
		}
		read += eventIOSize

		generationType := ""
		if format.genType {
			if len(payload) < read+1 {
				return nil, fmt.Errorf("data length too short")
			}
			generationType, err = tools.GetGenerationType(payload, int64(read), 1)
			if err != nil {
				return nil, err
			}
			read += 1
		}

		ios, next, err := decodeIOElements(payload, read, format)
		if err != nil {
			return nil, err
		}
		read = next

		records = append(records, AVLRecord{
			Timestamp:      *timestamp,
			Priority:       priority,
			GPS:            gps,
			EventIO:        eventIO,
			IOElements:     ios,
			GenerationType: generationType,
		})
	}

	// All bytes after the records are the trailing record count and, for TCP
	// frames, the four-byte CRC. Reject anything else.
	expectedRemainder := 1
	if protocol == ProtocolTCP {
		expectedRemainder = 5
	}
	if remaining := len(payload) - read; remaining != expectedRemainder {
		return nil, fmt.Errorf("unexpected trailing data: %d bytes", remaining)
	}
	if trailing := int(payload[read]); trailing != numberOfRecords {
		return nil, fmt.Errorf("record count mismatch: initial %d, trailing %d", numberOfRecords, trailing)
	}

	return records, nil
}

// decodeIOElements parses the I/O element block of a single AVL record.
// It returns the decoded elements and the next read offset.
func decodeIOElements(data []byte, read int, format ioFormat) ([]IOElement, int, error) {
	var ios []IOElement

	readCount := func() (int64, error) {
		if len(data) < read+format.countSize {
			return 0, fmt.Errorf("data length too short")
		}
		var count int64
		if format.countSize == 2 {
			count = int64(binary.BigEndian.Uint16(data[read : read+2]))
		} else {
			count = int64(data[read])
		}
		read += format.countSize
		return count, nil
	}

	readID := func() (int64, error) {
		if len(data) < read+format.idSize {
			return 0, fmt.Errorf("data length too short")
		}
		var id int64
		if format.idSize == 2 {
			id = int64(binary.BigEndian.Uint16(data[read : read+2]))
		} else {
			id = int64(data[read])
		}
		read += format.idSize
		return id, nil
	}

	if _, err := readCount(); err != nil { // total I/O count
		return nil, read, err
	}

	for _, valueSize := range []int{1, 2, 4, 8} {
		count, err := readCount()
		if err != nil {
			return nil, read, err
		}
		for range int(count) {
			id, err := readID()
			if err != nil {
				return nil, read, err
			}
			if len(data) < read+valueSize {
				return nil, read, fmt.Errorf("data length too short")
			}
			value := hexEncode(data[read : read+valueSize])
			read += valueSize
			ios = append(ios, IOElement{ID: id, Value: value})
		}
	}

	if format.hasXGroup {
		count, err := readCount()
		if err != nil {
			return nil, read, err
		}
		for range int(count) {
			id, err := readID()
			if err != nil {
				return nil, read, err
			}
			length, err := readCount()
			if err != nil {
				return nil, read, err
			}
			if len(data) < read+int(length) {
				return nil, read, fmt.Errorf("data length too short")
			}
			value := hexEncode(data[read : read+int(length)])
			read += int(length)
			ios = append(ios, IOElement{ID: id, Value: value})
		}
	}

	return ios, read, nil
}

// decodeCommandResponses decodes command response codecs 12, 13, 14 and 15.
func decodeCommandResponses(payload []byte, protocol Protocol, withTimestamp bool, withIMEI bool) ([]Command, error) {
	if protocol == ProtocolTCP && !crcValid(codecFor(withTimestamp, withIMEI), payload) {
		return nil, fmt.Errorf("CRC is not valid")
	}

	numberOfCommands := int(payload[0])
	read := 1
	if len(payload) < read+1 {
		return nil, fmt.Errorf("data length too short")
	}
	responseTypeNumber := int64(payload[read])
	read += 1

	var commandType string
	switch responseTypeNumber {
	case 5:
		commandType = "Command"
		if withTimestamp && !withIMEI {
			return nil, fmt.Errorf("codec 13 does not support command type: %s", commandType)
		}
	case 6:
		commandType = "Response"
	case 0x0B:
		// Codec 15 uses a raw response type byte (legacy decoders accept 0x0B).
		commandType = fmt.Sprintf("%d", responseTypeNumber)
	default:
		if withTimestamp && withIMEI {
			// Legacy Codec 15 decoder accepts any response type value.
			commandType = fmt.Sprintf("%d", responseTypeNumber)
			break
		}
		return nil, fmt.Errorf("unknown response type: %d", responseTypeNumber)
	}

	var responses []CommandResponse
	for range numberOfCommands {
		if len(payload) < read+4 {
			return nil, fmt.Errorf("data length too short")
		}
		responseSize := int64(binary.BigEndian.Uint32(payload[read : read+4]))
		read += 4
		if responseSize < 0 {
			return nil, fmt.Errorf("invalid response size")
		}

		minimum := int64(0)
		switch {
		case withTimestamp && withIMEI:
			// Codec 15: 4-byte seconds timestamp + 8-byte IMEI.
			minimum = 4 + 8
		case withTimestamp:
			// Codec 13: 8-byte millisecond timestamp.
			minimum = 8
		case withIMEI:
			// Codec 14: response embeds an IMEI-sized message.
			minimum = 8
		}
		if responseSize < minimum {
			return nil, fmt.Errorf("response size too small")
		}

		cmd := CommandResponse{CommandType: commandType}

		if withTimestamp {
			if withIMEI {
				// Timestamp is 4 bytes (seconds) for codec 15.
				if len(payload) < read+4 {
					return nil, fmt.Errorf("data length too short")
				}
				timestamp, err := tools.CalcTimestampSecondsBigEndian(payload[read : read+4])
				if err != nil {
					return nil, fmt.Errorf("error parsing timestamp: %w", err)
				}
				cmd.Timestamp = timestamp
				read += 4

				if len(payload) < read+8 {
					return nil, fmt.Errorf("data length too short")
				}
				cmd.IMEI, err = decodeIMEI(payload[read : read+8])
				if err != nil {
					return nil, fmt.Errorf("error parsing IMEI: %w", err)
				}
				read += 8
				commandSize := responseSize - 12
				if err := decodeCommandMessage(payload, read, commandSize, &cmd); err != nil {
					return nil, err
				}
				read += int(commandSize)
			} else {
				// Codec 13: 8-byte millisecond timestamp.
				if len(payload) < read+8 {
					return nil, fmt.Errorf("data length too short")
				}
				timestamp, err := tools.CalcTimestamp(payload[read : read+8])
				if err != nil {
					return nil, fmt.Errorf("error parsing timestamp: %w", err)
				}
				cmd.Timestamp = timestamp
				read += 8
				commandSize := responseSize - 8
				if err := decodeCommandMessage(payload, read, commandSize, &cmd); err != nil {
					return nil, err
				}
				read += int(commandSize)
			}
		} else if withIMEI {
			// Codec 14: response is an 8-byte IMEI followed by the command message.
			if len(payload) < read+8 {
				return nil, fmt.Errorf("data length too short")
			}
			imei, err := decodeIMEI(payload[read : read+8])
			if err != nil {
				return nil, fmt.Errorf("error parsing IMEI: %w", err)
			}
			cmd.IMEI = imei
			read += 8
			commandSize := int64(responseSize) - 8
			if err := decodeCommandMessage(payload, read, commandSize, &cmd); err != nil {
				return nil, err
			}
			read += int(commandSize)
		} else {
			// Codec 12: just the command message.
			if err := decodeCommandMessage(payload, read, responseSize, &cmd); err != nil {
				return nil, err
			}
			read += int(responseSize)
		}

		responses = append(responses, cmd)
	}

	if len(payload) < read+1 {
		return nil, fmt.Errorf("data length too short")
	}
	trailing := int64(payload[read])
	read += 1
	if int64(numberOfCommands) != trailing {
		return nil, fmt.Errorf("response type mismatch: %d != %d", numberOfCommands, trailing)
	}

	return []Command{{Type: commandType, Responses: responses}}, nil
}

func decodeCommandMessage(payload []byte, read int, size int64, cmd *CommandResponse) error {
	if len(payload) < read+int(size) {
		return fmt.Errorf("data length too short")
	}
	hexString, message, err := tools.DecodeToHexThenASCII(payload[read:read+int(size)], []byte{}, int(size))
	if err != nil {
		return fmt.Errorf("error parsing message: %w", err)
	}
	cmd.Response = message
	cmd.HexMessage = hexString
	return nil
}

// decodeIMEI decodes an IMEI byte slice to its hex string representation.
func decodeIMEI(data []byte) (string, error) {
	return tools.DecodeIMEI(data)
}

// IOBlock is the result of decoding a standalone I/O element block.
type IOBlock struct {
	Elements       []IOElement
	GenerationType string
	NextOffset     int
}

// DecodeCodecIO decodes a standalone I/O element block, optionally reading a
// Codec 16 generation type byte first. idSize is the I/O id width (1 for
// Codec 8, 2 for Codec 8E and 16); countSize is the count/group width (1 for
// Codecs 8 and 16, 2 for Codec 8E). hasXGroup enables the variable-length
// group of the extended Codec 8E format.
func DecodeCodecIO(data []byte, startByte int, idSize int, withGen bool, hasXGroup bool) (IOBlock, error) {
	format := ioFormat{idSize: idSize, countSize: 1, hasXGroup: hasXGroup}
	if hasXGroup {
		format.countSize = 2
	}

	read := startByte
	block := IOBlock{}

	if withGen {
		if len(data) < read+1 {
			return block, fmt.Errorf("data too short for generation type")
		}
		genType, err := tools.GetGenerationType(data, int64(read), 1)
		if err != nil {
			return block, err
		}
		block.GenerationType = genType
		read += 1
	}

	elements, next, err := decodeIOElements(data, read, format)
	if err != nil {
		return block, err
	}
	block.Elements = elements
	block.NextOffset = next
	return block, nil
}

// hexEncode returns the lowercase hexadecimal encoding of data.
func hexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// codecFor returns the command codec id given its variant flags.
func codecFor(withTimestamp bool, withIMEI bool) byte {
	switch {
	case withTimestamp && withIMEI:
		return 0x0F
	case withTimestamp:
		return 0x0D
	case withIMEI:
		return 0x0E
	default:
		return 0x0C
	}
}
