package teltonika_go

import (
	decoder_domain "github.com/danieljvsa/teltonika-go/internal/decoder"
	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// DecodeCodec8 decodes a Codec 08 payload and returns the legacy internal
// CodecData model. data excludes the transport header and codec id byte.
func DecodeCodec8(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec8, data, protocol)
}

// DecodeCodec8Ext decodes a Codec 8E payload and returns the legacy internal
// CodecData model.
func DecodeCodec8Ext(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec8Ext, data, protocol)
}

// DecodeCodec16 decodes a Codec 16 payload and returns the legacy internal
// CodecData model.
func DecodeCodec16(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec16, data, protocol)
}

// DecodeCodec12 decodes a Codec 12 payload and returns the legacy internal
// CodecData model.
func DecodeCodec12(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec12, data, protocol)
}

// DecodeCodec13 decodes a Codec 13 payload and returns the legacy internal
// CodecData model.
func DecodeCodec13(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec13, data, protocol)
}

// DecodeCodec14 decodes a Codec 14 payload and returns the legacy internal
// CodecData model.
func DecodeCodec14(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec14, data, protocol)
}

// DecodeCodec15 decodes a Codec 15 payload and returns the legacy internal
// CodecData model.
func DecodeCodec15(data []byte, protocol string) (*decoder_domain.CodecData, error) {
	return decodeDataViaPB(teltonika.Codec15, data, protocol)
}

func decodeDataViaPB(codec teltonika.CodecID, data []byte, protocol string) (*decoder_domain.CodecData, error) {
	pb, pbErr := teltonika.DecodeCodecData(data, codec, teltonika.Protocol(protocol))
	if pbErr != nil {
		return nil, pbErr
	}
	return codecDataFromPacket(pb), nil
}
