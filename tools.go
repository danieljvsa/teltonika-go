package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type GPSData struct {
	Latitude  float64
	Longitude float64
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
	if len(data) < 14 {
		fmt.Println("Invalid data length", len(data))
		return nil
	}

	lat, err := strconv.ParseUint(hex.EncodeToString(data[4:8]), 16, 32)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
		return nil
	}

	long, err := strconv.ParseUint(hex.EncodeToString(data[:4]), 16, 32)
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

	speed, err := strconv.ParseInt(hex.EncodeToString(data[13:14]), 16, 64)
	if err != nil {
		fmt.Println("Error parsing speed:", err)
		return nil
	}

	longitude := float64(int32(long)) / 10000000.0
	latitude := float64(int32(lat)) / 10000000.0

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
