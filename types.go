package jepego

type DataType int

const (
	TypeUnknown DataType = iota
	TypeString
	TypeNumber
	TypeObject
	TypeArray
	TypeBool
	TypeRaw
	TypeNull
)
