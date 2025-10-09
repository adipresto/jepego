package utils

// getTopLevelKey: ambil value untuk key top-level pada sebuah object JSON.
func GetTopLevelKey(obj []byte, key []byte) ([]byte, bool) {
	obj = TrimSpaceBytes(obj)
	if len(obj) == 0 || obj[0] != '{' {
		return nil, false
	}
	i := 1
	n := len(obj)

	for i < n {
		SkipWS(&i, obj)
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
		SkipWS(&i, obj)
		if i >= n || obj[i] != ':' {
			return nil, false
		}
		i++ // skip ':'
		SkipWS(&i, obj)
		if i >= n {
			return nil, false
		}

		// ambil value
		val, consumed := extractValue(obj[i:])
		if consumed == 0 {
			return nil, false
		}

		// bandingkan key (tanpa quotes)
		if BytesEqual(obj[keyStart:keyEnd], key) {
			return val, true
		}

		// lompat ke setelah value
		i += consumed
		SkipWS(&i, obj)
		if i < n && obj[i] == ',' {
			i++
			continue
		}
		// kalau bukan koma, mungkin '}' atau whitespace; loop lanjut.
	}
	return nil, false
}

// getArrayIndex: ambil elemen array ke-want (0-based).
func GetArrayIndex(arr []byte, want int) ([]byte, bool) {
	arr = TrimSpaceBytes(arr)
	if len(arr) == 0 || arr[0] != '[' {
		return nil, false
	}
	i := 1
	n := len(arr)
	idx := 0
	for i < n {
		SkipWS(&i, arr)
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
		SkipWS(&i, arr)
		if i < n && arr[i] == ',' {
			i++
		}
		idx++
	}
	return nil, false
}

// extractValue: ambil 1 JSON value utuh dari s[0:], return (slice, bytesConsumed)
func extractValue(s []byte) ([]byte, int) {
	s = TrimSpaceBytes(s)
	if len(s) == 0 {
		return nil, 0
	}
	switch s[0] {
	case '{':
		depth := 1
		inStr := false
		esc := false
		for i := 1; i < len(s); i++ { // mulai dari 1 karena s[0] == '{'
			c := s[i]
			if inStr {
				if esc {
					esc = false
					continue
				}
				if c == '\\' {
					esc = true
					continue
				}
				if c == '"' {
					inStr = false
				}
				continue
			}
			if c == '"' {
				inStr = true
				continue
			}
			if c == '{' {
				depth++
			} else if c == '}' {
				depth--
				if depth == 0 {
					return s[:i+1], i + 1
				}
			}
		}
	case '[':
		depth := 1
		inStr := false
		esc := false
		for i := 1; i < len(s); i++ { // mulai dari 1 karena s[0] == '['
			c := s[i]
			if inStr {
				if esc {
					esc = false
					continue
				}
				if c == '\\' {
					esc = true
					continue
				}
				if c == '"' {
					inStr = false
				}
				continue
			}
			if c == '"' {
				inStr = true
				continue
			}
			if c == '[' {
				depth++
			} else if c == ']' {
				depth--
				if depth == 0 {
					return s[:i+1], i + 1
				}
			}
		}
	case '"':
		ok, end := scanString(s, 1)
		if !ok {
			return nil, 0
		}
		return s[:end+1], end + 1
	default:
		// number / true / false / null
		i := 0
		for i < len(s) {
			c := s[i]
			if c == ',' || c == '}' || c == ']' {
				break
			}
			i++
		}
		// trim trailing spaces di dalam token
		j := i
		for j > 0 && IsSpace(s[j-1]) {
			j--
		}
		return s[:j], i
	}
	return nil, 0
}

func unwrap(val []byte) []byte {
	if len(val) > 1 && val[0] == '"' {
		return val[1 : len(val)-1]
	}
	return val
}

// SkipWS: majuin index selama whitespace
func SkipWS(i *int, b []byte) {
	for *i < len(b) && IsSpace(b[*i]) {
		*i++
	}
}

// scanString: mulai dari s[start], di mana s[start-1] == '"', cari closing quote (handle escape)
func scanString(s []byte, start int) (ok bool, end int) {
	esc := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if esc {
			esc = false
			continue
		}
		if c == '\\' {
			esc = true
			continue
		}
		if c == '"' {
			return true, i
		}
	}
	return false, 0
}
