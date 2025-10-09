package utils

type PathToken struct {
	key        []byte
	index      int
	isIdx      bool
	isWildcard bool
}

// Pecah path seperti "a.b[3].c"
// Tanpa clone → setiap key adalah subslice langsung dari path.
// Lebih hemat memori, tapi hasil valid hanya selama `path` masih ada.
func SplitPathBytes(path []byte) []PathToken {
	var toks []PathToken
	start := 0 // pointer awal token (key)

	for i := 0; i < len(path); i++ {
		c := path[i]
		switch c {
		case '.':
			// "a.b" → token pertama selesai di '.'
			if i > start {
				toks = append(toks, PathToken{key: path[start:i]})
			}
			start = i + 1 // geser start ke char setelah '.'

		case '[':
			// contoh: "arr[3]" → key = "arr"
			if i > start {
				toks = append(toks, PathToken{key: path[start:i]})
			}
			i++ // skip '['

			if i < len(path) && path[i] == ']' {
				// [] = wildcard
				toks = append(toks, PathToken{isWildcard: true})
				start = i + 1
				continue
			}

			// parse angka index di dalam bracket
			idx := 0
			for i < len(path) && path[i] != ']' {
				idx = idx*10 + int(path[i]-'0') // konversi char → int
				i++
			}
			// selesai bracket → simpan token index
			toks = append(toks, PathToken{isIdx: true, index: idx})
			start = i + 1
		}
	}

	// token terakhir (kalau masih ada sisa)
	if start < len(path) {
		toks = append(toks, PathToken{key: path[start:]})
	}
	return toks
}

// parseMultiExpr: "{a,b,c.d,arr[5].x}" -> [][]byte{"a","b","c.d","arr[5].x"}
func ParseMultiExpr(expr []byte) [][]byte {
	b := TrimSpaceBytes(expr)
	if len(b) >= 2 && b[0] == '{' && b[len(b)-1] == '}' {
		b = b[1 : len(b)-1]
	}
	return SplitCSVTopLevel(b)
}

// splitCSVTopLevel: pecah "a,b,c" jadi potongan, asumsi nama field tidak mengandung koma/kutip.
func SplitCSVTopLevel(b []byte) [][]byte {
	var parts [][]byte
	start := 0
	for i := 0; i <= len(b); i++ {
		if i == len(b) || b[i] == ',' {
			seg := TrimSpaceBytes(b[start:i])
			parts = append(parts, CloneBytes(seg))
			start = i + 1
		}
	}
	return parts
}
