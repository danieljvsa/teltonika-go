package decoders

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/danieljvsa/geolocation_server/tools"
)

func decodeCodec8(data []byte, dataLength int64) {
	numberOfRecords, err := strconv.ParseInt(hex.EncodeToString(data[:1]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing number of records:", err)
		return
	}

	timestamp := tools.CalcTimestamp(data[1:9])
	if timestamp == nil {
		fmt.Println("Error parsing number of records:", err)
		return
	}

	priority, err := strconv.ParseInt(hex.EncodeToString(data[9:10]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing number of records:", err)
		return
	}

	gps_data := tools.DecodeGPSData(data[10:25])
	if gps_data == nil {
		fmt.Println("Error parsing number of records:", err)
		return
	}

	io_data := data[27:dataLength]

	fmt.Println("Number Of Records:", numberOfRecords)
	fmt.Println("Timestamp:", timestamp)
	fmt.Println("priority:", priority)
	fmt.Println("GPS Data:", gps_data)
	fmt.Println("IO's:", io_data)
}
