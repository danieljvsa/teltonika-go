package teltonika

import (
	"bytes"
	"encoding/binary"
	"fmt"

	tools "github.com/danieljvsa/teltonika-go/tools"
)

// Encode serializes a Packet into a complete, wire-ready Teltonika frame
// (including transport header, codec id and CRC for TCP).
//
// The returned frame can be fed back to Decode, i.e. Decode(Encode(p)) == p.
func Encode(packet *Packet) ([]byte, error) {
	if packet == nil {
		return nil, fmt.Errorf("packet is nil")
	}

	switch packet.Kind {
	case KindLogin:
		return encodeLogin(packet)
	case KindData:
		return encodeDataFrame(packet)
	default:
		return nil, fmt.Errorf("unknown packet kind: %d", kind(packet.Kind))
	}
}

func encodeLogin(packet *Packet) ([]byte, error) {
	imei := []byte(packet.IMEI)
	if len(imei) == 0 {
		return nil, fmt.Errorf("login IMEI is empty")
	}
	if len(imei) > 65535 {
		return nil, fmt.Errorf("IMEI too long")
	}
	frame := make([]byte, 2+len(imei))
	binary.BigEndian.PutUint16(frame[0:2], uint16(len(imei)))
	copy(frame[2:], imei)
	return frame, nil
}

func encodeDataFrame(packet *Packet) ([]byte, error) {
	var payload []byte
	switch packet.Codec {
	case Codec8, Codec8Ext, Codec16:
		data, err := encodeAVL(packet.Codec, packet.Records)
		if err != nil {
			return nil, err
		}
		payload = data
	case Codec12, Codec13, Codec14, Codec15:
		data, err := encodeCommand(packet.Codec, packet.Commands)
		if err != nil {
			return nil, err
		}
		payload = data
	default:
		return nil, fmt.Errorf("unsupported codec: 0x%02X", byte(packet.Codec))
	}

	return wrapFrame(packet, byte(packet.Codec), payload)
}

// wrapFrame assembles the transport header, codec id and payload into a
// complete frame, appending the CRC for TCP packets.
func wrapFrame(packet *Packet, codec byte, payload []byte) ([]byte, error) {
	protocol := packet.Protocol
	if protocol == "" {
		protocol = ProtocolTCP
	}

	if protocol == ProtocolTCP {
		// frame = header(8) + codec + payload(with CRC)
		payload = appendCRC(append([]byte{codec}, payload...))
		frame := make([]byte, 8, 8+len(payload))
		binary.BigEndian.PutUint32(frame[4:8], uint32(len(payload)))
		return append(frame, payload...), nil
	}

	// UDP frame = length(2) packetId(2) version(1) avlPacketId(1)
	//            imeiLen(2) imei + codec + payload (no CRC)
	h := packet.Header.UDP
	imei := []byte{}
	imeiLen := uint16(0)
	packetID := uint16(0)
	avlPacketID := byte(0)
	if h != nil {
		imei = []byte(h.IMEI)
		imeiLen = uint16(h.IMEILength)
		if imeiLen == 0 && len(h.IMEI) > 0 {
			imeiLen = uint16(len(h.IMEI))
		}
		packetID = uint16(h.PacketID)
		avlPacketID = byte(h.AVLPacketID)
	}

	mid := make([]byte, 6, 6+len(imei))
	binary.BigEndian.PutUint16(mid[0:2], packetID)
	mid[2] = 0x01 // version byte
	mid[3] = avlPacketID
	binary.BigEndian.PutUint16(mid[4:6], imeiLen)
	mid = append(mid, imei...)

	codecPart := append([]byte{codec}, payload...)
	length := len(mid) + len(codecPart)
	if length > 65535 {
		return nil, fmt.Errorf("UDP frame too large")
	}

	frame := make([]byte, 2, 2+length)
	binary.BigEndian.PutUint16(frame[0:2], uint16(length))
	frame = append(frame, mid...)
	frame = append(frame, codecPart...)
	return frame, nil
}

