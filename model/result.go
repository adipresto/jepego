package model

import (
	"github.com/adipresto/jepego/utils"
)

type Result struct {
	Key      string
	Data     []byte
	DataType utils.DataType
	OK       bool
}
