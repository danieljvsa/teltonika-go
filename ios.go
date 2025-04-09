package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type IOData struct {
	IO    int64
	Value string
}

type ResponseDecode struct {
	IOs         []IOData
	NumberOfIOs int64
}

func decodeIos(data []byte, dataLength int64, startByte int64) (*ResponseDecode, error) {
	ios_data := []IOData{}
	if int(dataLength) < 4 {
		return &ResponseDecode{IOs: ios_data, NumberOfIOs: 0}, nil
	}
	fmt.Println("Data Length:", dataLength, len(data))
	if int(dataLength) != len(data) {
		return nil, fmt.Errorf("Invalid data length %d", len(data))
	}
	ios_read := 0
	byte := startByte
	ios_number, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing number of ios:", err)
		return nil, err
	}
	fmt.Println("Number of IO's: ", ios_number)
	byte += 1

	number_ios_one_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing id of ios one-byte:", err)
		return nil, err
	}
	byte += 1
	for i := 0; i < int(number_ios_one_byte); i++ {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			fmt.Println("Error parsing value of ios one-byte:", err)
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+1])

		byte += 1
		io := IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	if ios_read >= int(ios_number) {
		return &ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number}, nil
	}

	number_ios_two_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing id of ios two-byte:", err)
		return nil, err
	}
	byte += 1
	for i := 0; i < int(number_ios_two_byte); i++ {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			fmt.Println("Error parsing value of ios two-byte:", err)
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+2])

		byte += 2
		io := IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	if ios_read >= int(ios_number) {
		return &ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number}, nil
	}

	number_ios_four_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing id of ios two-byte:", err)
		return nil, err
	}
	byte += 1
	for i := 0; i < int(number_ios_four_byte); i++ {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			fmt.Println("Error parsing value of ios four-byte:", err)
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+4])

		byte += 4
		io := IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	if ios_read >= int(ios_number) {
		return &ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number}, nil
	}

	number_ios_eight_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing id of ios two-byte:", err)
		return nil, err
	}
	byte += 1
	for i := 0; i < int(number_ios_eight_byte); i++ {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			fmt.Println("Error parsing value of ios four-byte:", err)
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+8])

		byte += 8
		io := IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	return &ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number}, nil
}
