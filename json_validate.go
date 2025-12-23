package jepego

import "fmt"

// parseValue: parse 1 JSON value dari data[i:], update i ke posisi setelah value
func parseValue(data []byte, i *int) error {
	skipWS(i, data)
	if *i >= len(data) {
		return fmt.Errorf("unexpected end of data")
	}

	switch data[*i] {
	case '{':
		return parseObject(data, i)
	case '[':
		return parseArray(data, i)
	case '"':
		ok, end := scanString(data, *i+1)
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
	skipWS(i, data)

	if *i < len(data) && data[*i] == '}' {
		*i++
		return nil
	}

	for {
		skipWS(i, data)
		if *i >= len(data) || data[*i] != '"' {
			return fmt.Errorf("expected string key")
		}
		ok, end := scanString(data, *i+1)
		if !ok {
			return fmt.Errorf("unterminated key string")
		}
		*i = end + 1

		skipWS(i, data)
		if *i >= len(data) || data[*i] != ':' {
			return fmt.Errorf("expected ':' after key")
		}
		*i++
		if err := parseValue(data, i); err != nil {
			return err
		}

		skipWS(i, data)
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
	skipWS(i, data)

	if *i < len(data) && data[*i] == ']' {
		*i++
		return nil
	}

	for {
		if err := parseValue(data, i); err != nil {
			return err
		}
		skipWS(i, data)
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
		skipWS(i, data)
	}
}

func parsePrimitive(data []byte, i *int) error {
	start := *i
	for *i < len(data) {
		c := data[*i]
		if c == ',' || c == '}' || c == ']' || isSpace(c) {
			break
		}
		*i++
	}
	if *i == start {
		return fmt.Errorf("invalid primitive")
	}
	return nil
}
