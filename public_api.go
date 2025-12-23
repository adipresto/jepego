package jepego

import (
	"fmt"
	"strings"
)

// DeleteField menghapus sebuah field dari JSON berdasarkan path.
//
// Parameter:
//   - json: data JSON mentah dalam bentuk byte slice.
//   - field: path field yang akan dihapus (misalnya "user.profile.name").
//
// Perilaku:
//   - JSON akan dibersihkan dari komentar sebelum diproses.
//   - Path field akan di-split menjadi token-token.
//   - Jika field tidak ditemukan, JSON original akan dikembalikan tanpa perubahan.
//   - Jika field ditemukan, field tersebut akan dihapus secara rekursif.
//
// Return:
//   - JSON hasil penghapusan field, atau JSON original jika field tidak ditemukan.
func DeleteField(json []byte, field string) []byte {
	json = removeCommentsBytes(json)
	toks := splitPathBytes([]byte(field))
	if len(toks) == 0 {
		return json
	}

	out, ok := deleteNested(json, toks)
	if !ok {
		return json // field tidak ditemukan → return original
	}
	return out
}

// UpsertRemapJSON membangun JSON baru berdasarkan mapping field dari JSON sumber.
//
// Fungsi ini mengekstrak satu atau lebih field dari JSON sumber,
// kemudian memetakan hasilnya ke key baru pada JSON output.
//
// Parameter:
//   - jsonData: JSON sumber dalam bentuk byte slice.
//   - mapping: map dengan key berupa path JSON (dapat lebih dari satu, dipisah koma)
//     dan value berupa nama field baru di output JSON.
//     Contoh:
//     {
//     "user.first_name,user.last_name": "full_name",
//     "age": "umur"
//     }
//   - sep: separator yang digunakan jika lebih dari satu field digabung.
//   - caseStyle: gaya penamaan field output (misalnya snake_case, camelCase, dll).
//
// Perilaku:
//   - Jika beberapa path didefinisikan dalam satu mapping key,
//     semua nilai yang ditemukan akan digabung menggunakan separator.
//   - Jika tidak ada field yang ditemukan, field output tetap dibuat
//     dengan nilai string kosong.
//   - Tipe data (string, number, object, array, dll) dipertahankan
//     jika hanya satu nilai yang ditemukan.
//
// Return:
//   - JSON baru hasil transformasi dan remapping field.
func UpsertRemapJSON(jsonData []byte, mapping map[string]string, sep string, caseStyle CaseStyle) []byte {
	out := make(map[string]Result, len(mapping))

	for pathExpr, newKey := range mapping {
		paths := strings.Split(pathExpr, ",")
		for i := range paths {
			paths[i] = strings.TrimSpace(paths[i])
		}

		results := GetMany(jsonData, paths)
		values := make([]string, 0, len(paths))
		var dtype DataType = TypeString

		for _, p := range paths {
			if r, ok := results[p]; ok && r.OK {
				values = append(values, string(r.Data))
				dtype = r.DataType
			}
		}

		// konversi gaya nama field baru (jika diminta)
		convertedKey := convertCase(newKey, caseStyle)

		if len(values) == 0 {
			out[convertedKey] = Result{
				Key:      convertedKey,
				Data:     []byte(""),
				DataType: TypeString,
				OK:       true,
			}
			continue
		}

		if len(values) > 1 {
			out[convertedKey] = Result{
				Key:      convertedKey,
				Data:     []byte(strings.Join(values, sep)),
				DataType: TypeString,
				OK:       true,
			}
			continue
		}

		out[convertedKey] = Result{
			Key:      convertedKey,
			Data:     []byte(values[0]),
			DataType: dtype,
			OK:       true,
		}
	}

	return BuildJSONFromResults(out)
}

