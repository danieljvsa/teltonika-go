package teltonika_go

import (
	"encoding/hex"
	"testing"

	"github.com/danieljvsa/teltonika-go/tools"
)

// Benchmarks for Timestamp Functions

// BenchmarkCalcTimestamp benchmarks the 8-byte millisecond timestamp conversion.
// This test measures performance of parsing hex-encoded 8-byte timestamps.
func BenchmarkCalcTimestamp(b *testing.B) {
	data, _ := hex.DecodeString("00000001795F8F00")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.CalcTimestamp(data)
	}
}

// BenchmarkCalcTimestampSeconds benchmarks 4-byte second timestamp conversion using hex.
// Tests the performance of string-based hex parsing for seconds timestamps.
func BenchmarkCalcTimestampSeconds(b *testing.B) {
	data, _ := hex.DecodeString("56D826A0")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.CalcTimestampSeconds(data)
	}
}

// BenchmarkCalcTimestampSecondsBigEndian benchmarks big-endian 4-byte timestamp conversion.
// Tests performance of binary.BigEndian parsing for seconds timestamps.
func BenchmarkCalcTimestampSecondsBigEndian(b *testing.B) {
	data, _ := hex.DecodeString("56D826A0")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.CalcTimestampSecondsBigEndian(data)
	}
}

// BenchmarkCalcTimestampSecondsLittleEndian benchmarks little-endian 4-byte timestamp conversion.
// Tests performance of binary.LittleEndian parsing for seconds timestamps.
func BenchmarkCalcTimestampSecondsLittleEndian(b *testing.B) {
	data, _ := hex.DecodeString("A026D856")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.CalcTimestampSecondsLittleEndian(data)
	}
}

// Benchmarks for CRC Functions

// BenchmarkCrc16IBM benchmarks the CRC-16 IBM checksum calculation.
// Tests the performance of bitwise CRC computation on various data sizes.
func BenchmarkCrc16IBM(b *testing.B) {
	data, _ := hex.DecodeString("0F080100000000000000000000000001795F8F00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.Crc16IBM(data)
	}
}

// BenchmarkCrc16IBMSmallData benchmarks CRC-16 on small data chunks.
// Useful for understanding overhead on minimal payloads.
func BenchmarkCrc16IBMSmallData(b *testing.B) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.Crc16IBM(data)
	}
}

// BenchmarkIsValidTram benchmarks complete tram validation including CRC check.
// Measures end-to-end performance of frame validation with CRC verification.
func BenchmarkIsValidTram(b *testing.B) {
	// Valid Codec 08 TCP frame with proper CRC
	tram, _ := hex.DecodeString("000000000000003C08010000016A571F000000000000000000000000000000000000000000000000000000000000000000000000000000000000000101095B66")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.IsValidTram(tram)
	}
}

// Benchmarks for GPS Functions

// BenchmarkDecodeGPSData benchmarks GPS data decoding from 14-byte blocks.
// Tests performance of parsing latitude, longitude, altitude, and related fields.
func BenchmarkDecodeGPSData(b *testing.B) {
	// 14 bytes of GPS data
	gpsData, _ := hex.DecodeString("0000000100000001000100020304")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.DecodeGPSData(gpsData)
	}
}

// Benchmarks for IMEI Functions

// BenchmarkDecodeIMEI benchmarks IMEI extraction and hex encoding.
// Measures performance of hex string conversion for device identification.
func BenchmarkDecodeIMEI(b *testing.B) {
	imeiData, _ := hex.DecodeString("356423401234567890ABCDEF12345678")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.DecodeIMEI(imeiData)
	}
}

// Benchmarks for Login Functions

// BenchmarkIsLogin benchmarks login frame detection.
// Tests performance of quick validation on potential login frames.
func BenchmarkIsLogin(b *testing.B) {
	loginData, _ := hex.DecodeString("000F3335363432333430313238393536")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.IsLogin(loginData)
	}
}

// BenchmarkLogin benchmarks complete login frame parsing.
// Measures full login extraction including length and IMEI decoding.
func BenchmarkLogin(b *testing.B) {
	loginData, _ := hex.DecodeString("000F3335363432333430313238393536")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.Login(loginData)
	}
}

// Composite Benchmarks for realistic scenarios

// BenchmarkFullFrameProcessing simulates processing a complete AVL frame.
// Combines CRC validation, GPS decoding, and timestamp parsing.
func BenchmarkFullFrameProcessing(b *testing.B) {
	// Complete Codec 08 frame
	frame, _ := hex.DecodeString("000000000000003C08010000016A571F000000000000000000000000000000000000000000000000000000000000000000000000000000000000000101095B66")
	gpsData, _ := hex.DecodeString("0000000100000001000100020304")
	timestamp, _ := hex.DecodeString("00000001795F8F00")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.IsValidTram(frame)
		tools.DecodeGPSData(gpsData)
		tools.CalcTimestamp(timestamp)
	}
}

// BenchmarkParallelCRC measures CRC calculation on parallel frames.
// Tests performance scaling with concurrent operations.
func BenchmarkParallelCRC(b *testing.B) {
	data, _ := hex.DecodeString("0F080100000000000000000000000001795F8F00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tools.Crc16IBM(data)
		}
	})
}
