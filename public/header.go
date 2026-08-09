package teltonika

import (
	"encoding/binary"
	"fmt"
)

type decodedHeader struct {
	protocol Protocol
	header   Header
	lastByte int
}

// isTCP reports whether data begins with the fixed TCP AVL header magic.
func isTCP(data []byte) bool {
	return len(data) >= 4 && data[0] == 0x00 && data[1] == 0x00 && data[2] == 0x00 && data[3] == 0x00
}

// decodeHeader parses the transport header at the start of a data packet.
func decodeHeader(data []byte) (*decodedHeader, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("header is too small")
	}

	if isTCP(data) {
		length := int64(binary.BigEndian.Uint32(data[4:8]))
		return &decodedHeader{
			protocol: ProtocolTCP,
			header: Header{TCP: &HeaderTCP{
				DataLength: length,
			}},
			lastByte: 8,
		}, nil
	}

	// UDP header: length(2) packetId(2) [version skipped] avlId(1)
	//            imeiLen(2) imei(...)
	length := int64(binary.BigEndian.Uint16(data[0:2]))
	read := 2

	if len(data) < read+3 {
		return nil, fmt.Errorf("header is too small")
	}
	packetID := int64(binary.BigEndian.Uint16(data[read : read+2]))
	read += 3 // 2 for the packet id plus one version byte

	if len(data) < read+3 {
		return nil, fmt.Errorf("header is too small")
	}
	avlPacketID := int64(data[read])
	read += 1
	imeiLength := int64(binary.BigEndian.Uint16(data[read : read+2]))
	read += 2

	if len(data) < read+int(imeiLength) {
		return nil, fmt.Errorf("header length exceeds available data")
	}
	imei := string(data[read : read+int(imeiLength)])
	read += int(imeiLength)

	return &decodedHeader{
		protocol: ProtocolUDP,
		header: Header{UDP: &HeaderUDP{
			Length:      length,
			PacketID:    packetID,
			AVLPacketID: avlPacketID,
			IMEILength:  imeiLength,
			IMEI:        imei,
		}},
		lastByte: read,
	}, nil
}
