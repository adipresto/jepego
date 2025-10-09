package apis

import (
	"fmt"

	"github.com/adipresto/jepego/utils"
)

// Upsert: update atau insert value di path tertentu
func Upsert(json []byte, path string, value []byte, valueType utils.DataType) ([]byte, error) {
	if len(path) == 0 {
		return json, fmt.Errorf("empty path")
	}

	tokens := utils.SplitPathBytes([]byte(path))
	return upsertNested(json, tokens, value, valueType)
}

func upsertNested(cur []byte, tokens []utils.PathToken, value []byte, valueType utils.DataType) ([]byte, error) {
	if len(tokens) == 0 {
		// reached leaf → return value as JSON fragment
		switch valueType {
		case utils.TypeString:
			return append(append([]byte{'"'}, utils.EscapeString(value)...), '"'), nil
		case utils.TypeNumber, utils.TypeBool, utils.TypeNull:
			return value, nil
		case utils.TypeObject, utils.TypeArray:
			return value, nil
		default:
			return append(append([]byte{'"'}, utils.EscapeString(value)...), '"'), nil
		}
	}

	p := tokens[0]
	rest := tokens[1:]

	if len(cur) == 0 {
		// init object or array if missing
		if p.IsIdx || p.IsWildcard {
			cur = []byte("[]")
		} else {
			cur = []byte("{}")
		}
	}

	switch cur[0] {
	case '{':
		// object
		v, ok := utils.GetTopLevelKey(cur, p.Key)
		if ok {
			newVal, err := upsertNested(v, rest, value, valueType)
			if err != nil {
				return nil, err
			}
			return replaceObjectKey(cur, p.Key, newVal)
		} else {
			// key missing → add new
			newLeaf, err := upsertNested(nil, rest, value, valueType)
			if err != nil {
				return nil, err
			}
			return addObjectKey(cur, p.Key, newLeaf)
		}

	case '[':
		// array
		if !p.IsIdx {
			return nil, fmt.Errorf("cannot upsert wildcard directly")
		}
		v, ok := utils.GetArrayIndex(cur, p.Index)
		if ok {
			newVal, err := upsertNested(v, rest, value, valueType)
			if err != nil {
				return nil, err
			}
			return replaceArrayIndex(cur, p.Index, newVal)
		} else {
			// pad array sampai index → insert value
			return padAndInsertArray(cur, p.Index, rest, value, valueType)
		}
	default:
		return nil, fmt.Errorf("unexpected JSON type at this path")
	}
}

// replaceObjectKey: ganti value key existing
func replaceObjectKey(obj []byte, key []byte, newVal []byte) ([]byte, error) {
	i := 1
	n := len(obj)
	out := make([]byte, 0, len(obj)+len(newVal))
	out = append(out, '{')

	for i < n {
		utils.SkipWS(&i, obj)
		if i >= n || obj[i] == '}' {
			break
		}
		if obj[i] != '"' {
			return nil, fmt.Errorf("invalid object")
		}
		i++
		keyStart := i
		ok, keyEnd := utils.ScanString(obj, i)
		if !ok {
			return nil, fmt.Errorf("unterminated key")
		}
		i = keyEnd + 1
		utils.SkipWS(&i, obj)
		if i >= n || obj[i] != ':' {
			return nil, fmt.Errorf("expected ':' after key")
		}
		i++
		utils.SkipWS(&i, obj)
		val, consumed := utils.ExtractValue(obj[i:])
		if consumed == 0 {
			return nil, fmt.Errorf("invalid value")
		}

		if utils.BytesEqual(obj[keyStart:keyEnd], key) {
			// key match → tulis key + newVal
			if len(out) > 1 {
				out = append(out, ',')
			}
			out = append(out, '"')
			out = append(out, obj[keyStart:keyEnd]...)
			out = append(out, '"', ':')
			out = append(out, newVal...)
		} else {
			if len(out) > 1 {
				out = append(out, ',')
			}
			out = append(out, '"')
			out = append(out, obj[keyStart:keyEnd]...)
			out = append(out, '"', ':')
			out = append(out, val...)
		}

		i += consumed
		utils.SkipWS(&i, obj)
		if i < n && obj[i] == ',' {
			i++
		}
	}
	out = append(out, '}')
	return out, nil
}

// addObjectKey: append key baru
func addObjectKey(obj []byte, key []byte, val []byte) ([]byte, error) {
	if len(obj) < 2 || obj[0] != '{' || obj[len(obj)-1] != '}' {
		return nil, fmt.Errorf("invalid object")
	}
	out := make([]byte, 0, len(obj)+len(key)+len(val)+4)
	out = append(out, obj[:len(obj)-1]...) // copy tanpa '}'

	if len(obj) > 2 {
		out = append(out, ',')
	}

	out = append(out, '"')
	out = append(out, key...)
	out = append(out, '"', ':')
	out = append(out, val...)
	out = append(out, '}')
	return out, nil
}

// replaceArrayIndex: ganti elemen array existing
func replaceArrayIndex(arr []byte, idx int, val []byte) ([]byte, error) {
	i := 1
	n := len(arr)
	out := []byte{'['}
	curIdx := 0

	for i < n {
		utils.SkipWS(&i, arr)
		if i >= n || arr[i] == ']' {
			break
		}
		elem, consumed := utils.ExtractValue(arr[i:])
		if consumed == 0 {
			return nil, fmt.Errorf("invalid array")
		}

		if curIdx > 0 {
			out = append(out, ',')
		}
		if curIdx == idx {
			out = append(out, val...)
		} else {
			out = append(out, elem...)
		}

		i += consumed
		utils.SkipWS(&i, arr)
		if i < n && arr[i] == ',' {
			i++
		}
		curIdx++
	}
	out = append(out, ']')
	return out, nil
}

// padAndInsertArray: jika index > len, pad dengan null
func padAndInsertArray(arr []byte, idx int, rest []utils.PathToken, value []byte, valueType utils.DataType) ([]byte, error) {
	i := 1
	n := len(arr)
	elements := [][]byte{}
	curIdx := 0

	for i < n {
		utils.SkipWS(&i, arr)
		if i >= n || arr[i] == ']' {
			break
		}
		elem, consumed := utils.ExtractValue(arr[i:])
		if consumed == 0 {
			return nil, fmt.Errorf("invalid array")
		}
		elements = append(elements, elem)
		i += consumed
		utils.SkipWS(&i, arr)
		if i < n && arr[i] == ',' {
			i++
		}
		curIdx++
	}

	// pad sampai idx
	for len(elements) <= idx {
		elements = append(elements, []byte("null"))
	}

	// upsert leaf value
	newVal, err := upsertNested(elements[idx], rest, value, valueType)
	if err != nil {
		return nil, err
	}
	elements[idx] = newVal

	// rebuild array
	out := []byte{'['}
	for j, e := range elements {
		if j > 0 {
			out = append(out, ',')
		}
		out = append(out, e...)
	}
	out = append(out, ']')
	return out, nil
}
