package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type GPSData struct {
	Latitude  float64
	Longitude float64
	Altitude  int64
	Angle     int64
	Satelites int64
	Speed     int64
}

type LoginData struct {
	Length int64
	IMEI   string
}

type ProtocolData struct {
	Protocol string
}

type HeaderDataTCP struct {
	Header     string
	DataLength int64
	LastByte   int
}

type HeaderDataUDP struct {
	Length      int64
	PacketID    int64
	AVLPacketID int64
	IMEILength  int64
	IMEI        string
	LastByte    int
}

type HeaderData struct {
	HeaderTCP *HeaderDataTCP
	HeaderUDP *HeaderDataUDP
	Protocol  string
	LastByte  int
}

func CalcTimestamp(data []byte) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return nil, err
	}
	date := time.UnixMilli(timestamp).UTC()
	return &date, nil
}

func DecodeGPSData(data []byte) (*GPSData, error) {
	if len(data) < 14 {
		fmt.Println("Invalid data length", len(data))
		return nil, fmt.Errorf("invalid data length %d", len(data))
	}

	lat, err := strconv.ParseUint(hex.EncodeToString(data[4:8]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
		return nil, err
	}

	long, err := strconv.ParseUint(hex.EncodeToString(data[:4]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
		return nil, err
	}
	altitude, err := strconv.ParseInt(hex.EncodeToString(data[8:10]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing altitude:", err)
		return nil, err
	}
	angle, err := strconv.ParseInt(hex.EncodeToString(data[10:12]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing angle:", err)
		return nil, err
	}
	satelites, err := strconv.ParseInt(hex.EncodeToString(data[12:13]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing satelites:", err)
		return nil, err
	}

	speed, err := strconv.ParseInt(hex.EncodeToString(data[13:14]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing speed:", err)
		return nil, err
	}

	longitude := float64(int32(long)) / 10000000.0
	latitude := float64(int32(lat)) / 10000000.0

	return &GPSData{
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
		Angle:     angle,
		Satelites: satelites,
		Speed:     speed,
	}, nil
}

// CRC16 IBM/ARC implementation (poly = 0x8005, reflected = true, init = 0x0000)
func crc16IBM(data []byte) uint16 {
	var crc uint16 = 0x0000
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// Validate tram with last 2 bytes as CRC (little-endian)
func isValidTram(tram []byte) bool {
	if len(tram) < 2 {
		return false
	}
	length := len(tram)
	data := tram[:len(tram)-4]
	receivedCRC := uint16(binary.BigEndian.Uint32(tram[length-4:]))
	calculatedCRC := crc16IBM(data)
	return receivedCRC == calculatedCRC
}

func decodeHeaderTCP(header []byte) (*HeaderDataTCP, error) {
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

	headerData := &HeaderDataTCP{
		Header:     "00000000",
		DataLength: dataLength,
		LastByte:   read,
	}
	return headerData, nil
}

func decodeHeaderUDP(header []byte) (*HeaderDataUDP, error) {
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

	data := &HeaderDataUDP{
		Length:      length,
		PacketID:    packetID,
		AVLPacketID: avlPacketID,
		IMEILength:  imeiLength,
		IMEI:        imei,
		LastByte:    read,
	}

	return data, nil
}

func isLogin(tram []byte) (bool, error) {
	if len(tram) < 2 {
		return false, fmt.Errorf("tram is too small")
	}

	length, err := strconv.ParseInt(hex.EncodeToString(tram[:2]), 16, 64)
	if err != nil {
		return false, err
	}

	if length == 0 {
		return false, nil
	}

	return true, nil
}

func login(tram []byte) (*LoginData, error) {
	read := 0
	valid, err := isLogin(tram)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("login message is not valid")
	}

	length, err := strconv.ParseInt(hex.EncodeToString(tram[read:read+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	read += 2

	hexImei := hex.EncodeToString(tram[read:])
	bytesImei, _ := hex.DecodeString(hexImei)
	imei := string(bytesImei)

	return &LoginData{Length: length, IMEI: imei}, nil
}

func getProtocol(header []byte) (*ProtocolData, error) {
	if len(header) < 4 {
		return nil, fmt.Errorf("header is too small")
	}

	if hex.EncodeToString(header[:4]) == "00000000" {
		return &ProtocolData{Protocol: "TCP"}, nil
	}

	return &ProtocolData{Protocol: "UDP"}, nil
}

func decodeHeader(header []byte) (*HeaderData, error) {
	protocolData, err := getProtocol(header)
	if err != nil {
		return nil, err
	}

	if protocolData.Protocol == "UDP" {
		decodedUDPHeader, err := decodeHeaderUDP(header)
		if err != nil {
			return nil, err
		}

		headerData := &HeaderData{
			HeaderTCP: nil,
			HeaderUDP: decodedUDPHeader,
			Protocol:  protocolData.Protocol,
			LastByte:  decodedUDPHeader.LastByte,
		}

		return headerData, nil
	} else if protocolData.Protocol == "TCP" {
		decodedTCPHeader, err := decodeHeaderTCP(header)
		if err != nil {
			return nil, err
		}

		headerData := &HeaderData{
			HeaderTCP: decodedTCPHeader,
			HeaderUDP: nil,
			Protocol:  protocolData.Protocol,
			LastByte:  decodedTCPHeader.LastByte,
		}

		return headerData, nil
	}

	headerData := &HeaderData{
		HeaderTCP: nil,
		HeaderUDP: nil,
		Protocol:  "Unknown",
		LastByte:  0,
	}

	return headerData, nil
}
