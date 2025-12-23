package jepego

import (
	"fmt"
	"strings"
	"unicode"
)

type CaseStyle int

const (
	CaseKeep CaseStyle = iota
	CaseCamel
	CasePascal
	CaseSnake
	KeepItIs
)

func convertCase(s string, style CaseStyle) string {
	switch style {
	case CaseCamel:
		return toCamel(s)
	case CasePascal:
		return toPascal(s)
	case CaseSnake:
		return toSnake(s)
	default:
		return s
	}
}

func toCamel(s string) string {
	// ubah dulu ke PascalCase untuk normalisasi
	pascal := toPascal(s)
	if len(pascal) == 0 {
		return ""
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func toPascal(s string) string {
	if s == "" {
		return s
	}

	// normalisasi: ubah snake_case jadi spasi
	s = strings.ReplaceAll(s, "_", " ")

	var out []rune
	capNext := true

	for _, r := range s {
		if r == ' ' {
			capNext = true
			continue
		}
		if capNext {
			out = append(out, unicode.ToUpper(r))
			capNext = false
		} else {
			out = append(out, unicode.ToLower(r))
		}
	}

	return string(out)
}

func toSnake(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	var out []rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsUpper(r) {
			// jika bukan huruf pertama dan sebelumnya bukan underscore
			if i > 0 && runes[i-1] != '_' {
				// tambahkan underscore jika berikutnya bukan uppercase
				// atau kalau berikutnya ada lowercase (menandakan akhir blok uppercase)
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					out = append(out, '_')
				} else if !unicode.IsUpper(runes[i-1]) {
					out = append(out, '_')
				}
			}
			out = append(out, unicode.ToLower(r))
		} else {
			out = append(out, r)
		}
	}

	return string(out)
}

func transformValue(data []byte, i *int, out *[]byte, style CaseStyle) error {
	skipWS(i, data)
	if *i >= len(data) {
		return fmt.Errorf("unexpected end")
	}

	switch data[*i] {
	case '{':
		return transformObject(data, i, out, style)
	case '[':
		return transformArray(data, i, out, style)
	case '"':
		// string primitive: copy raw string token
		val, consumed := extractValue(data[*i:])
		if consumed == 0 {
			return fmt.Errorf("invalid string at %d", *i)
		}
		*out = append(*out, val...)
		*i += consumed
		return nil
	default:
		// primitive number / true / false / null
		val, consumed := extractValue(data[*i:])
		if consumed == 0 {
			return fmt.Errorf("invalid primitive at %d", *i)
		}
		*out = append(*out, val...)
		*i += consumed
		return nil
	}
}

func transformObject(data []byte, i *int, out *[]byte, style CaseStyle) error {
	// Expect '{'
	if data[*i] != '{' {
		return fmt.Errorf("expected '{' at %d", *i)
	}
	// write opening brace
	*out = append(*out, '{')
	*i++
	skipWS(i, data)

	// empty object
	if *i < len(data) && data[*i] == '}' {
		*out = append(*out, '}')
		*i++
		return nil
	}

	first := true
	for *i < len(data) {
		skipWS(i, data)
		if *i >= len(data) {
			return fmt.Errorf("unterminated object")
		}
		// expect key string
		if data[*i] != '"' {
			return fmt.Errorf("expected string key at %d", *i)
		}

		if !first {
			*out = append(*out, ',')
		}
		first = false

		// parse raw key string (including quotes)
		ok, keyEnd := scanString(data, *i+1)
		if !ok {
			return fmt.Errorf("unterminated key string")
		}
		// key content between *i+1 .. keyEnd-1
		keyRaw := data[*i+1 : keyEnd]
		// move index to position after closing quote
		*i = keyEnd + 1

		// unescape key to get logical characters (best-effort)
		keyUnescaped := unescapeKey(keyRaw)

		// convert case
		keyConv := convertCase(string(keyUnescaped), style)

		// write quoted, escaped converted key
		*out = append(*out, '"')
		*out = append(*out, escapeString([]byte(keyConv))...)
		*out = append(*out, '"')

		// skip whitespace and expect colon
		skipWS(i, data)
		if *i >= len(data) || data[*i] != ':' {
			return fmt.Errorf("expected ':' after key at %d", *i)
		}
		*out = append(*out, ':')
		*i++ // skip ':'

		// write value (recursively)
		if err := transformValue(data, i, out, style); err != nil {
			return err
		}

		skipWS(i, data)
		// if next is ',', consume and continue (loop will write comma earlier)
		if *i < len(data) && data[*i] == ',' {
			*i++ // consume comma in source; we've already appended comma before next pair
			// continue loop
			continue
		}
		// if next is '}', close object
		if *i < len(data) && data[*i] == '}' {
			*out = append(*out, '}')
			*i++
			return nil
		}
		// else error
		return fmt.Errorf("expected ',' or '}' after object entry at %d", *i)
	}
	return fmt.Errorf("unterminated object")
}

func transformArray(data []byte, i *int, out *[]byte, style CaseStyle) error {
	// expect '['
	if data[*i] != '[' {
		return fmt.Errorf("expected '[' at %d", *i)
	}
	*out = append(*out, '[')
	*i++
	skipWS(i, data)

	// empty array
	if *i < len(data) && data[*i] == ']' {
		*out = append(*out, ']')
		*i++
		return nil
	}

	first := true
	for *i < len(data) {
		if !first {
			*out = append(*out, ',')
		}
		first = false

		// parse next value
		if err := transformValue(data, i, out, style); err != nil {
			return err
		}

		skipWS(i, data)
		if *i < len(data) && data[*i] == ',' {
			*i++ // consume comma and continue
			continue
		}
		if *i < len(data) && data[*i] == ']' {
			*out = append(*out, ']')
			*i++
			return nil
		}
		return fmt.Errorf("expected ',' or ']' in array at %d", *i)
	}
	return fmt.Errorf("unterminated array")
}

// unescapeKey: best-effort unescape string content (logical characters).
// input: bytes between quotes (no surrounding quotes).
// handles: \", \\, \/, \b, \f, \n, \r, \t and copies \uXXXX as-is.
func unescapeKey(s []byte) []byte {
	out := make([]byte, 0, len(s))
	i := 0
	for i < len(s) {
		c := s[i]
		if c != '\\' {
			out = append(out, c)
			i++
			continue
		}
		// backslash sequence
		i++
		if i >= len(s) {
			// trailing backslash — copy it
			out = append(out, '\\')
			break
		}
		switch s[i] {
		case '"':
			out = append(out, '"')
			i++
		case '\\':
			out = append(out, '\\')
			i++
		case '/':
			out = append(out, '/')
			i++
		case 'b':
			out = append(out, '\b')
			i++
		case 'f':
			out = append(out, '\f')
			i++
		case 'n':
			out = append(out, '\n')
			i++
		case 'r':
			out = append(out, '\r')
			i++
		case 't':
			out = append(out, '\t')
			i++
		case 'u':
			// copy as-is: \uXXXX (6 bytes if valid) — we don't decode here
			// prefer copying the entire sequence to avoid mis-decoding
			if i+4 < len(s) {
				out = append(out, '\\', 'u')
				out = append(out, s[i+1:i+5]...)
				i += 5
			} else {
				// malformed — copy leftover
				out = append(out, '\\', 'u')
				i++
			}
		default:
			// unknown escape, copy both
			out = append(out, '\\', s[i])
			i++
		}
	}
	return out
}
