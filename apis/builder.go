package apis

func BuildJSON(m map[string]string) []byte {
	if len(m) == 0 {
		return []byte("{}")
	}

	buf := make([]byte, 0, 256)
	buf = append(buf, '{')

	first := true
	for k, v := range m {
		if !first {
			buf = append(buf, ',')
		}
		first = false

		buf = append(buf, '"')
		buf = append(buf, escapeString([]byte(k))...)
		buf = append(buf, '"', ':')

		// coba deteksi tipe sederhana (angka/bool/null)
		if isNumberString(v) || v == "true" || v == "false" || v == "null" {
			buf = append(buf, v...)
		} else {
			buf = append(buf, '"')
			buf = append(buf, escapeString([]byte(v))...)
			buf = append(buf, '"')
		}
	}
	buf = append(buf, '}')
	return buf
}
