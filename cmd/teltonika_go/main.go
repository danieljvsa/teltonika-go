package teltonika_go

import (
	"encoding/hex"
	"fmt"

	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	pkg "github.com/danieljvsa/teltonika-go/pkg"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func LoginDecoder(request []byte) *decoder_domain.CodecDecoded {
	isLogin, err := tools.IsLogin(request)
	if err != nil {
		return &decoder_domain.CodecDecoded{Response: nil, Error: err}
	}

	if isLogin {
		login, err := tools.Login(request)
		if err != nil {
			return &decoder_domain.CodecDecoded{Response: nil, Error: err}
		}
		res := &decoder_domain.ResponseType{Result: login, Type: "Login"}
		return &decoder_domain.CodecDecoded{Response: res, Error: err}
	}

	return &decoder_domain.CodecDecoded{Response: nil, Error: fmt.Errorf("login is not valid")}
}

func TramDecoder(request []byte) *decoder_domain.CodecDecoded {
	read := 0

	headerData, err := pkg.DecodeHeader(request)
	if err != nil {
		return &decoder_domain.CodecDecoded{Response: nil, Error: err}
	}

	read += headerData.LastByte
	codec := hex.EncodeToString(request[read : read+1])

	read += 1
	data := request[read:]
	response := &decoder_domain.ResponseType{Result: "Codec not supported", Type: "Tram"}
	switch string(codec) {
	case "08":
		res, err := pkg.DecodeCodec8(data, headerData.Protocol)
		response.Result = &decoder_domain.CodecHeaderResponse{CodecData: res, HeaderData: headerData}
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	case "8e":
		res, err := pkg.DecodeCodec8Ext(data, headerData.Protocol)
		response.Result = &decoder_domain.CodecHeaderResponse{CodecData: res, HeaderData: headerData}
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	case "0C":
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	case "0D":
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	case "0E":
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	case "0F":
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	case "10":
		return &decoder_domain.CodecDecoded{Response: response, Error: err}
	default:
		return &decoder_domain.CodecDecoded{Response: response, Error: fmt.Errorf("unknown codec: %s", codec)}
	}

}

func TramEncoder(request []byte) *decoder_domain.CodecDecoded {
	return &decoder_domain.CodecDecoded{Response: nil, Error: nil}
}
