package jepego

// extractValue: ambil 1 JSON value utuh dari s[0:], return (slice, bytesConsumed)
func extractValue(s []byte) ([]byte, int) {
	s = trimSpaceBytes(s)
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
		for j > 0 && isSpace(s[j-1]) {
			j--
		}
		return s[:j], i
	}
	return nil, 0
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
