package jepego

// Pecah path seperti "a.b[3].c"
// Tanpa clone → setiap key adalah subslice langsung dari path.
// Lebih hemat memori, tapi hasil valid hanya selama `path` masih ada.
func splitPathBytes(path []byte) []pathToken {
	var toks []pathToken
	start := 0 // pointer awal token (key)

	for i := 0; i < len(path); i++ {
		c := path[i]
		switch c {
		case '.':
			// "a.b" → token pertama selesai di '.'
			if i > start {
				toks = append(toks, pathToken{key: path[start:i]})
			}
			start = i + 1 // geser start ke char setelah '.'

		case '[':
			// contoh: "arr[3]" → key = "arr"
			if i > start {
				toks = append(toks, pathToken{key: path[start:i]})
			}
			i++ // skip '['

			if i < len(path) && path[i] == ']' {
				// [] = wildcard
				toks = append(toks, pathToken{isWildcard: true})
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
			toks = append(toks, pathToken{isIdx: true, index: idx})
			start = i + 1
		}
	}

	// token terakhir (kalau masih ada sisa)
	if start < len(path) {
		toks = append(toks, pathToken{key: path[start:]})
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
