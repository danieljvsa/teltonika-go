package main

import (
	"codec8"
	"fmt"
	"net"
	"os"
)

func router(request []byte) {
	header := request[:8]
	dataLength := request[8:12]
	codec := request[12:14]
	data := request[14:]
	switch string(codec) {
	case "08":
		decodeCodec8(data)
	case "8E":
    decodeCodec8Ext(data)
	case "0C":
    decodeCodec12(data)
	case "0D":
    decodeCodec13(data)
	case "0E":
    decodeCodec14(data)
	case "0F":
    decodeCodec15(data)
	case "10":
    decodeCodec16(data)
	default:
    throw "Unknown codec"
	}
}
