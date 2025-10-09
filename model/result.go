package constant

import (
	jpbyte "github.com/adipresto/jepego/model/utils/byte"
)

type Result struct {
	Key      string
	Data     []byte
	DataType jpbyte.DataType
	OK       bool
}
