package tools

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type GPSData struct {
	latitude  int64
	longitude int64
	altitude  int64
	angle     int64
	satelites int64
	speed     int64
}

func CalcTimestamp(data []byte) *int64 {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return nil
	}

	return &timestamp
}

func DecodeGPSData(data []byte) *GPSData {
	latitude, err := strconv.ParseInt(hex.EncodeToString(data[:4]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
		return nil
	}
	longitude, err := strconv.ParseInt(hex.EncodeToString(data[4:8]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
		return nil
	}
	altitude, err := strconv.ParseInt(hex.EncodeToString(data[8:10]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing altitude:", err)
		return nil
	}
	angle, err := strconv.ParseInt(hex.EncodeToString(data[10:12]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing angle:", err)
		return nil
	}
	satelites, err := strconv.ParseInt(hex.EncodeToString(data[12:13]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing satelites:", err)
		return nil
	}
	speed, err := strconv.ParseInt(hex.EncodeToString(data[14:15]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing speed:", err)
		return nil
	}

	fmt.Println("latitude: ", latitude)
	fmt.Println("longitude: ", longitude)
	fmt.Println("altitude: ", altitude)
	fmt.Println("angle: ", angle)
	fmt.Println("satelites: ", satelites)
	fmt.Println("speed: ", speed)

	return &GPSData{
		latitude,
		longitude,
		altitude,
		angle,
		satelites,
		speed,
	}
}
