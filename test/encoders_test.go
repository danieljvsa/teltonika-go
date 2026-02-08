package teltonika_go_test

import (
	"encoding/hex"
	"testing"
	"time"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
	pkg "github.com/danieljvsa/teltonika-go/pkg"
)

func TestEncodeDecodeCodec8RoundTrip(t *testing.T) {
	timestamp := time.Unix(1700000000, 0).UTC()
	priority := int64(1)
	eventIO := int64(5)
	gps := &tool_domain.GPSData{
		Latitude:  52.520008,
		Longitude: 13.404954,
		Altitude:  120,
		Angle:     25,
		Satelites: 7,
		Speed:     60,
	}
	ios := []io_domain.IOData{
		{IO: 1, Value: "01"},
		{IO: 2, Value: "0001"},
		{IO: 3, Value: "00000002"},
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

	encoded, err := pkg.EncodeCodec8(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec8(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.NumberOfRecords != 1 {
		t.Fatalf("expected 1 record, got %d", decoded.NumberOfRecords)
	}
	if decoded.Records[0].GPSData == nil || decoded.Records[0].EventIO == nil {
		t.Fatalf("decoded record missing required fields")
	}
	if decoded.Records[0].GPSData.Speed != gps.Speed {
		t.Errorf("expected speed %d, got %d", gps.Speed, decoded.Records[0].GPSData.Speed)
	}
}

func TestEncodeDecodeCodec8ExtRoundTrip(t *testing.T) {
	timestamp := time.Unix(1700001000, 0).UTC()
	priority := int64(1)
	eventIO := int64(258)
	gps := &tool_domain.GPSData{
		Latitude:  -33.8688,
		Longitude: 151.2093,
		Altitude:  15,
		Angle:     120,
		Satelites: 9,
		Speed:     80,
	}
	ios := []io_domain.IOData{
		{IO: 100, Value: "01"},
		{IO: 101, Value: "DEADBEEF"},
		{IO: 102, Value: "ABCDEF"},
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

	encoded, err := pkg.EncodeCodec8Ext(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec8Ext(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.NumberOfRecords != 1 {
		t.Fatalf("expected 1 record, got %d", decoded.NumberOfRecords)
	}
}

func TestEncodeDecodeCodec16RoundTrip(t *testing.T) {
	timestamp := time.Unix(1700002000, 0).UTC()
	priority := int64(2)
	eventIO := int64(512)
	gps := &tool_domain.GPSData{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  10,
		Angle:     45,
		Satelites: 6,
		Speed:     50,
	}
	ios := []io_domain.IOData{{IO: 300, Value: "05"}}
	attributes := map[string]any{"generation_type": "On Exit"}

	record := decoder_domain.Record{
		Timestamp:  &timestamp,
		Priority:   &priority,
		GPSData:    gps,
		EventIO:    &eventIO,
		IOs:        &ios,
		Attributes: &attributes,
	}
	codecData := &decoder_domain.CodecData{
		NumberOfRecords: 1,
		Records:         []decoder_domain.Record{record},
	}

	encoded, err := pkg.EncodeCodec16(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec16(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.NumberOfRecords != 1 {
		t.Fatalf("expected 1 record, got %d", decoded.NumberOfRecords)
	}
}

func TestEncodeDecodeCodec12RoundTrip(t *testing.T) {
	commandResponses := []tool_domain.CommandResponse{
		{Response: "OK"},
	}
	commandType := "Response"
	record := decoder_domain.Record{
		CommandType:      &commandType,
		CommandResponses: &commandResponses,
	}
	codecData := &decoder_domain.CodecData{
		NumberOfRecords: 1,
		Records:         []decoder_domain.Record{record},
	}

	encoded, err := pkg.EncodeCodec12(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec12(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.Records[0].CommandResponses == nil {
		t.Fatalf("decoded responses are nil")
	}
	if (*decoded.Records[0].CommandResponses)[0].Response != "OK" {
		t.Errorf("expected response OK, got %s", (*decoded.Records[0].CommandResponses)[0].Response)
	}
}

func TestEncodeDecodeCodec13RoundTrip(t *testing.T) {
	ts := time.Unix(1700003000, 0).UTC()
	commandResponses := []tool_domain.CommandResponse{
		{Response: "status", Timestamp: &ts},
	}
	commandType := "Response"
	record := decoder_domain.Record{
		CommandType:      &commandType,
		CommandResponses: &commandResponses,
	}
	codecData := &decoder_domain.CodecData{
		NumberOfRecords: 1,
		Records:         []decoder_domain.Record{record},
	}

	encoded, err := pkg.EncodeCodec13(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec13(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.Records[0].CommandResponses == nil {
		t.Fatalf("decoded responses are nil")
	}
	if (*decoded.Records[0].CommandResponses)[0].Timestamp == nil {
		t.Fatalf("expected timestamp")
	}
}

func TestEncodeDecodeCodec14RoundTrip(t *testing.T) {
	commandResponses := []tool_domain.CommandResponse{
		{Response: "OK", IMEI: "0123456789ABCDEF"},
	}
	commandType := "Response"
	record := decoder_domain.Record{
		CommandType:      &commandType,
		CommandResponses: &commandResponses,
	}
	codecData := &decoder_domain.CodecData{
		NumberOfRecords: 1,
		Records:         []decoder_domain.Record{record},
	}

	encoded, err := pkg.EncodeCodec14(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec14(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.Records[0].CommandResponses == nil {
		t.Fatalf("decoded responses are nil")
	}

	payload := append([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, []byte("OK")...)
	expectedHex := hex.EncodeToString(payload)
	if (*decoded.Records[0].CommandResponses)[0].HexMessage != expectedHex {
		t.Errorf("expected hex %s, got %s", expectedHex, (*decoded.Records[0].CommandResponses)[0].HexMessage)
	}
}

func TestEncodeDecodeCodec15RoundTrip(t *testing.T) {
	ts := time.Unix(1700004000, 0).UTC()
	commandResponses := []tool_domain.CommandResponse{
		{Response: "PING", IMEI: "0123456789ABCDEF", Timestamp: &ts},
	}
	commandType := "Response"
	record := decoder_domain.Record{
		CommandType:      &commandType,
		CommandResponses: &commandResponses,
	}
	codecData := &decoder_domain.CodecData{
		NumberOfRecords: 1,
		Records:         []decoder_domain.Record{record},
	}

	encoded, err := pkg.EncodeCodec15(codecData)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := pkg.DecodeCodec15(encoded, "TCP")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.Records[0].CommandResponses == nil {
		t.Fatalf("decoded responses are nil")
	}
	if (*decoded.Records[0].CommandResponses)[0].Timestamp == nil {
		t.Fatalf("expected timestamp")
	}
}
