package jsonParser_v2

import (
	"fmt"
	"strconv"
	"strings"
)

// Result menyimpan JSON string
type Result struct {
	Raw        string
	Collection []string
}

type ParserOption struct {
	IsRaw bool
}

func removeComments(jsonStr string) string {
	lines := strings.Split(jsonStr, "\n")

	var b strings.Builder
	for _, line := range lines {
		if line != "" {
			if idx := strings.Index(line, "//"); idx != -1 {
				line = line[:idx]
			}
			line = strings.TrimSpace(line)
			b.WriteString(line)
		}
	}
	return strings.TrimSpace(b.String())
}

// Get mengambil nilai JSON berdasarkan path tanpa unmarshal
func Get(jsonStr string, path string, args ...ParserOption) Result {
	// sest default option
	o := ParserOption{
		IsRaw: false,
	}
	// kalau ada opt
	if len(args) > 0 {
		o = args[0]
	}

	jsonStr = removeComments(jsonStr)

	// Modifier @ / !
	if strings.HasPrefix(path, "@") || strings.HasPrefix(path, "!") {
		val, ok := getNestedValue(jsonStr, splitPath(path[1:]))
		if val == "" || !ok {
			if o.IsRaw {
				return Result{Raw: "null"}
			}
			return Result{
				Collection: []string{"null"},
			}
		}
		if o.IsRaw {
			return Result{Raw: val}
		}
		return Result{Collection: []string{val}}
	}

	// Subselector object {field1,field2,...}
	// support nested path & array
	if strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}") {
		fields := strings.Split(path[1:len(path)-1], ",")
		sub := []string{}
		for _, f := range fields {
			f = strings.TrimSpace(f)
			val, ok := getNestedValue(jsonStr, splitPath(f))
			if !ok {
				continue
			}
			// Jika value object/array, wrap jadi string JSON
			if val != "" && (val[0] == '{' || val[0] == '[') {
				val = fmt.Sprintf(`"%s"`, val)
			}

			if o.IsRaw {
				// disimpan sebagai keyname:value
				sub = append(sub, fmt.Sprintf(`"%s":%s`, f, val))
			} else {
				// disimpan hanya value
				sub = append(sub, val)
			}
		}
		if len(sub) == 0 {
			if o.IsRaw {
				return Result{Raw: "null"}
			}
			return Result{Collection: []string{"null"}}
		}
		if o.IsRaw {
			return Result{Raw: "{" + strings.Join(sub, ",") + "}"}
		}
		return Result{Collection: sub}
	}

	// Subselector array [0,2]
	if strings.HasPrefix(path, "[") && strings.HasSuffix(path, "]") {
		parts := strings.Split(path[1:len(path)-1], ",")
		arrVals := []string{}

		if jsonStr[0] == '[' {
			arr := getTopLevelArray(jsonStr)
			for _, s := range parts {
				s = strings.TrimSpace(s)
				if idx, err := strconv.Atoi(s); err == nil && idx >= 0 && idx < len(arr) {
					arrVals = append(arrVals, arr[idx])
				}
			}
		} else {
			for _, s := range parts {
				s = strings.TrimSpace(s)
				val, ok := getNestedValue(jsonStr, splitPath(s))

				if !ok {
					continue
				}

				// Jika value object/array, wrap jadi string JSON
				if val != "" && (val[0] == '{' || val[0] == '[') {
					val = fmt.Sprintf(`"%s"`, val)
				}

				if o.IsRaw {
					// disimpan sebagai keyname:value
					arrVals = append(arrVals, fmt.Sprintf(`"%s":%s`, s, val))
				} else {
					// disimpan hanya value
					arrVals = append(arrVals, val)
				}
			}
		}
		if len(arrVals) == 0 {
			if o.IsRaw {
				return Result{Raw: "null"}
			}
			return Result{Collection: []string{"null"}}
		}
		if o.IsRaw {
			return Result{Raw: "[" + strings.Join(arrVals, ",") + "]"}
		}
		return Result{Collection: arrVals}
	}

	// Nested path (a.b[0].c)
	val, ok := getNestedValue(jsonStr, splitPath(path))
	if !ok {
		if o.IsRaw {
			return Result{Raw: "null"}
		}
		return Result{Collection: []string{"null"}}
	}
	if o.IsRaw {
		return Result{Raw: val}
	}
	return Result{Collection: []string{val}}
}

