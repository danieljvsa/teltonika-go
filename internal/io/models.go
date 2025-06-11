package io

type IOData struct {
	IO    int64
	Value string
}

type ResponseDecode struct {
	IOs            []IOData
	NumberOfIOs    int64
	LastByte       int64
	GenerationType string
}
