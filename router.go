package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func RouterDecoder(request []byte) {
	header := hex.EncodeToString(request[:4])
	if header != "00000000" {
		fmt.Println("Invalid header:", header)
		return
	}

	dataLength, err := strconv.ParseInt(hex.EncodeToString(request[4:8]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing data length:", err)
		return
	}

	codec := hex.EncodeToString(request[8:9])
	data := request[9:]
	fmt.Println("Codec:", codec)
	fmt.Println("Data Length:", dataLength)
	switch string(codec) {
	case "08":
		decodeCodec8(data, dataLength)
	case "8E":
		decodeCodec8Ext(data, dataLength)
	case "0C":
		decodeCodec12(data, dataLength)
	case "0D":
		decodeCodec13(data, dataLength)
	case "0E":
		decodeCodec14(data, dataLength)
	case "0F":
		decodeCodec15(data, dataLength)
	case "10":
		decodeCodec16(data, dataLength)
	default:
		fmt.Printf("Unknown codec: %s", codec)
	}
}

func RouterEncoder(request []byte) {
	header := hex.EncodeToString(request[:4])
	if string(header) != "00000000" {
		fmt.Println("Invalid header")
		return
	}

	dataLength, err := strconv.ParseInt(hex.EncodeToString(request[4:8]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing data length:", err)
		return
	}

	codec := hex.EncodeToString(request[8:9])
	data := request[9:]
	fmt.Println("Codec:", codec)
	fmt.Println("Data Length:", dataLength)
	fmt.Println("Data:", data)
	switch string(codec) {
	case "08":
		encodeCodec8(data, dataLength)
	case "8E":
		encodeCodec8(data, dataLength)
	case "0C":
		encodeCodec8(data, dataLength)
	case "0D":
		encodeCodec8(data, dataLength)
	case "0E":
		encodeCodec8(data, dataLength)
	case "0F":
		encodeCodec8(data, dataLength)
	case "10":
		encodeCodec8(data, dataLength)
	default:
		fmt.Printf("Unknown codec: %s", codec)
	}
}
