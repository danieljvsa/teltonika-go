package teltonika_go

import (
	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// EncodeCodec8 encodes CodecData into the legacy Codec 08 payload shape
// (record count, records, trailing count and CRC, excluding the codec id
// and transport header).
func EncodeCodec8(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec8, codecData)
}

// EncodeCodec8Ext encodes CodecData into the Codec 8E payload shape.
func EncodeCodec8Ext(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec8Ext, codecData)
}

// EncodeCodec12 encodes CodecData into the Codec 12 payload shape.
func EncodeCodec12(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec12, codecData)
}

// EncodeCodec13 encodes CodecData into the Codec 13 payload shape.
func EncodeCodec13(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec13, codecData)
}

// EncodeCodec14 encodes CodecData into the Codec 14 payload shape.
func EncodeCodec14(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec14, codecData)
}

// EncodeCodec15 encodes CodecData into the Codec 15 payload shape.
func EncodeCodec15(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec15, codecData)
}

// EncodeCodec16 encodes CodecData into the Codec 16 payload shape.
func EncodeCodec16(codecData *decoder_domain.CodecData) ([]byte, error) {
	return encodeCodecData(teltonika.Codec16, codecData)
}

func encodeCodecData(codec teltonika.CodecID, codecData *decoder_domain.CodecData) ([]byte, error) {
	packet, err := packetFromCodecData(codec, codecData)
	if err != nil {
		return nil, err
	}
	frame, err := teltonika.Encode(packet)
	if err != nil {
		return nil, err
	}
	// Strip the 8-byte TCP header and the codec id byte to restore the
	// legacy payload shape (records + trailing count + CRC).
	return frame[9:], nil
}

// packetFromCodecData converts the legacy internal CodecData model into a
// public Packet.
func packetFromCodecData(codec teltonika.CodecID, codecData *decoder_domain.CodecData) (*teltonika.Packet, error) {
	packet := &teltonika.Packet{
		Kind:     teltonika.KindData,
		Protocol: teltonika.ProtocolTCP,
		Codec:    codec,
	}

	if len(codecData.Records) == 0 {
		return nil, nil // let Encode report "no records"
	}

	for _, rec := range codecData.Records {
		if rec.CommandResponses != nil {
			var responses []teltonika.CommandResponse
			for _, cr := range *rec.CommandResponses {
				responses = append(responses, teltonika.CommandResponse{
					Timestamp:   cr.Timestamp,
					Response:    cr.Response,
					HexMessage:  cr.HexMessage,
					CommandType: cr.CommandType,
					IMEI:        cr.IMEI,
				})
			}
			cmdType := ""
			if rec.CommandType != nil {
				cmdType = *rec.CommandType
			} else if len(responses) > 0 {
				cmdType = responses[0].CommandType
			}
			packet.Commands = append(packet.Commands, teltonika.Command{
				Type:      cmdType,
				Responses: responses,
			})
			continue
		}

		gps := teltonika.GPSData{
			Latitude:   rec.GPSData.Latitude,
			Longitude:  rec.GPSData.Longitude,
			Altitude:   rec.GPSData.Altitude,
			Angle:      rec.GPSData.Angle,
			Satellites: rec.GPSData.Satelites,
			Speed:      rec.GPSData.Speed,
		}
		var elements []teltonika.IOElement
		if rec.IOs != nil {
			for _, io := range *rec.IOs {
				elements = append(elements, teltonika.IOElement{ID: io.IO, Value: io.Value})
			}
		}
		record := teltonika.AVLRecord{
			Priority:   ptrValue(rec.Priority),
			GPS:        gps,
			EventIO:    ptrValue(rec.EventIO),
			IOElements: elements,
		}
		if rec.Timestamp != nil {
			record.Timestamp = *rec.Timestamp
		}
		if codec == teltonika.Codec16 && rec.Attributes != nil {
			if v, ok := (*rec.Attributes)["generation_type"]; ok {
				record.GenerationType, _ = v.(string)
			}
		}
		packet.Records = append(packet.Records, record)
	}

	return packet, nil
}

func ptrValue(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
