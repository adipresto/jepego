package core

import {
	"github.com/adipresto/jepego/model/constant"
	"github.com/adipresto/jepego/model/utils/byte"
	"github.com/adipresto/jepego/model/utils/extract"
	"github.com/adipresto/jepego/model/utils/path"
	"github.com/adipresto/jepego/model/utils/typecheck"
}

// Get mengambil value berdasarkan path (contoh: "a.b[3].c").
func Get(json []byte, path string) Result {
	json = removeCommentsBytes(json)
	if len(path) == 0 {
		return Result{Key: "", OK: false}
	}

	val, ok := getNestedValue(json, splitPathBytes([]byte(path)))
	if !ok {
		return Result{Key: path, OK: false}
	}

	var resVal []byte
	if val[0] == '"' {
		resVal = val[1 : len(val)-1]
	} else {
		resVal = val
	}

	return Result{
		Key:      path,
		Data:     resVal,
		DataType: detectType(val),
		OK:       true,
	}
}

// GetMany menerima ekspresi dalam bentuk []string{"a","b","c.d","arr[5].x"}.
// Mengembalikan hasil sesuai urutan field yang diminta.
func GetMany(json []byte, exprs []string) map[string]Result {
	json = removeCommentsBytes(json)
	out := make(map[string]Result, len(exprs))
	for _, p := range exprs {
		if len(p) == 0 {
			out[p] = Result{Key: "", OK: false}
			continue
		}

		val, ok := getNestedValue(json, splitPathBytes([]byte(p)))
		if !ok {
			out[p] = Result{Key: p, OK: false}
			continue
		}

		var retVal []byte
		if val[0] == '"' {
			retVal = val[1 : len(val)-1]
		} else {
			retVal = val
		}

		out[p] = Result{
			Key:      p,
			Data:     retVal,
			DataType: detectType(val),
			OK:       true,
		}
	}
	return out
}

func GetAll(json []byte, path string) []Result {
	json = removeCommentsBytes(json)
	if len(path) == 0 {
		return nil
	}

	toks := splitPathBytes([]byte(path))
	return getNestedValues(json, toks, path)
}

// GetManyAll menerima banyak ekspresi path (contoh: "a", "b[0].c", "arr[].x").
// Berbeda dengan GetMany yang hanya ambil satu value per key,
// GetManyAll akan expand semua hasil jika ada wildcard [] di path.
// Hasil dikembalikan dalam bentuk map[string][]Result agar tiap key bisa punya banyak item.
func GetManyAll(json []byte, exprs []string) map[string][]Result {
	json = removeCommentsBytes(json)
	out := make(map[string][]Result, len(exprs))

	for _, p := range exprs {
		if len(p) == 0 {
			out[p] = []Result{{Key: "", OK: false}}
			continue
		}

		toks := splitPathBytes([]byte(p))

		// kalau ada wildcard, pakai getNestedValues
		hasWildcard := false
		for _, t := range toks {
			if t.isWildcard {
				hasWildcard = true
				break
			}
		}

		if hasWildcard {
			results := getNestedValues(json, toks, p)
			if len(results) == 0 {
				out[p] = []Result{{Key: p, OK: false}}
			} else {
				out[p] = results
			}
			continue
		}

		// kalau tanpa wildcard, fallback ke getNestedValue
		val, ok := getNestedValue(json, toks)
		if !ok {
			out[p] = []Result{{Key: p, OK: false}}
			continue
		}
		out[p] = []Result{{
			Key:      p,
			Data:     unwrap(val),
			DataType: detectType(val),
			OK:       true,
		}}
	}

	return out
}
