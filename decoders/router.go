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
    codec8.decode()
	case "8E":
	case "0C":
	case "0D":
	case "0E":
	case "0F":
	case "10":
	default:
	}
}