// BuildJSONFromResults membangun JSON object dari map Result.
//
// Fungsi ini digunakan sebagai tahap akhir untuk mengubah hasil ekstraksi
// dan transformasi data menjadi JSON valid.
//
// Parameter:
//   - m: map dengan key berupa nama field dan value berupa Result
//     yang berisi data mentah, tipe data, dan status OK.
//
// Perilaku:
//   - Hanya Result dengan OK == true yang akan dimasukkan ke output JSON.
//   - Field name yang mengandung "." hanya akan mengambil segmen terakhir.
//   - Tipe data akan ditulis sesuai DataType:
//   - String → di-escape dan dibungkus tanda petik
//   - Number, Bool, Null → ditulis langsung
//   - Object, Array → ditulis sebagai JSON mentah
//
// Return:
//   - JSON object dalam bentuk byte slice.
//   - Jika map kosong, akan mengembalikan "{}".
func BuildJSONFromResults(m map[string]Result) []byte {
	if len(m) == 0 {
		return []byte("{}")
	}

	buf := make([]byte, 0, 256)
	buf = append(buf, '{')

	first := true
	for k, r := range m {
		if !r.OK {
			continue
		}

		splitFields := strings.Split(k, ".")
		field := splitFields[len(splitFields)-1]

		if !first {
			buf = append(buf, ',')
		}
		first = false

		buf = append(buf, '"')
		buf = append(buf, escapeString([]byte(field))...)
		buf = append(buf, '"', ':')

		switch r.DataType {
		case TypeString:
			buf = append(buf, '"')
			buf = append(buf, escapeString(r.Data)...)
			buf = append(buf, '"')

		case TypeNumber, TypeBool, TypeNull:
			// angka, true/false, null → langsung tulis apa adanya
			buf = append(buf, r.Data...)

		case TypeObject, TypeArray:
			// data mentah object/array langsung masuk
			buf = append(buf, r.Data...)

		default:
			// fallback → string
			buf = append(buf, '"')
			buf = append(buf, escapeString(r.Data)...)
			buf = append(buf, '"')
		}
	}

	buf = append(buf, '}')
	return buf
}

// BuildJSON membangun JSON object sederhana dari map string ke string.
//
// Fungsi ini cocok untuk pembuatan JSON ringan tanpa struktur bersarang.
//
// Parameter:
//   - m: map[string]string yang berisi pasangan key-value JSON.
//
// Perilaku:
//   - Nilai akan dicoba dideteksi sebagai number, boolean, atau null.
//   - Jika tidak cocok, nilai akan diperlakukan sebagai string dan di-escape.
//   - Tidak mendukung nested object atau array.
//
// Return:
//   - JSON object dalam bentuk byte slice.
//   - Jika map kosong, akan mengembalikan "{}".
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

// ValidateJSONRobust memeriksa apakah data JSON valid.
// Mengembalikan error dengan info posisi byte kalau ada kesalahan.
func ValidateJSON(data []byte) error {
	data = removeCommentsBytes(data)
	data = trimSpaceBytes(data)

	if len(data) == 0 {
		return fmt.Errorf("empty JSON")
	}

	i := 0
	if err := parseValue(data, &i); err != nil {
		return fmt.Errorf("invalid JSON at byte %d: %v", i, err)
	}

	// pastikan tidak ada sisa byte setelah value valid
	skipWS(&i, data)
	if i != len(data) {
		return fmt.Errorf("invalid JSON: extra content after valid value at byte %d", i)
	}

	return nil
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

// Upsert: update atau insert value di path tertentu
func Upsert(json []byte, path string, value []byte, valueType DataType) ([]byte, error) {
	if len(path) == 0 {
		return json, fmt.Errorf("empty path")
	}

	tokens := splitPathBytes([]byte(path))
	return upsertNested(json, tokens, value, valueType)
}

// TransformCaseJSON: streaming transform case tanpa menggunakan encoding/json
func TransformCaseJSON(jsonData []byte, style CaseStyle) ([]byte, error) {
	if len(jsonData) == 0 {
		return []byte("{}"), nil
	}
	i := 0
	out := make([]byte, 0, len(jsonData))
	if err := transformValue(jsonData, &i, &out, style); err != nil {
		return nil, err
	}
	// skip trailing whitespace
	skipWS(&i, jsonData)
	if i != len(jsonData) {
		return nil, fmt.Errorf("extra content after JSON at byte %d", i)
	}
	return out, nil
}
