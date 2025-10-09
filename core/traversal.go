package core

import (
	"github.com/adipresto/jepego/model"
	"github.com/adipresto/jepego/utils"
)

// getNestedValue: berjalan mengikuti tokens (field / [index])
func GetNestedValue(json []byte, parts []utils.PathToken) ([]byte, bool) {
	cur := utils.TrimSpaceBytes(json)
	for _, p := range parts {
		if len(cur) == 0 {
			return nil, false
		}
		switch cur[0] {
		case '{':
			if p.IsIdx {
				return nil, false
			}
			v, ok := utils.GetTopLevelKey(cur, p.Key)
			if !ok {
				return nil, false
			}
			cur = utils.TrimSpaceBytes(v)
		case '[':
			if !p.IsIdx {
				return nil, false
			}
			v, ok := utils.GetArrayIndex(cur, p.Index)
			if !ok {
				return nil, false
			}
			cur = utils.TrimSpaceBytes(v)
		default:
			return nil, false
		}
	}
	return cur, true
}

// getNestedValues: berjalan mengikuti tokens (field / [index / []]) dan mengembalikan semua hasil yang ditemukan.
func GetNestedValues(json []byte, parts []utils.PathToken, fullPath string) []model.Result {
	cur := utils.TrimSpaceBytes(json)
	if len(parts) == 0 {
		return []model.Result{{
			Key:      fullPath,
			Data:     utils.Unwrap(cur),
			DataType: utils.DetectType(cur),
			OK:       true,
		}}
	}

	p := parts[0]
	rest := parts[1:]

	switch cur[0] {
	case '{':
		if p.IsIdx {
			return nil
		}
		v, ok := utils.GetTopLevelKey(cur, p.Key)
		if !ok {
			return nil
		}
		return GetNestedValues(v, rest, fullPath)

	case '[':
		if p.IsIdx {
			// index spesifik
			v, ok := utils.GetArrayIndex(cur, p.Index)
			if !ok {
				return nil
			}
			return GetNestedValues(v, rest, fullPath)
		}

		// wildcard []
		if len(p.Key) == 0 && !p.IsIdx {
			var results []model.Result
			i := 1
			n := len(cur)
			for i < n {
				utils.SkipWS(&i, cur)
				if i >= n || cur[i] == ']' {
					break
				}
				val, consumed := utils.ExtractValue(cur[i:])
				if consumed == 0 {
					break
				}
				sub := GetNestedValues(val, rest, fullPath)
				results = append(results, sub...)
				i += consumed
				utils.SkipWS(&i, cur)
				if i < n && cur[i] == ',' {
					i++
				}
			}
			return results
		}
	}
	return nil
}
