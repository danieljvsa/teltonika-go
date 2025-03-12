package main

import (
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
	case "8E":
	case "0C":
	case "0D":
	case "0E":
	case "0F":
	case "10":
	default:
	}
}
