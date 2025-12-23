package jepego

import "fmt"

func upsertNested(cur []byte, tokens []pathToken, value []byte, valueType DataType) ([]byte, error) {
	if len(tokens) == 0 {
		// reached leaf → return value as JSON fragment
		switch valueType {
		case TypeString:
			return append(append([]byte{'"'}, escapeString(value)...), '"'), nil
		case TypeNumber, TypeBool, TypeNull:
			return value, nil
		case TypeObject, TypeArray:
			return value, nil
		case TypeRaw:
			return value, nil
		default:
			return append(append([]byte{'"'}, escapeString(value)...), '"'), nil
		}
	}

	p := tokens[0]
	rest := tokens[1:]

	if len(cur) == 0 {
		// init object or array if missing
		if p.isIdx || p.isWildcard {
			cur = []byte("[]")
		} else {
			cur = []byte("{}")
		}
	}

	switch cur[0] {
	case '{':
		// object
		v, ok := getTopLevelKey(cur, p.key)
		if ok {
			newVal, err := upsertNested(v, rest, value, valueType)
			if err != nil {
				return nil, err
			}
			return replaceObjectKey(cur, p.key, newVal)
		} else {
			// key missing → add new
			newLeaf, err := upsertNested(nil, rest, value, valueType)
			if err != nil {
				return nil, err
			}
			return addObjectKey(cur, p.key, newLeaf)
		}

	case '[':
		// array
		if !p.isIdx {
			return nil, fmt.Errorf("cannot upsert wildcard directly")
		}
		v, ok := getArrayIndex(cur, p.index)
		if ok {
			newVal, err := upsertNested(v, rest, value, valueType)
			if err != nil {
				return nil, err
			}
			return replaceArrayIndex(cur, p.index, newVal)
		} else {
			// pad array sampai index → insert value
			return padAndInsertArray(cur, p.index, rest, value, valueType)
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
		skipWS(&i, obj)
		if i >= n || obj[i] == '}' {
			break
		}
		if obj[i] != '"' {
			return nil, fmt.Errorf("invalid object")
		}
		i++
		keyStart := i
		ok, keyEnd := scanString(obj, i)
		if !ok {
			return nil, fmt.Errorf("unterminated key")
		}
		i = keyEnd + 1
		skipWS(&i, obj)
		if i >= n || obj[i] != ':' {
			return nil, fmt.Errorf("expected ':' after key")
		}
		i++
		skipWS(&i, obj)
		val, consumed := extractValue(obj[i:])
		if consumed == 0 {
			return nil, fmt.Errorf("invalid value")
		}

		if bytesEqual(obj[keyStart:keyEnd], key) {
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
		skipWS(&i, obj)
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
		skipWS(&i, arr)
		if i >= n || arr[i] == ']' {
			break
		}
		elem, consumed := extractValue(arr[i:])
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
		skipWS(&i, arr)
		if i < n && arr[i] == ',' {
			i++
		}
		curIdx++
	}
	out = append(out, ']')
	return out, nil
}

// padAndInsertArray: jika index > len, pad dengan null
func padAndInsertArray(arr []byte, idx int, rest []pathToken, value []byte, valueType DataType) ([]byte, error) {
	i := 1
	n := len(arr)
	elements := [][]byte{}
	curIdx := 0

	for i < n {
		skipWS(&i, arr)
		if i >= n || arr[i] == ']' {
			break
		}
		elem, consumed := extractValue(arr[i:])
		if consumed == 0 {
			return nil, fmt.Errorf("invalid array")
		}
		elements = append(elements, elem)
		i += consumed
		skipWS(&i, arr)
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
