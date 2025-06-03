package teltonicaGo

import (
	"encoding/hex"
	"fmt"
)

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
	HeaderData *HeaderData
}

func LoginDecoder(request []byte) *CodecDecoded {
	isLogin, err := isLogin(request)
	if err != nil {
		return &CodecDecoded{Response: nil, Error: err}
	}

	if isLogin {
		login, err := login(request)
		if err != nil {
			return &CodecDecoded{Response: nil, Error: err}
		}
		fmt.Println("Decoded (Login):", login)
		res := &ResponseType{Result: login, Type: "Login"}
		return &CodecDecoded{Response: res, Error: err}
	}

	return &CodecDecoded{Response: nil, Error: fmt.Errorf("login is not valid")}
}

func TramDecoder(request []byte) *CodecDecoded {
	read := 0

	headerData, err := decodeHeader(request)
	if err != nil {
		return &CodecDecoded{Response: nil, Error: err}
	}

	read += headerData.LastByte
	codec := hex.EncodeToString(request[read : read+1])

	read += 1
	data := request[read:]
	response := &ResponseType{Result: "Codec not supported", Type: "Tram"}
	switch string(codec) {
	case "08":
		res, err := decodeCodec8(data, headerData.Protocol)
		response.Result = &CodecHeaderResponse{CodecData: res, HeaderData: headerData}
		fmt.Println("Decoded (res): ", res)
		fmt.Println("Decoded (headerData): ", headerData, headerData.HeaderTCP, headerData.HeaderUDP)
		return &CodecDecoded{Response: response, Error: err}
	case "8e":
		res, err := decodeCodec8Ext(data, headerData.Protocol)
		response.Result = &CodecHeaderResponse{CodecData: res, HeaderData: headerData}
		fmt.Println("Decoded (res): ", res)
		fmt.Println("Decoded (headerData): ", headerData, headerData.HeaderTCP, headerData.HeaderUDP)
		return &CodecDecoded{Response: response, Error: err}
	case "0C":
		return &CodecDecoded{Response: response, Error: err}
	case "0D":
		return &CodecDecoded{Response: response, Error: err}
	case "0E":
		return &CodecDecoded{Response: response, Error: err}
	case "0F":
		return &CodecDecoded{Response: response, Error: err}
	case "10":
		return &CodecDecoded{Response: response, Error: err}
	default:
		return &CodecDecoded{Response: response, Error: fmt.Errorf("unknown codec: %s", codec)}
	}

}

func TramEncoder(request []byte) {}