// encodeAVL serializes AVL records into the codec payload (excluding the
// codec id byte; trailing record count plus the CRC are added by wrapFrame).
func encodeAVL(codec CodecID, records []AVLRecord) ([]byte, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no records to encode")
	}
	if len(records) > 255 {
		return nil, fmt.Errorf("record count exceeds uint8 range")
	}

	var format ioFormat
	switch codec {
	case Codec8:
		format = codec8IO()
	case Codec8Ext:
		format = codec8ExtIO()
	case Codec16:
		format = codec16IO()
	default:
		return nil, fmt.Errorf("unsupported AVL codec: 0x%02X", byte(codec))
	}

	buffer := &bytes.Buffer{}
	buffer.WriteByte(byte(len(records)))

	for _, record := range records {
		tsBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(tsBytes, uint64(record.Timestamp.UTC().UnixMilli()))
		buffer.Write(tsBytes)

		if record.Priority < 0 || record.Priority > 255 {
			return nil, fmt.Errorf("record priority out of range: %d", record.Priority)
		}
		buffer.WriteByte(byte(record.Priority))

		gps, err := encodeGPSData(record.GPS)
		if err != nil {
			return nil, err
		}
		buffer.Write(gps)

		if format.idSize == 2 {
			if record.EventIO < 0 || record.EventIO > 65535 {
				return nil, fmt.Errorf("event IO out of range: %d", record.EventIO)
			}
			eventBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(eventBytes, uint16(record.EventIO))
			buffer.Write(eventBytes)
		} else {
			if record.EventIO < 0 || record.EventIO > 255 {
				return nil, fmt.Errorf("event IO out of range: %d", record.EventIO)
			}
			buffer.WriteByte(byte(record.EventIO))
		}

		if format.genType {
			_ = record.GenerationType // validated inside encodeIOElements
		}

		ioBytes, err := encodeIOElements(record.IOElements, format, record.GenerationType)
		if err != nil {
			return nil, err
		}
		buffer.Write(ioBytes)
	}

	buffer.WriteByte(byte(len(records)))
	return buffer.Bytes(), nil
}

// encodeIOElements serializes IO elements in the width/order expected by the
// given codec's I/O block format.
func encodeIOElements(ios []IOElement, format ioFormat, generationType string) ([]byte, error) {
	var oneByte, twoByte, fourByte, eightByte, xByte []IOElement

	for _, io := range ios {
		valueBytes, err := decodeHex(io.Value)
		if err != nil {
			return nil, err
		}
		switch len(valueBytes) {
		case 1:
			oneByte = append(oneByte, io)
		case 2:
			twoByte = append(twoByte, io)
		case 4:
			fourByte = append(fourByte, io)
		case 8:
			eightByte = append(eightByte, io)
		default:
			xByte = append(xByte, io)
		}
	}

	if !format.hasXGroup && len(xByte) > 0 {
		return nil, fmt.Errorf("codec does not support variable-length IO values")
	}

	maxCount := 255
	if format.countSize == 2 {
		maxCount = 65535
	}
	if len(ios) > maxCount {
		return nil, fmt.Errorf("IO count exceeds range")
	}

	buffer := &bytes.Buffer{}
	if format.genType {
		genByte, err := tools.EncodeGenerationType(generationType)
		if err != nil {
			return nil, err
		}
		buffer.WriteByte(genByte)
	}

	writeCount := countWriter(buffer, format)

	writeCount(len(ios))

	for _, group := range [][]IOElement{oneByte, twoByte, fourByte, eightByte} {
		writeCount(len(group))
		for _, io := range group {
			if err := writeIO(io, format, false, buffer); err != nil {
				return nil, err
			}
		}
	}

	if format.hasXGroup {
		writeCount(len(xByte))
		for _, io := range xByte {
			if err := writeIO(io, format, true, buffer); err != nil {
				return nil, err
			}
		}
	}

	return buffer.Bytes(), nil
}

// writeIO serializes one I/O element. withLength writes a length prefix
// before the value (used by the 8E variable-length X group).
func writeIO(io IOElement, format ioFormat, withLength bool, buffer *bytes.Buffer) error {
	if io.ID < 0 || io.ID > int64(maxIDFor(format)) {
		return fmt.Errorf("IO ID out of range: %d", io.ID)
	}
	if format.idSize == 2 {
		idBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(idBytes, uint16(io.ID))
		buffer.Write(idBytes)
	} else {
		buffer.WriteByte(byte(io.ID))
	}
	valueBytes, err := decodeHex(io.Value)
	if err != nil {
		return err
	}
	if withLength {
		if len(valueBytes) > 65535 {
			return fmt.Errorf("IO value too large: %d bytes", len(valueBytes))
		}
		lengthBytes := make([]byte, format.countSize)
		if format.countSize == 2 {
			binary.BigEndian.PutUint16(lengthBytes, uint16(len(valueBytes)))
		} else {
			lengthBytes[0] = byte(len(valueBytes))
		}
		buffer.Write(lengthBytes)
	}
	buffer.Write(valueBytes)
	return nil
}

