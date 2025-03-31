package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func decodeCodec8(data []byte, dataLength int64) {
	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing number of records:", err)
		return
	}

	timestamp := CalcTimestamp(data[2:9])
	if timestamp == nil {
		fmt.Println("Error parsing Timestamp:", timestamp)
		return
	}

	priority, err := strconv.ParseInt(hex.EncodeToString(data[9:10]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing Priority:", err)
		return
	}

	gps_data := DecodeGPSData(data[10:25])
	if gps_data == nil {
		fmt.Println("Error parsing GPS Data:", err)
		return
	}

	event_io, err := strconv.ParseInt(hex.EncodeToString(data[27:28]), 16, 64)
		if err != nil {
			fmt.Println("Error parsing number of records:", err)
			return
		}
	
	ios_number, err := strconv.ParseInt(hex.EncodeToString(data[28:29]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing number of records:", err)
		return
	}
	
	io_data := data[27:dataLength]

	fmt.Println("Number Of Records:", numberOfRecords)
	fmt.Println("Timestamp:", timestamp)
	fmt.Println("priority:", priority)
	fmt.Println("GPS Data:", gps_data)
	fmt.Println("Event IO:", event_io)
	fmt.Println("IOs Number:", ios_number)
	fmt.Println("IO's:", io_data)
}

func decodeCodec8Ext(data []byte, dataLength int64) {}
func decodeCodec12(data []byte, dataLength int64)   {}
func decodeCodec13(data []byte, dataLength int64)   {}
func decodeCodec14(data []byte, dataLength int64)   {}
func decodeCodec15(data []byte, dataLength int64)   {}
func decodeCodec16(data []byte, dataLength int64)   {}
