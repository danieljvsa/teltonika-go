package decoder

import (
	"time"

	header "github.com/danieljvsa/teltonika-go/internal/header"
	io "github.com/danieljvsa/teltonika-go/internal/io"
	tool "github.com/danieljvsa/teltonika-go/internal/tool"
)

type Record struct {
	// Optional fields to support different codec payloads.
	// Use pointers where a value may be absent for some codecs.
	Timestamp        *time.Time
	Priority         *int64
	GPSData          *tool.GPSData
	EventIO          *int64
	NumberOfIOs      *int64
	IOs              *[]io.IOData
	CommandResponses *[]tool.CommandResponse
	CommandType      *string

	// Codec-specific metadata
	CodecID    *int //will be nil for some codecs
	RawData    *[]byte
	Attributes *map[string]any // catch-all for extra codec-specific values
}

type CodecData struct {
	NumberOfRecords int64
	Records         []Record
}

type CodecDecoded struct {
	Response *ResponseType
	Error    error
}

type ResponseType struct {
	Type   string
	Result CodecHeaderResponse
}

type CodecHeaderResponse struct {
	// Tram type fields
	CodecData  *CodecData
	HeaderData *header.HeaderData

	// Login packet fields
	Length *int64
	IMEI   *string
}
