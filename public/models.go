package teltonika

import "time"

// Protocol identifies the transport layer used by a Teltonika packet.
type Protocol string

const (
	// ProtocolTCP identifies TCP Transport Protocol AVL packets.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP identifies UDP Transport Protocol AVL packets.
	ProtocolUDP Protocol = "UDP"
)

// CodecID identifies the codec used to encode a packet payload.
type CodecID byte

const (
	// Codec8 is the base AVL data codec.
	Codec8 CodecID = 0x08
	// Codec8Ext is the extended AVL data codec (Codec 8E).
	Codec8Ext CodecID = 0x8E
	// Codec12 is a command response codec.
	Codec12 CodecID = 0x0C
	// Codec13 is a command response codec with timestamps.
	Codec13 CodecID = 0x0D
	// Codec14 is a command response codec with IMEI information.
	Codec14 CodecID = 0x0E
	// Codec15 is a command response codec with timestamps and IMEI information.
	Codec15 CodecID = 0x0F
	// Codec16 is a GPRS/IQ codec with device type and generation type.
	Codec16 CodecID = 0x10
)

// Kind describes what a Packet represents.
type Kind int

const (
	// KindData represents an AVL or command data packet.
	KindData Kind = iota
	// KindLogin represents an IMEI login packet.
	KindLogin
)

// Header infoses the transport header carried by a data packet.
type Header struct {
	// TCP holds the TCP header when Protocol is ProtocolTCP.
	TCP *HeaderTCP
	// UDP holds the UDP header when Protocol is ProtocolUDP.
	UDP *HeaderUDP
}

// HeaderTCP is the fixed 8-byte TCP AVL header.
type HeaderTCP struct {
	// DataLength is the number of bytes that follow the header
	// (codec id plus payload).
	DataLength int64
}

// HeaderUDP is the variable-length UDP AVL header.
type HeaderUDP struct {
	Length      int64
	PacketID    int64
	AVLPacketID int64
	IMEILength  int64
	IMEI        string
}

// Packet is the top-level decoded representation of a Teltonika frame.
type Packet struct {
	// Kind discriminates login packets from data packets.
	Kind Kind
	// Protocol is the transport protocol of the packet.
	Protocol Protocol
	// Header carries the transport header for data packets.
	Header Header
	// IMEI is set for KindLogin packets.
	IMEI string
	// Codec is the codec identifier for data packets.
	Codec CodecID
	// Records holds GPS/telemetry records for AVL codecs (08, 8E, 16).
	Records []AVLRecord
	// Commands holds command responses for command codecs (12, 13, 14, 15).
	Commands []Command
}

// AVLRecord is a single decoded AVL record.
type AVLRecord struct {
	Timestamp  time.Time
	Priority   int64
	GPS        GPSData
	EventIO    int64
	IOElements []IOElement
	// GenerationType is only populated for Codec 16 records.
	GenerationType string
}

// GPSData is the decoded GPS portion of an AVL record.
type GPSData struct {
	Latitude   float64
	Longitude  float64
	Altitude   int64
	Angle      int64
	Satellites int64
	Speed      int64
}

// IOElement is a single I/O element (id and hex-encoded value).
type IOElement struct {
	ID    int64
	Value string
}

// Command is a command or command response collection.
type Command struct {
	// Type is "Command" or "Response".
	Type      string
	Responses []CommandResponse
}

// CommandResponse is a single command response payload.
type CommandResponse struct {
	Timestamp   *time.Time
	Response    string
	HexMessage  string
	CommandType string
	IMEI        string
}
