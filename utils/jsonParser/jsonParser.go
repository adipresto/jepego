package jsonParser

import (
	"encoding/json"
	"fmt"
)

// Get ambil satu key dari JSON pakai unmarshal
func Get(data []byte, key string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", false
	}

	if val, ok := m[key]; ok {
		return toString(val), true
	}
	return "", false
}

// GetMany ambil beberapa key sekaligus
func GetMany(data []byte, keys []string) map[string]string {
	var m map[string]interface{}
	_ = json.Unmarshal(data, &m) // kalau gagal, hasil kosong

	result := make(map[string]string, len(keys))
	for _, k := range keys {
		if val, ok := m[k]; ok {
			result[k] = toString(val)
		}
	}
	return result
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%v", t)
	case bool:
		return fmt.Sprintf("%v", t)
	default:
		return ""
	}
}
