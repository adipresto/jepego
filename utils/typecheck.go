package utils

type DataType int

const (
	TypeUnknown DataType = iota
	TypeString
	TypeNumber
	TypeObject
	TypeArray
	TypeBool
	TypeNull
)

func detectType(raw []byte) DataType {
	if len(raw) == 0 {
		return TypeUnknown
	}
	switch raw[0] {
	case '"':
		return TypeString
	case '{':
		return TypeObject
	case '[':
		return TypeArray
	case 't', 'f':
		return TypeBool
	case 'n':
		return TypeNull
	default:
		return TypeNumber
	}
}