func maxIDFor(format ioFormat) int {
	if format.idSize == 2 {
		return 65535
	}
	return 255
}

// countWriter returns a writer that serializes counts in the codec's width.
func countWriter(buffer *bytes.Buffer, format ioFormat) func(int) {
	return func(count int) {
		if format.countSize == 2 {
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(count))
			buffer.Write(b)
		} else {
			buffer.WriteByte(byte(count))
		}
	}
}

func encodeCommand(codec CodecID, commands []Command) ([]byte, error) {
	if len(commands) == 0 {
		return nil, fmt.Errorf("no commands to encode")
	}
	if len(commands) > 255 {
		return nil, fmt.Errorf("command count exceeds uint8 range")
	}

	responses := commands[0].Responses
	commandType := commands[0].Type

	responseType, err := resolveCommandType(codec, commandType)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	buffer.WriteByte(byte(len(commands)))
	buffer.WriteByte(responseType)

	for _, response := range responses {
		commandBytes, err := encodeHexMessage(response.Response, response.HexMessage)
		if err != nil {
			return nil, err
		}

		var payload []byte
		switch codec {
		case Codec12:
			payload = commandBytes
		case Codec13:
			if response.Timestamp == nil {
				return nil, fmt.Errorf("codec 13 requires timestamp")
			}
			ts, err := tools.EncodeTimestampMillis(response.Timestamp)
			if err != nil {
				return nil, err
			}
			payload = append(ts, commandBytes...)
		case Codec14:
			imei, err := tools.EncodeIMEI(response.IMEI)
			if err != nil {
				return nil, err
			}
			payload = append(imei, commandBytes...)
		case Codec15:
			if response.Timestamp == nil {
				return nil, fmt.Errorf("codec 15 requires timestamp")
			}
			ts, err := tools.EncodeTimestampSeconds(response.Timestamp)
			if err != nil {
				return nil, err
			}
			imei, err := tools.EncodeIMEI(response.IMEI)
			if err != nil {
				return nil, err
			}
			payload = append(ts, imei...)
			payload = append(payload, commandBytes...)
		default:
			return nil, fmt.Errorf("unsupported command codec: 0x%02X", byte(codec))
		}

		sizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sizeBytes, uint32(len(payload)))
		buffer.Write(sizeBytes)
		buffer.Write(payload)
	}

	buffer.WriteByte(byte(len(commands)))
	return buffer.Bytes(), nil
}

func encodeHexMessage(response string, hexMessage string) ([]byte, error) {
	if hexMessage != "" {
		data, err := decodeHex(hexMessage)
		if err != nil {
			return nil, fmt.Errorf("invalid hex message: %w", err)
		}
		return data, nil
	}
	if response == "" {
		return nil, fmt.Errorf("response is empty")
	}
	return []byte(response), nil
}

func decodeHex(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("even-length hex required")
	}
	var out []byte
	for i := 0; i < len(s); i += 2 {
		b, err := hexByte(s[i], s[i+1])
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func hexByte(hi, lo byte) (byte, error) {
	h, ok := unhex(hi)
	if !ok {
		return 0, fmt.Errorf("invalid hex digit %q", hi)
	}
	l, ok := unhex(lo)
	if !ok {
		return 0, fmt.Errorf("invalid hex digit %q", lo)
	}
	return h<<4 | l, nil
}

func unhex(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// resolveCommandType maps the friendly command type to its codec byte value.
func resolveCommandType(codec CodecID, commandType string) (byte, error) {
	switch commandType {
	case "", "Response", "6":
		return 6, nil
	case "Command", "5":
		if codec == Codec13 {
			return 0, fmt.Errorf("codec 13 only supports response type")
		}
		return 5, nil
	default:
		return 0, fmt.Errorf("unknown command type: %s", commandType)
	}
}

// kind returns the integer value of kind for error messages.
func kind(k Kind) int {
	return int(k)
}
