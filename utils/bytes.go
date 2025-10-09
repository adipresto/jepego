package utils

// trimSpaceBytes menghapus spasi di awal dan akhir slice byte.
func TrimSpaceBytes(b []byte) []byte {
	n := len(b)
	start := 0

	// Cari awal bukan spasi
	for start < n && IsSpace(b[start]) {
		start++
	}
	if start == n {
		return b[:0] // semua spasi
	}

	end := n - 1
	// Cari akhir bukan spasi
	for end >= start && IsSpace(b[end]) {
		end--
	}

	// Jika tidak ada prefix spasi, kembalikan slice asli yang sudah dipotong suffix
	if start == 0 {
		return b[:end+1]
	}

	// Ada prefix spasi, geser data ke awal slice
	copy(b[0:], b[start:end+1])
	return b[:end-start+1]
}

func IsSpace(c byte) bool {
	switch c {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

func CloneBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

func BytesEqual(a, b []byte) bool {
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
func EscapeString(s []byte) []byte {
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

// Hapus komentar // tanpa alokasi baru sebesar len(data).
// Beda dengan versi lama, ini nulis ulang ke buffer yang sama (in-place).
// Return slice hasil trimming, masih pointing ke buffer asli.
func RemoveCommentsBytes(data []byte) []byte {
	n := len(data)     // panjang total input
	write := 0         // posisi tulis (writer index)
	inComment := false // flag: sedang dalam komentar atau tidak

	for i := 0; i < n; i++ { // loop semua byte
		// deteksi awal komentar //
		if !inComment && i+1 < n && data[i] == '/' && data[i+1] == '/' {
			inComment = true // masuk mode komentar
			i++              // skip char kedua dari "//"
			continue
		}
		// kalau lagi dalam komentar
		if inComment {
			// komentar selesai hanya kalau ketemu newline
			if data[i] == '\n' {
				inComment = false
				data[write] = data[i] // tulis newline ke posisi write
				write++
			}
			continue // skip semua char lain di komentar
		}
		// kalau bukan komentar → salin byte ke posisi write
		data[write] = data[i]
		write++
	}
	// hasil akhir = slice yang dipangkas + di-trim spasi
	return TrimSpaceBytes(data[:write])
}