// =================== Helper Functions =======================

// Mendapatkan key top-level di object JSON
func getTopLevelKey(jsonStr, key string) string {
	// cache panjang jsonStr
	jsonStrLen := len(jsonStr)
	pos := 0
	for pos < jsonStrLen {
		c := jsonStr[pos]
		if c == '"' {
			// buat posisi awal
			start := pos + 1
			pos++
			for pos < jsonStrLen && jsonStr[pos] != '"' {
				pos++
			}

			// ambil key dari posisi start dan saat ini
			k := jsonStr[start:pos]
			pos++
			// cari ':'
			for pos < jsonStrLen && (jsonStr[pos] == ' ' || jsonStr[pos] == ':') {
				if jsonStr[pos] == ':' {
					pos++
					break
				}
				pos++
			}
			// jika key cocok
			if k == key {
				// _ harusnya next
				val, _ := extractValue(jsonStr[pos:])
				return val
			}
		} else {
			pos++
		}
	}
	return ""
}

// Mendapatkan array top-level sebagai slice string
func getTopLevelArray(jsonStr string) []string {
	jsonStr = strings.TrimSpace(jsonStr)
	if len(jsonStr) == 0 || jsonStr[0] != '[' {
		return nil
	}
	var arr []string
	pos := 1
	for pos < len(jsonStr) {
		val, next := extractValue(jsonStr[pos:])
		if val != "" {
			arr = append(arr, val)
		}
		pos += next
		if pos >= len(jsonStr) {
			break
		}
		if jsonStr[pos] == ',' {
			pos++
		} else if jsonStr[pos] == ']' {
			break
		}
	}
	return arr
}

// Ambil value JSON dari posisi awal string (object, array, string, number, boolean, null)
func extractValue(s string) (string, int) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", 0
	}
	switch s[0] {
	case '{':
		// ambil object utuh
		count := 0
		for i, c := range s {
			if c == '{' {
				count++
			} else if c == '}' {
				count--
				if count == 0 {
					return s[:i+1], i + 1
				}
			} else if c == '"' {
				// skip string
				i++
			}
		}
	case '[':
		count := 0
		for i, c := range s {
			if c == '[' {
				count++
			} else if c == ']' {
				count--
				if count == 0 {
					return s[:i+1], i + 1
				}
			} else if c == '"' {
				i++
			}
		}
	case '"':
		sLen := len(s)
		// ambil string literal
		for i := 1; i < sLen; i++ {
			if s[i] == '"' {
				return s[:i+1], i + 1
			}
		}
	default:
		sLen := len(s)
		// number, boolean, null
		for i := 0; i < sLen; i++ {
			if s[i] == ',' || s[i] == '}' || s[i] == ']' {
				return strings.TrimSpace(s[:i]), i
			}
		}
		return strings.TrimSpace(s), sLen
	}
	return "", 0
}

// Split path nested (a.b[0].c -> ["a","b","0","c"])
func splitPath(path string) []string {
	parts := []string{}
	buf := ""
	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '.' {
			if buf != "" {
				parts = append(parts, buf)
				buf = ""
			}
		} else if c == '[' {
			if buf != "" {
				parts = append(parts, buf)
				buf = ""
			}
			j := i + 1
			for ; j < len(path) && path[j] != ']'; j++ {
			}
			parts = append(parts, path[i+1:j])
			i = j
		} else {
			buf += string(c)
		}
	}
	if buf != "" {
		parts = append(parts, buf)
	}
	return parts
}

// Ambil nested value berdasarkan slice path
func getNestedValue(jsonStr string, parts []string) (string, bool) {
	current := jsonStr
	for _, p := range parts {
		current = strings.TrimSpace(current)
		if current == "" {
			return "", false
		}
		if current[0] == '{' {
			v := getTopLevelKey(current, p)
			if v == "" {
				return "", false
			}
			current = v
		} else if current[0] == '[' {
			arr := getTopLevelArray(current)
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 || idx >= len(arr) {
				return "", false
			}
			current = arr[idx]
		} else {
			return "", false
		}
	}
	return current, true
}
