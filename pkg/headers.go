package teltonika_go

import (
	"encoding/hex"
	"fmt"
	"strconv"

	header_domain "github.com/danieljvsa/teltonika-go/internal/header"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func DecodeHeaderTCP(header []byte) (*header_domain.HeaderDataTCP, error) {
	read := 0
	zeroBytes := hex.EncodeToString(header[:read+4])
	if zeroBytes != "00000000" {
		return nil, fmt.Errorf("header is not valid")
	}
	read += 4

	dataLength, err := strconv.ParseInt(hex.EncodeToString(header[read:read+4]), 16, 64)
	if err != nil {
		return nil, err
	}
	read += 4

	headerData := &header_domain.HeaderDataTCP{
		Header:     "00000000",
		DataLength: dataLength,
		LastByte:   read,
	}
	return headerData, nil
}

func DecodeHeaderUDP(header []byte) (*header_domain.HeaderDataUDP, error) {
	read := 0
	if len(header) < 8 {
		return nil, fmt.Errorf("header is too small")
	}

	length, err := strconv.ParseInt(hex.EncodeToString(header[:read+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	read += 2

	packetID, err := strconv.ParseInt(hex.EncodeToString(header[read:read+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	read += 3

	avlPacketID, err := strconv.ParseInt(hex.EncodeToString(header[read:read+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	read += 1

	imeiLength, err := strconv.ParseInt(hex.EncodeToString(header[read:read+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	read += 2

	hexImei := hex.EncodeToString(header[read : read+int(imeiLength)])
	bytesImei, _ := hex.DecodeString(hexImei)
	imei := string(bytesImei)

	read += int(imeiLength)

	data := &header_domain.HeaderDataUDP{
		Length:      length,
		PacketID:    packetID,
		AVLPacketID: avlPacketID,
		IMEILength:  imeiLength,
		IMEI:        imei,
		LastByte:    read,
	}

	return data, nil
}

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

		headerData := &header_domain.HeaderData{
			HeaderTCP: nil,
			HeaderUDP: decodedUDPHeader,
			Protocol:  protocolData.Protocol,
			LastByte:  decodedUDPHeader.LastByte,
		}

		return headerData, nil
	} else if protocolData.Protocol == "TCP" {
		decodedTCPHeader, err := DecodeHeaderTCP(header)
		if err != nil {
			return nil, err
		}

		headerData := &header_domain.HeaderData{
			HeaderTCP: decodedTCPHeader,
			HeaderUDP: nil,
			Protocol:  protocolData.Protocol,
			LastByte:  decodedTCPHeader.LastByte,
		}

		return headerData, nil
	}

	headerData := &header_domain.HeaderData{
		HeaderTCP: nil,
		HeaderUDP: nil,
		Protocol:  "Unknown",
		LastByte:  0,
	}

	return headerData, nil
}
