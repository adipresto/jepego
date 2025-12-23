package jepego

// -------------------- Core traversal --------------------

type pathToken struct {
	key        []byte
	index      int
	isIdx      bool
	isWildcard bool
}

// getNestedValue: berjalan mengikuti tokens (field / [index])
func getNestedValue(json []byte, parts []pathToken) ([]byte, bool) {
	cur := trimSpaceBytes(json)
	for _, p := range parts {
		if len(cur) == 0 {
			return nil, false
		}
		switch cur[0] {
		case '{':
			if p.isIdx {
				return nil, false
			}
			v, ok := getTopLevelKey(cur, p.key)
			if !ok {
				return nil, false
			}
			cur = trimSpaceBytes(v)
		case '[':
			if !p.isIdx {
				return nil, false
			}
			v, ok := getArrayIndex(cur, p.index)
			if !ok {
				return nil, false
			}
			cur = trimSpaceBytes(v)
		default:
			return nil, false
		}
	}
	return cur, true
}

// getNestedValues: berjalan mengikuti tokens (field / [index / []]) dan mengembalikan semua hasil yang ditemukan.
func getNestedValues(json []byte, parts []pathToken, fullPath string) []Result {
	cur := trimSpaceBytes(json)
	if len(parts) == 0 {
		return []Result{{
			Key:      fullPath,
			Data:     unwrap(cur),
			DataType: detectType(cur),
			OK:       true,
		}}
	}

	p := parts[0]
	rest := parts[1:]

	switch cur[0] {
	case '{':
		if p.isIdx {
			return nil
		}
		v, ok := getTopLevelKey(cur, p.key)
		if !ok {
			return nil
		}
		return getNestedValues(v, rest, fullPath)

	case '[':
		if p.isIdx {
			// index spesifik
			v, ok := getArrayIndex(cur, p.index)
			if !ok {
				return nil
			}
			return getNestedValues(v, rest, fullPath)
		}

		// wildcard []
		if len(p.key) == 0 && !p.isIdx {
			var results []Result
			i := 1
			n := len(cur)
			for i < n {
				skipWS(&i, cur)
				if i >= n || cur[i] == ']' {
					break
				}
				val, consumed := extractValue(cur[i:])
				if consumed == 0 {
					break
				}
				sub := getNestedValues(val, rest, fullPath)
				results = append(results, sub...)
				i += consumed
				skipWS(&i, cur)
				if i < n && cur[i] == ',' {
					i++
				}
			}
			return results
		}
	}
	return nil
}

func unwrap(val []byte) []byte {
	if len(val) > 1 && val[0] == '"' {
		return val[1 : len(val)-1]
	}
	return val
}

// getTopLevelKey: ambil value untuk key top-level pada sebuah object JSON.
func getTopLevelKey(obj []byte, key []byte) ([]byte, bool) {
	obj = trimSpaceBytes(obj)
	if len(obj) == 0 || obj[0] != '{' {
		return nil, false
	}
	i := 1
	n := len(obj)

	for i < n {
		skipWS(&i, obj)
		if i >= n {
			break
		}
		if obj[i] == '}' {
			return nil, false
		}

		// key harus string
		if obj[i] != '"' {
			// JSON tidak valid untuk kasus ini; gagal.
			return nil, false
		}
		i++ // ke isi string
		keyStart := i
		ok, keyEnd := scanString(obj, i) // keyEnd = index dari closing quote
		if !ok {
			return nil, false
		}
		i = keyEnd + 1 // pos setelah closing quote
		skipWS(&i, obj)
		if i >= n || obj[i] != ':' {
			return nil, false
		}
		i++ // skip ':'
		skipWS(&i, obj)
		if i >= n {
			return nil, false
		}

		// ambil value
		val, consumed := extractValue(obj[i:])
		if consumed == 0 {
			return nil, false
		}

		// bandingkan key (tanpa quotes)
		if bytesEqual(obj[keyStart:keyEnd], key) {
			return val, true
		}

		// lompat ke setelah value
		i += consumed
		skipWS(&i, obj)
		if i < n && obj[i] == ',' {
			i++
			continue
		}
		// kalau bukan koma, mungkin '}' atau whitespace; loop lanjut.
	}
	return nil, false
}

// getArrayIndex: ambil elemen array ke-want (0-based).
func getArrayIndex(arr []byte, want int) ([]byte, bool) {
	arr = trimSpaceBytes(arr)
	if len(arr) == 0 || arr[0] != '[' {
		return nil, false
	}
	i := 1
	n := len(arr)
	idx := 0
	for i < n {
		skipWS(&i, arr)
		if i >= n {
			break
		}
		if arr[i] == ']' {
			return nil, false
		}
		val, consumed := extractValue(arr[i:])
		if consumed == 0 {
			return nil, false
		}
		if idx == want {
			return val, true
		}
		i += consumed
		skipWS(&i, arr)
		if i < n && arr[i] == ',' {
			i++
		}
		idx++
	}
	return nil, false
}
