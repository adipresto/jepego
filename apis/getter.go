package apis

import (
	"github.com/adipresto/jepego/core"
	"github.com/adipresto/jepego/model"
	"github.com/adipresto/jepego/utils"
)

// Get mengambil value berdasarkan path (contoh: "a.b[3].c").
func Get(json []byte, path string) model.Result {
	json = utils.RemoveCommentsBytes(json)
	if len(path) == 0 {
		return model.Result{Key: "", OK: false}
	}

	val, ok := core.GetNestedValue(json, utils.SplitPathBytes([]byte(path)))
	if !ok {
		return model.Result{Key: path, OK: false}
	}

	var resVal []byte
	if val[0] == '"' {
		resVal = val[1 : len(val)-1]
	} else {
		resVal = val
	}

	return model.Result{
		Key:      path,
		Data:     resVal,
		DataType: utils.DetectType(val),
		OK:       true,
	}
}

// GetMany menerima ekspresi dalam bentuk []string{"a","b","c.d","arr[5].x"}.
// Mengembalikan hasil sesuai urutan field yang diminta.
func GetMany(json []byte, exprs []string) map[string]model.Result {
	json = utils.RemoveCommentsBytes(json)
	out := make(map[string]model.Result, len(exprs))
	for _, p := range exprs {
		if len(p) == 0 {
			out[p] = model.Result{Key: "", OK: false}
			continue
		}

		val, ok := core.GetNestedValue(json, utils.SplitPathBytes([]byte(p)))
		if !ok {
			out[p] = model.Result{Key: p, OK: false}
			continue
		}

		var retVal []byte
		if val[0] == '"' {
			retVal = val[1 : len(val)-1]
		} else {
			retVal = val
		}

		out[p] = model.Result{
			Key:      p,
			Data:     retVal,
			DataType: utils.DetectType(val),
			OK:       true,
		}
	}
	return out
}

func GetAll(json []byte, path string) []model.Result {
	json = utils.RemoveCommentsBytes(json)
	if len(path) == 0 {
		return nil
	}

	toks := utils.SplitPathBytes([]byte(path))
	return core.GetNestedValues(json, toks, path)
}

// GetManyAll menerima banyak ekspresi path (contoh: "a", "b[0].c", "arr[].x").
// Berbeda dengan GetMany yang hanya ambil satu value per key,
// GetManyAll akan expand semua hasil jika ada wildcard [] di path.
// Hasil dikembalikan dalam bentuk map[string][]Result agar tiap key bisa punya banyak item.
func GetManyAll(json []byte, exprs []string) map[string][]model.Result {
	json = utils.RemoveCommentsBytes(json)
	out := make(map[string][]model.Result, len(exprs))

	for _, p := range exprs {
		if len(p) == 0 {
			out[p] = []model.Result{{Key: "", OK: false}}
			continue
		}

		toks := utils.SplitPathBytes([]byte(p))

		// kalau ada wildcard, pakai getNestedValues
		hasWildcard := false
		for _, t := range toks {
			if t.IsWildcard {
				hasWildcard = true
				break
			}
		}

		if hasWildcard {
			results := core.GetNestedValues(json, toks, p)
			if len(results) == 0 {
				out[p] = []model.Result{{Key: p, OK: false}}
			} else {
				out[p] = results
			}
			continue
		}

		// kalau tanpa wildcard, fallback ke getNestedValue
		val, ok := core.GetNestedValue(json, toks)
		if !ok {
			out[p] = []model.Result{{Key: p, OK: false}}
			continue
		}
		out[p] = []model.Result{{
			Key:      p,
			Data:     utils.Unwrap(val),
			DataType: utils.DetectType(val),
			OK:       true,
		}}
	}

	return out
}
