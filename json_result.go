package jepego

type Result struct {
	Key      string
	Data     []byte
	DataType DataType
	OK       bool
}
