package teltonika_go_test

import (
	"encoding/binary"
	"testing"
	"time"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
	pkg "github.com/danieljvsa/teltonika-go/pkg"
)

func buildTCPTram(codecID byte, payload []byte) []byte {
	dataLength := uint32(len(payload) + 1)
	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[4:8], dataLength)
	tram := append(header, codecID)
	tram = append(tram, payload...)
	return tram
}

func TestTramEncodeDecodeRoundTripCodec8(t *testing.T) {
	timestamp := time.Unix(1701000000, 0).UTC()
	priority := int64(1)
	eventIO := int64(5)
	gps := &tool_domain.GPSData{
		Latitude:  40.4168,
		Longitude: -3.7038,
		Altitude:  667,
		Angle:     180,
		Satelites: 10,
		Speed:     90,
	}
	ios := []io_domain.IOData{
		{IO: 1, Value: "01"},
		{IO: 2, Value: "0001"},
	}

	record := decoder_domain.Record{
		Timestamp: &timestamp,
		Priority:  &priority,
		GPSData:   gps,
		EventIO:   &eventIO,
		IOs:       &ios,
	}
	codecData := &decoder_domain.CodecData{
		NumberOfRecords: 1,
		Records:         []decoder_domain.Record{record},
	}

	payload, err := pkg.EncodeCodec8(codecData)
	if err != nil {
		t.Fatalf("EncodeCodec8 failed: %v", err)
	}

	tram := buildTCPTram(0x08, payload)
	decoded := pkg.TramDecoder(tram)
	if decoded.Error != nil {
		t.Fatalf("TramDecoder failed: %v", decoded.Error)
	}

	result := decoded.Response.Result.CodecData
	if result == nil || len(result.Records) != 1 {
		t.Fatalf("expected 1 decoded record")
	}
	decodedRecord := result.Records[0]
	if decodedRecord.GPSData == nil || decodedRecord.Timestamp == nil || decodedRecord.EventIO == nil {
		t.Fatalf("decoded record missing fields")
	}
	if decodedRecord.GPSData.Speed != gps.Speed {
		t.Errorf("expected speed %d, got %d", gps.Speed, decodedRecord.GPSData.Speed)
	}
	if decodedRecord.GPSData.Altitude != gps.Altitude {
		t.Errorf("expected altitude %d, got %d", gps.Altitude, decodedRecord.GPSData.Altitude)
	}
	if decodedRecord.EventIO == nil || *decodedRecord.EventIO != eventIO {
		t.Errorf("expected event IO %d, got %v", eventIO, decodedRecord.EventIO)
	}
}
