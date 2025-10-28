package tool

type GPSData struct {
	Latitude  float64
	Longitude float64
	Altitude  int64
	Angle     int64
	Satelites int64
	Speed     int64
}

type LoginData struct {
	Length int64
	IMEI   string
}

type ProtocolData struct {
	Protocol string
}

type CommandResponse struct {
	Response    string
	HexMessage  string
	CommandType string
	IMEI        string
}
