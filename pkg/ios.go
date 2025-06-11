package teltonika_go

import (
	"encoding/hex"
	"strconv"

	io_domain "github.com/danieljvsa/teltonika-go/internal/io"
	tools "github.com/danieljvsa/teltonika-go/tools"
)

func DecodeIos8(data []byte, startByte int64) (*io_domain.ResponseDecode, error) {
	ios_data := []io_domain.IOData{}
	if len(data) < 4 {
		return &io_domain.ResponseDecode{IOs: ios_data, NumberOfIOs: 0, LastByte: 0}, nil
	}

	ios_read := 0
	byte := startByte
	ios_number, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1

	number_ios_one_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1
	for range int(number_ios_one_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+1])

		byte += 1
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_two_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1
	for range int(number_ios_two_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+2])

		byte += 2
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_four_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1
	for range int(number_ios_four_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+4])

		byte += 4
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_eight_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1

	for range int(number_ios_eight_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 1
		value := hex.EncodeToString(data[byte : byte+8])

		byte += 8
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	return &io_domain.ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number, LastByte: byte}, nil
}

func DecodeIos8Extended(data []byte, startByte int64) (*io_domain.ResponseDecode, error) {
	ios_data := []io_domain.IOData{}
	if len(data) < 4 {
		return &io_domain.ResponseDecode{IOs: ios_data, NumberOfIOs: 0}, nil
	}

	ios_read := 0
	byte := startByte
	ios_number, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 2
	number_ios_one_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 2
	for range int(number_ios_one_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 2
		value := hex.EncodeToString(data[byte : byte+1])

		byte += 1
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_two_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 2
	for range int(number_ios_two_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 2
		value := hex.EncodeToString(data[byte : byte+2])

		byte += 2
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_four_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 2
	for range int(number_ios_four_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 2
		value := hex.EncodeToString(data[byte : byte+4])

		byte += 4
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_eight_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 2
	for range int(number_ios_eight_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 2
		value := hex.EncodeToString(data[byte : byte+8])

		byte += 8
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_x_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 2
	for range int(number_ios_x_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 2

		io_length, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+2]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += 2

		value := hex.EncodeToString(data[byte : byte+io_length])

		byte += io_length
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	return &io_domain.ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number, LastByte: byte}, nil
}

func DecodeIos16(data []byte, startByte int64) (*io_domain.ResponseDecode, error) {
	ios_read := 0
	ios_id_length := int64(2)
	byte := startByte
	ios_data := []io_domain.IOData{}

	if len(data) < 4 {
		return &io_domain.ResponseDecode{IOs: ios_data, NumberOfIOs: 0, LastByte: 0}, nil
	}

	generation_type, err := tools.GetGenerationType(data[byte:byte+1], 0, 1)
	if err != nil {
		return nil, err
	}
	byte += 1

	ios_number, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1

	number_ios_one_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1
	for range int(number_ios_one_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+ios_id_length]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += ios_id_length
		value := hex.EncodeToString(data[byte : byte+1])

		byte += 1
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_two_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1
	for range int(number_ios_two_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+ios_id_length]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += ios_id_length
		value := hex.EncodeToString(data[byte : byte+2])

		byte += 2
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_four_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1
	for range int(number_ios_four_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+ios_id_length]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += ios_id_length
		value := hex.EncodeToString(data[byte : byte+4])

		byte += 4
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	number_ios_eight_byte, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+1]), 16, 64)
	if err != nil {
		return nil, err
	}
	byte += 1

	for range int(number_ios_eight_byte) {
		id, err := strconv.ParseInt(hex.EncodeToString(data[byte:byte+ios_id_length]), 16, 64)
		if err != nil {
			return nil, err
		}
		byte += ios_id_length
		value := hex.EncodeToString(data[byte : byte+8])

		byte += 8
		io := io_domain.IOData{IO: id, Value: value}
		ios_data = append(ios_data, io)
		ios_read += 1
	}

	return &io_domain.ResponseDecode{IOs: ios_data, NumberOfIOs: ios_number, LastByte: byte, GenerationType: generation_type}, nil
}
