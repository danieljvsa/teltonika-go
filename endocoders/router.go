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
	case "8E":
	case "0C":
	case "0D":
	case "0E":
	case "0F":
	case "10":
	default:
	}
}
