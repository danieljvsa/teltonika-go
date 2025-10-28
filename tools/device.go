package tools

import "encoding/hex"

func DecodeIMEI(data []byte) (string, error) {
	hexImei := hex.EncodeToString(data)
	bytesImei, err := hex.DecodeString(hexImei)
	if err != nil {
		return "", err
	}
	imei := string(bytesImei)
	return imei, nil
}
