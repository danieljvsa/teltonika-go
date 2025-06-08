package header

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
