package jepego

// removeCommentsBytes: hapus komentar // hanya jika di luar string literal.
// Tetap bekerja in-place tanpa alokasi besar baru.
func removeCommentsBytes(data []byte) []byte {
	n := len(data)
	write := 0

	inString := false  // sedang dalam tanda kutip "
	escaped := false   // apakah char sebelumnya adalah '\'
	inComment := false // sedang dalam komentar //

	for i := 0; i < n; i++ {
		c := data[i]

		// --- Saat sedang di dalam komentar ---
		if inComment {
			if c == '\n' { // komentar berakhir di newline
				inComment = false
				data[write] = c
				write++
			}
			continue
		}

		// --- Saat sedang di dalam string literal ---
		if inString {
			data[write] = c
			write++

			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inString = false
			}
			continue
		}

		// --- Saat di luar string ---
		if c == '"' {
			inString = true
			data[write] = c
			write++
			continue
		}

		// Deteksi awal komentar // hanya di luar string
		if c == '/' && i+1 < n && data[i+1] == '/' {
			inComment = true
			i++ // skip '/'
			continue
		}

		// salin byte normal
		data[write] = c
		write++
	}

	// hasil akhir, di-trim whitespace
	return trimSpaceBytes(data[:write])
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
