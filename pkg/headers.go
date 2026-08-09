package teltonika_go

import (
	"encoding/binary"
	"fmt"

	header_domain "github.com/danieljvsa/teltonika-go/internal/header"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

// DecodeHeaderTCP decodes an 8-byte TCP AVL header into the legacy internal
// header model.
func DecodeHeaderTCP(header []byte) (*header_domain.HeaderDataTCP, error) {
	if len(header) < 8 {
		return nil, fmt.Errorf("header is too small")
	}
	if !isTCPMagic(header) {
		return nil, fmt.Errorf("header is not valid")
	}
	dataLength := int64(binary.BigEndian.Uint32(header[4:8]))
	return &header_domain.HeaderDataTCP{
		Header:     "00000000",
		DataLength: dataLength,
		LastByte:   8,
	}, nil
}

// DecodeHeaderUDP decodes a UDP AVL header into the legacy internal model.
func DecodeHeaderUDP(header []byte) (*header_domain.HeaderDataUDP, error) {
	if len(header) < 8 {
		return nil, fmt.Errorf("header is too small")
	}
	if isTCPMagic(header) {
		return nil, fmt.Errorf("header is not a UDP header")
	}
	read := 2
	length := int64(binary.BigEndian.Uint16(header[0:2]))
	if len(header) < read+3 {
		return nil, fmt.Errorf("header is too small")
	}
	packetID := int64(binary.BigEndian.Uint16(header[read : read+2]))
	read += 3 // 2 for packet id + 1 version byte
	if len(header) < read+3 {
		return nil, fmt.Errorf("header is too small")
	}
	avlPacketID := int64(header[read])
	read += 1
	imeiLength := int64(binary.BigEndian.Uint16(header[read : read+2]))
	read += 2
	if len(header) < read+int(imeiLength) {
		return nil, fmt.Errorf("header length exceeds available data")
	}
	imei := string(header[read : read+int(imeiLength)])
	read += int(imeiLength)

	return &header_domain.HeaderDataUDP{
		Length:      length,
		PacketID:    packetID,
		AVLPacketID: avlPacketID,
		IMEILength:  imeiLength,
		IMEI:        imei,
		LastByte:    read,
	}, nil
}

// DecodeHeader detects the transport protocol and decodes the header into the
// legacy internal header model.
func DecodeHeader(header []byte) (*header_domain.HeaderData, error) {
	protocolData, err := tools.GetProtocol(header)
	if err != nil {
		return nil, err
	}

	if protocolData.Protocol == "UDP" {
		decodedUDPHeader, err := DecodeHeaderUDP(header)
		if err != nil {
			return nil, err
		}
		return &header_domain.HeaderData{
			HeaderTCP: nil,
			HeaderUDP: decodedUDPHeader,
			Protocol:  "UDP",
			LastByte:  decodedUDPHeader.LastByte,
		}, nil
	}

	decodedTCPHeader, err := DecodeHeaderTCP(header)
	if err != nil {
		return nil, err
	}
	return &header_domain.HeaderData{
		HeaderTCP: decodedTCPHeader,
		HeaderUDP: nil,
		Protocol:  "TCP",
		LastByte:  decodedTCPHeader.LastByte,
	}, nil
}

// isTCPMagic reports whether header begins with the TCP zero magic.
func isTCPMagic(header []byte) bool {
	return len(header) >= 4 && header[0] == 0x00 && header[1] == 0x00 && header[2] == 0x00 && header[3] == 0x00
}
