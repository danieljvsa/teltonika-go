package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"

	tool_domain "github.com/danieljvsa/teltonika-go/internal/tool"
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

func Login(tram []byte) (*tool_domain.LoginData, error) {
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

	return &tool_domain.LoginData{Length: length, IMEI: imei}, nil
}
