package apis

import (
	"fmt"

	"github.com/adipresto/jepego/utils"
)

// ValidateJSONRobust memeriksa apakah data JSON valid.
// Mengembalikan error dengan info posisi byte kalau ada kesalahan.
func ValidateJSON(data []byte) error {
	data = utils.RemoveCommentsBytes(data)
	data = utils.TrimSpaceBytes(data)
	if len(data) == 0 {
		return fmt.Errorf("empty JSON")
	}

	i := 0
	if err := parseValue(data, &i); err != nil {
		return fmt.Errorf("invalid JSON at byte %d: %v", i, err)
	}

	// pastikan tidak ada sisa byte setelah value valid
	utils.SkipWS(&i, data)
	if i != len(data) {
		return fmt.Errorf("invalid JSON: extra content after valid value at byte %d", i)
	}

	return nil
}

// parseValue: parse 1 JSON value dari data[i:], update i ke posisi setelah value
func parseValue(data []byte, i *int) error {
	utils.SkipWS(i, data)
	if *i >= len(data) {
		return fmt.Errorf("unexpected end of data")
	}

	switch data[*i] {
	case '{':
		return parseObject(data, i)
	case '[':
		return parseArray(data, i)
	case '"':
		ok, end := utils.ScanString(data, *i+1)
		if !ok {
			return fmt.Errorf("unterminated string")
		}
		*i = end + 1
		return nil
	case 't', 'f', 'n', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return parsePrimitive(data, i)
	default:
		return fmt.Errorf("unexpected character '%c'", data[*i])
	}
}

func parseObject(data []byte, i *int) error {
	if data[*i] != '{' {
		return fmt.Errorf("expected '{'")
	}
	*i++
	utils.SkipWS(i, data)

	if *i < len(data) && data[*i] == '}' {
		*i++
		return nil
	}

	for {
		utils.SkipWS(i, data)
		if *i >= len(data) || data[*i] != '"' {
			return fmt.Errorf("expected string key")
		}
		ok, end := utils.ScanString(data, *i+1)
		if !ok {
			return fmt.Errorf("unterminated key string")
		}
		*i = end + 1

		utils.SkipWS(i, data)
		if *i >= len(data) || data[*i] != ':' {
			return fmt.Errorf("expected ':' after key")
		}
		*i++
		if err := parseValue(data, i); err != nil {
			return err
		}

		utils.SkipWS(i, data)
		if *i >= len(data) {
			return fmt.Errorf("unexpected end of object")
		}
		if data[*i] == '}' {
			*i++
			return nil
		}
		if data[*i] != ',' {
			return fmt.Errorf("expected ',' between object items")
		}
		*i++
	}
}

func parseArray(data []byte, i *int) error {
	if data[*i] != '[' {
		return fmt.Errorf("expected '['")
	}
	*i++
	utils.SkipWS(i, data)

	if *i < len(data) && data[*i] == ']' {
		*i++
		return nil
	}

	for {
		if err := parseValue(data, i); err != nil {
			return err
		}
		utils.SkipWS(i, data)
		if *i >= len(data) {
			return fmt.Errorf("unexpected end of array")
		}
		if data[*i] == ']' {
			*i++
			return nil
		}
		if data[*i] != ',' {
			return fmt.Errorf("expected ',' between array items")
		}
		*i++
		utils.SkipWS(i, data)
	}
}

func parsePrimitive(data []byte, i *int) error {
	start := *i
	for *i < len(data) {
		c := data[*i]
		if c == ',' || c == '}' || c == ']' || utils.IsSpace(c) {
			break
		}
		*i++
	}
	if *i == start {
		return fmt.Errorf("invalid primitive")
	}
	return nil
}

// escapeString: ubah kutip & backslash agar valid JSON string
func escapeString(s []byte) []byte {
	out := make([]byte, 0, len(s))
	for _, c := range s {
		switch c {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		case '\r':
			out = append(out, '\\', 'r')
		case '\t':
			out = append(out, '\\', 't')
		default:
			out = append(out, c)
		}
	}
	return out
}

// isNumberString: cek sederhana apakah string berupa angka valid
func isNumberString(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, c := range s {
		if (c < '0' || c > '9') && c != '.' && c != '-' && c != '+' && c != 'e' && c != 'E' {
			return false
		}
		if (c == 'e' || c == 'E') && (i == 0 || i == len(s)-1) {
			return false
		}
	}
	return true
}
