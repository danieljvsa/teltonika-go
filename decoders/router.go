package main

import (
	"fmt"
	"strconv"
)

func router(request []byte) {
	header := request[:8]
  if string(header) != "00000000" {
    fmt.Println("Invalid header")
    return
  }
  
	dataLength, err := strconv.ParseInt(string(request[8:12]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing data length:", err)
		return
	}
  
	codec := request[12:14]
	data := request[14:]
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
		fmt.Sprintf("Unknown codec: %s", codec)
	}
}
