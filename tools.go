package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type GPSData struct {
	Latitude  int64
	Longitude int64
	Altitude  int64
	Angle     int64
	Satelites int64
	Speed     int64
}

func CalcTimestamp(data []byte) *time.Time {
	timestamp, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return nil
	}
	date := time.UnixMilli(timestamp).UTC()
	return &date
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

	longitude = longitude / 10000000.0
	latitude = latitude / 10000000.0

	fmt.Println("Latitude:", latitude)
	fmt.Println("Longitude:", longitude)
	fmt.Println("Altitude:", altitude)
	fmt.Println("Angle:", angle)
	fmt.Println("Satelites:", satelites)
	fmt.Println("Speed:", speed)

	return &GPSData{
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
		Angle:     angle,
		Satelites: satelites,
		Speed:     speed,
	}
}
