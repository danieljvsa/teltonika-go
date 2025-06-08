package decoder

import (
	"time"

	header "github.com/danieljvsa/teltonika-go/internal/header"
	io "github.com/danieljvsa/teltonika-go/internal/io"
	tool "github.com/danieljvsa/teltonika-go/internal/tool"
)

type Record struct {
	Timestamp   time.Time
	Priority    int64
	GPSData     tool.GPSData
	EventIO     int64
	NumberOfIOs int64
	IOs         []io.IOData
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
	Result any
}

type CodecHeaderResponse struct {
	CodecData  *CodecData
	HeaderData *header.HeaderData
}
