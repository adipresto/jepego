package jsonParser_v3

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

type Result struct {
	Key      string
	Data     []byte
	DataType DataType
	OK       bool
}

// -------------------- Public API --------------------

// Get mengambil value berdasarkan path (contoh: "a.b[3].c").
func Get(json []byte, path string) Result {
	json = removeCommentsBytes(json)
	if len(path) == 0 {
		return Result{Key: "", OK: false}
	}

	val, ok := getNestedValue(json, splitPathBytes([]byte(path)))
	if !ok {
		return Result{Key: path, OK: false}
	}

	var resVal []byte
	if val[0] == '"' {
		resVal = val[1 : len(val)-1]
	} else {
		resVal = val
	}

	return Result{
		Key:      path,
		Data:     resVal,
		DataType: detectType(val),
		OK:       true,
	}
}

// GetMany menerima ekspresi dalam bentuk []string{"a","b","c.d","arr[5].x"}.
// Mengembalikan hasil sesuai urutan field yang diminta.
func GetMany(json []byte, exprs []string) []Result {
	json = removeCommentsBytes(json)
	out := make([]Result, 0, len(exprs))
	for _, p := range exprs {
		if len(p) == 0 {
			out = append(out, Result{Key: "", OK: false})
			continue
		}

		val, ok := getNestedValue(json, splitPathBytes([]byte(p)))
		if !ok {
			out = append(out, Result{Key: p, OK: false})
			continue
		}

		var resVal []byte
		if val[0] == '"' {
			resVal = val[1 : len(val)-1]
		} else {
			resVal = val
		}

		out = append(out, Result{
			Key:      p,
			Data:     resVal,
			DataType: detectType(val),
			OK:       true,
		})
	}
	return out
}

// -------------------- Core traversal --------------------

type pathToken struct {
	key   []byte
	index int
	isIdx bool
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

// extractValue: ambil 1 JSON value utuh dari s[0:], return (slice, bytesConsumed)
func extractValue(s []byte) ([]byte, int) {
	s = trimSpaceBytes(s)
	if len(s) == 0 {
		return nil, 0
	}
	switch s[0] {
	case '{':
		depth := 0
		inStr := false
		esc := false
		for i := 0; i < len(s); i++ {
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
		depth := 0
		inStr := false
		esc := false
		for i := 0; i < len(s); i++ {
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
		for j > 0 && isSpace(s[j-1]) {
			j--
		}
		return s[:j], i
	}
	return nil, 0
}

// -------------------- Path & multi expr --------------------

// splitPathBytes: "a.b[3].c" -> tokens {key:"a"}, {idx:3}, {key:"c"}
func splitPathBytes(path []byte) []pathToken {
	var toks []pathToken
	var buf []byte
	for i := 0; i < len(path); i++ {
		c := path[i]
		switch c {
		case '.':
			if len(buf) > 0 {
				toks = append(toks, pathToken{key: cloneBytes(buf)})
				buf = buf[:0]
			}
		case '[':
			if len(buf) > 0 {
				toks = append(toks, pathToken{key: cloneBytes(buf)})
				buf = buf[:0]
			}
			// parse index
			i++
			idx := 0
			for i < len(path) && path[i] != ']' {
				idx = idx*10 + int(path[i]-'0')
				i++
			}
			toks = append(toks, pathToken{isIdx: true, index: idx})
		default:
			buf = append(buf, c)
		}
	}
	if len(buf) > 0 {
		toks = append(toks, pathToken{key: cloneBytes(buf)})
	}
	return toks
}

// parseMultiExpr: "{a,b,c.d,arr[5].x}" -> [][]byte{"a","b","c.d","arr[5].x"}
func parseMultiExpr(expr []byte) [][]byte {
	b := trimSpaceBytes(expr)
	if len(b) >= 2 && b[0] == '{' && b[len(b)-1] == '}' {
		b = b[1 : len(b)-1]
	}
	return splitCSVTopLevel(b)
}

// -------------------- Utils --------------------

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

func removeCommentsBytes(data []byte) []byte {
	// Hapus komentar "// ..." per-baris. (Komentar block tidak didukung)
	out := make([]byte, 0, len(data))
	n := len(data)
	i := 0
	for i < n {
		// deteksi // (pastikan tidak sedang di dalam string)
		if i+1 < n && data[i] == '/' && data[i+1] == '/' {
			i += 2
			for i < n && data[i] != '\n' {
				i++
			}
			continue
		}
		out = append(out, data[i])
		i++
	}
	return trimSpaceBytes(out)
}

// skipWS: majuin index selama whitespace
func skipWS(i *int, b []byte) {
	for *i < len(b) && isSpace(b[*i]) {
		*i++
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
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

// splitCSVTopLevel: pecah "a,b,c" jadi potongan, asumsi nama field tidak mengandung koma/kutip.
func splitCSVTopLevel(b []byte) [][]byte {
	var parts [][]byte
	start := 0
	for i := 0; i <= len(b); i++ {
		if i == len(b) || b[i] == ',' {
			seg := trimSpaceBytes(b[start:i])
			parts = append(parts, cloneBytes(seg))
			start = i + 1
		}
	}
	return parts
}

func trimSpaceBytes(b []byte) []byte {
	start, end := 0, len(b)
	for start < end && isSpace(b[start]) {
		start++
	}
	for end > start && isSpace(b[end-1]) {
		end--
	}
	return b[start:end]
}

func cloneBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
