package jsonParser_v3

import (
	"strconv"
)

// Hasil pencarian
type Result struct {
	Raw   []byte
	Str   string
	Valid bool
}

// Remove comments (// ...) dari JSON bytes
func removeComments(data []byte) []byte {
	out := make([]byte, 0, len(data))
	n := len(data)
	i := 0
	for i < n {
		// kalau ketemu //
		if i+1 < n && data[i] == '/' && data[i+1] == '/' {
			// skip sampai newline
			i += 2
			for i < n && data[i] != '\n' {
				i++
			}
			continue
		}
		out = append(out, data[i])
		i++
	}
	return out
}

// Get ambil field dari JSON object satu level
// contoh: {"name":"Alice","age":20}
// Get(data, "name") -> "Alice"
func Get(data []byte, field string) Result {
	data = removeComments(data)

	n := len(data)
	i := 0

	// cari {
	for i < n && data[i] != '{' {
		i++
	}
	if i >= n {
		return Result{Valid: false}
	}
	i++

	for i < n {
		// skip spasi
		for i < n && (data[i] == ' ' || data[i] == '\n' || data[i] == '\t' || data[i] == ',') {
			i++
		}
		if i >= n || data[i] == '}' {
			break
		}

		// baca key
		if data[i] != '"' {
			return Result{Valid: false}
		}
		i++
		keyStart := i
		for i < n && data[i] != '"' {
			i++
		}
		if i >= n {
			return Result{Valid: false}
		}
		keyEnd := i
		key := data[keyStart:keyEnd]
		i++ // skip closing "

		// skip spasi + colon
		for i < n && (data[i] == ' ' || data[i] == ':') {
			i++
		}
		if i >= n {
			return Result{Valid: false}
		}

		// ambil value
		var valStart, valEnd int
		if data[i] == '"' { // string value
			i++
			valStart = i
			for i < n && data[i] != '"' {
				i++
			}
			if i >= n {
				return Result{Valid: false}
			}
			valEnd = i
			i++ // skip closing "
		} else { // number / bool / null
			valStart = i
			for i < n && data[i] != ',' && data[i] != '}' {
				i++
			}
			valEnd = i
		}

		// cocokkan key
		if string(key) == field {
			val := data[valStart:valEnd]
			return Result{
				Raw:   val,
				Str:   string(val),
				Valid: true,
			}
		}
	}
	return Result{Valid: false}
}

// contoh parsing integer
func (r Result) Int() (int, bool) {
	if !r.Valid {
		return 0, false
	}
	v, err := strconv.Atoi(r.Str)
	if err != nil {
		return 0, false
	}
	return v, true
}
