package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
)

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
