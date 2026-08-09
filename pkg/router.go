package teltonika_go

import (
	"fmt"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	header_domain "github.com/danieljvsa/teltonika-go/internal/header"
	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// LoginDecoder parses a login packet and returns a CodecDecoded wrapper.
func LoginDecoder(request []byte) *decoder_domain.CodecDecoded {
	packet, err := teltonika.DecodeLogin(request)
	if err != nil {
		return &decoder_domain.CodecDecoded{Response: nil, Error: err}
	}

	length := int64(len(packet.IMEI))
	res := &decoder_domain.ResponseType{
		Type: "Login",
		Result: decoder_domain.CodecHeaderResponse{
			Length: &length,
			IMEI:   &packet.IMEI,
		},
	}
	return &decoder_domain.CodecDecoded{Response: res, Error: nil}
}

// TramDecoder parses a complete AVL data packet carrying an AVL or command
// codec and returns a CodecDecoded wrapper.
func TramDecoder(request []byte) *decoder_domain.CodecDecoded {
	packet, err := teltonika.Decode(request)
	if err != nil {
		return &decoder_domain.CodecDecoded{Response: nil, Error: err}
	}
	if packet.Kind == teltonika.KindLogin {
		return &decoder_domain.CodecDecoded{Response: nil, Error: fmt.Errorf("login packet is not valid for tram decoder")}
	}

	response := &decoder_domain.ResponseType{Result: decoder_domain.CodecHeaderResponse{}, Type: "Tram"}
	headerData := headerFromPacket(packet)
	codecDataIntern := codecDataFromPacket(packet)
	response.Result = decoder_domain.CodecHeaderResponse{CodecData: codecDataIntern, HeaderData: headerData}
	return &decoder_domain.CodecDecoded{Response: response, Error: nil}
}

// TramEncoder is a legacy placeholder; it returns an empty response.
func TramEncoder(request []byte) *decoder_domain.CodecDecoded {
	return &decoder_domain.CodecDecoded{Response: nil, Error: nil}
}

// headerFromPacket converts a public packet header to the legacy internal
// header model used by pkg.
func headerFromPacket(packet *teltonika.Packet) *header_domain.HeaderData {
	if packet.Protocol == teltonika.ProtocolUDP {
		udp := packet.Header.UDP
		var imei string
		var imeiLength int64
		if udp != nil {
			imei = udp.IMEI
			imeiLength = udp.IMEILength
			if imeiLength == 0 && len(imei) > 0 {
				imeiLength = int64(len(imei))
			}
		}
		return &header_domain.HeaderData{
			HeaderTCP: nil,
			HeaderUDP: &header_domain.HeaderDataUDP{
				Length:      lengthOrDefault(udp),
				PacketID:    packetIDOrDefault(udp),
				AVLPacketID: avlPacketIDOrDefault(udp),
				IMEILength:  imeiLength,
				IMEI:        imei,
				LastByte:    headerUDPLastByte(packet),
			},
			Protocol: "UDP",
			LastByte: headerUDPLastByte(packet),
		}
	}
	return &header_domain.HeaderData{
		HeaderTCP: &header_domain.HeaderDataTCP{
			Header:     "00000000",
			DataLength: tcpDataLengthOrDefault(packet),
			LastByte:   8,
		},
		HeaderUDP: nil,
		Protocol:  "TCP",
		LastByte:  8,
	}
}

func lengthOrDefault(udp *teltonika.HeaderUDP) int64 {
	if udp == nil {
		return 0
	}
	return udp.Length
}

func packetIDOrDefault(udp *teltonika.HeaderUDP) int64 {
	if udp == nil {
		return 0
	}
	return udp.PacketID
}

func avlPacketIDOrDefault(udp *teltonika.HeaderUDP) int64 {
	if udp == nil {
		return 0
	}
	return udp.AVLPacketID
}

func headerUDPLastByte(packet *teltonika.Packet) int {
	udp := packet.Header.UDP
	imeiLen := 0
	if udp != nil {
		imeiLen = len(udp.IMEI)
	}
	return 8 + imeiLen
}

func tcpDataLengthOrDefault(packet *teltonika.Packet) int64 {
	tcp := packet.Header.TCP
	if tcp == nil {
		return 0
	}
	return tcp.DataLength
}

// codecDataFromPacket converts a public Packet to the legacy internal
// decoder CodecData model.
func codecDataFromPacket(packet *teltonika.Packet) *decoder_domain.CodecData {
	if packet.Kind == teltonika.KindLogin {
		return &decoder_domain.CodecData{NumberOfRecords: 0, Records: []decoder_domain.Record{}}
	}

	if len(packet.Commands) > 0 {
		cmd := packet.Commands[0]
		var responses []tool_domain.CommandResponse
		for _, r := range cmd.Responses {
			responses = append(responses, tool_domain.CommandResponse{
				Timestamp:   r.Timestamp,
				Response:    r.Response,
				HexMessage:  r.HexMessage,
				CommandType: r.CommandType,
				IMEI:        r.IMEI,
			})
		}
		return &decoder_domain.CodecData{
			NumberOfRecords: int64(len(packet.Commands)),
			Records: []decoder_domain.Record{
				{
					CommandType:      &cmd.Type,
					CommandResponses: &responses,
				},
			},
		}
	}

	var records []decoder_domain.Record
	for _, avl := range packet.Records {
		timestamp := avl.Timestamp
		priority := avl.Priority
		eventIO := avl.EventIO
		numberOfIOs := int64(len(avl.IOElements))
		gps := tool_domain.GPSData{
			Latitude:  avl.GPS.Latitude,
			Longitude: avl.GPS.Longitude,
			Altitude:  avl.GPS.Altitude,
			Angle:     avl.GPS.Angle,
			Satelites: avl.GPS.Satellites,
			Speed:     avl.GPS.Speed,
		}
		ios := make([]io_domain.IOData, 0, len(avl.IOElements))
		for _, io := range avl.IOElements {
			ios = append(ios, io_domain.IOData{IO: io.ID, Value: io.Value})
		}
		records = append(records, decoder_domain.Record{
			Timestamp:   &timestamp,
			Priority:    &priority,
			GPSData:     &gps,
			EventIO:     &eventIO,
			NumberOfIOs: &numberOfIOs,
			IOs:         &ios,
		})
	}
	return &decoder_domain.CodecData{
		NumberOfRecords: int64(len(records)),
		Records:         records,
	}
}
