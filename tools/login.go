package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
)

// IsLogin checks if a byte sequence is a valid login frame by examining
// the length field at the beginning of the tram.
//
// A login frame must have:
//   - At least 2 bytes for the length field
//   - Non-zero length value
//
// Parameters:
//   - tram: byte slice that may contain a login frame
//
// Returns:
//   - bool: true if this appears to be a login frame
//   - error: if tram is too small to be a valid frame
//
// Example:
//
//	if isLogin, err := IsLogin(frameBytes); err == nil && isLogin {
//		fmt.Println("This is a login frame")
//	}
func IsLogin(tram []byte) (bool, error) {
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

// Login decodes a login frame from Teltonika protocol and extracts
// the message length and IMEI (International Mobile Equipment Identity).
//
// Login frame format:
//   - Bytes 0-1: Length (big-endian uint16)
//   - Bytes 2+: IMEI as ASCII/hex string
//
// Parameters:
//   - tram: complete login frame bytes
//
// Returns:
//   - *decoder_domain.CodecHeaderResponse: structure with Length and IMEI fields
//   - error: if frame is invalid or parsing fails
//
// Example:
//
//	response, err := Login(loginFrameBytes)
//	if err == nil {
//		fmt.Printf("Length: %d, IMEI: %s\n", *response.Length, *response.IMEI)
//	}
func Login(tram []byte) (*decoder_domain.CodecHeaderResponse, error) {
	read := 0
	valid, err := IsLogin(tram)
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

	return &decoder_domain.CodecHeaderResponse{Length: &length, IMEI: &imei}, nil
}
