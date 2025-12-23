package jepego

func deleteNested(json []byte, parts []pathToken) ([]byte, bool) {
	cur := trimSpaceBytes(json)

	// base case: hapus field di object ini
	if len(parts) == 1 {
		if cur[0] != '{' {
			return json, false
		}
		return deleteKeyFromObject(cur, parts[0].key)
	}

	p := parts[0]
	rest := parts[1:]

	if cur[0] != '{' || p.isIdx {
		return json, false
	}

	// ambil value child
	val, ok := getTopLevelKey(cur, p.key)
	if !ok {
		return json, false
	}

	// recursive delete
	newChild, changed := deleteNested(val, rest)
	if !changed {
		return json, false
	}

	// rebuild object dengan child baru
	return replaceKeyInObject(cur, p.key, newChild), true
}

func deleteKeyFromObject(obj []byte, key []byte) ([]byte, bool) {
	i := 1
	n := len(obj)

	buf := make([]byte, 0, len(obj))
	buf = append(buf, '{')

	first := true
	deleted := false

	for i < n {
		skipWS(&i, obj)
		if i >= n || obj[i] == '}' {
			break
		}

		// parse key
		if obj[i] != '"' {
			return obj, false
		}
		keyStart := i + 1
		ok, keyEnd := scanString(obj, keyStart)
		if !ok {
			return obj, false
		}
		i = keyEnd + 1

		skipWS(&i, obj)
		if obj[i] != ':' {
			return obj, false
		}
		i++
		skipWS(&i, obj)

		val, consumed := extractValue(obj[i:])
		if consumed == 0 {
			return obj, false
		}

		isTarget := bytesEqual(obj[keyStart:keyEnd], key)

		if !isTarget {
			if !first {
				buf = append(buf, ',')
			}
			first = false

			buf = append(buf, '"')
			buf = append(buf, obj[keyStart:keyEnd]...)
			buf = append(buf, '"', ':')

			buf = append(buf, val...)
		} else {
			deleted = true
		}

		i += consumed
		skipWS(&i, obj)
		if i < n && obj[i] == ',' {
			i++
		}
	}

	buf = append(buf, '}')

	if !deleted {
		return obj, false
	}
	return buf, true
}

func replaceKeyInObject(obj []byte, key []byte, newVal []byte) []byte {
	i := 1
	n := len(obj)

	buf := make([]byte, 0, len(obj))
	buf = append(buf, '{')

	first := true

	for i < n {
		skipWS(&i, obj)
		if i >= n || obj[i] == '}' {
			break
		}

		keyStart := i + 1
		_, keyEnd := scanString(obj, keyStart)
		i = keyEnd + 1

		skipWS(&i, obj)
		i++ // skip :
		skipWS(&i, obj)

		val, consumed := extractValue(obj[i:])

		if !first {
			buf = append(buf, ',')
		}
		first = false

		buf = append(buf, '"')
		buf = append(buf, obj[keyStart:keyEnd]...)
		buf = append(buf, '"', ':')

		if bytesEqual(obj[keyStart:keyEnd], key) {
			buf = append(buf, newVal...)
		} else {
			buf = append(buf, val...)
		}

		i += consumed
		skipWS(&i, obj)
		if i < n && obj[i] == ',' {
			i++
		}
	}

	buf = append(buf, '}')
	return buf
}
