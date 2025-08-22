package jsonParser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	Raw string
}

// menghapus semua "// comment" dari JSON string jika ada / saat test
func RemoveComments(jsonStr string) string {
	lines := strings.Split(jsonStr, "\n")
	var b strings.Builder
	for _, line := range lines {
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx] // hapus bagian setelah //
		}
		line = strings.TrimSpace(line)
		if line != "" {
			b.WriteString(line)
		}
	}
	return b.String()
}

// Get mengambil nilai JSON berdasarkan path
func Get(jsonStr, path string) Result {
	// Bersihin dulu comment jika ada
	jsonStr = RemoveComments(jsonStr)

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		fmt.Printf("error: %s\n", err)
		return Result{Raw: "null"}
	}

	// Modifier @ / !
	if strings.HasPrefix(path, "@") || strings.HasPrefix(path, "!") {
		key := path[1:]
		if m, ok := data.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				raw, _ := json.Marshal(val)
				return Result{Raw: string(raw)}
			}
		}
		return Result{Raw: "null"}
	}

	// Subselector object {field1,field2}
	if strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}") {
		fields := strings.Split(path[1:len(path)-1], ",")
		if m, ok := data.(map[string]interface{}); ok {
			sub := map[string]interface{}{}
			for _, f := range fields {
				f = strings.TrimSpace(f)
				if val, exists := m[f]; exists {
					sub[f] = val
				}
			}

			// Build JSON string manually untuk jaga urutan
			parts := []string{}
			for _, f := range fields {
				f = strings.TrimSpace(f)
				if val, exists := sub[f]; exists {
					b, _ := json.Marshal(val)
					parts = append(parts, fmt.Sprintf(`"%s":%s`, f, string(b)))
				}
			}
			return Result{Raw: "{" + strings.Join(parts, ",") + "}"}
		}
		return Result{Raw: "null"}
	}

	// Subselector array [0,2]
	if strings.HasPrefix(path, "[") && strings.HasSuffix(path, "]") {
		idxStr := path[1 : len(path)-1]
		idxParts := strings.Split(idxStr, ",")
		if arr, ok := data.([]interface{}); ok {
			subArr := []interface{}{}
			for _, s := range idxParts {
				s = strings.TrimSpace(s)
				idx, err := strconv.Atoi(s)
				if err == nil && idx >= 0 && idx < len(arr) {
					subArr = append(subArr, arr[idx])
				}
			}
			raw, _ := json.Marshal(subArr)
			return Result{Raw: string(raw)}
		}
		return Result{Raw: "null"}
	}

	// Nested path dan array index
	parts := splitPath(path)
	current := data

	for _, part := range parts {
		switch node := current.(type) {
		case map[string]interface{}:
			if val, ok := node[part]; ok {
				current = val
			} else {
				return Result{Raw: "null"}
			}
		case []interface{}:
			// array index [i]
			if idx, err := strconv.Atoi(part); err == nil && idx >= 0 && idx < len(node) {
				current = node[idx]
			} else {
				return Result{Raw: "null"}
			}
		default:
			return Result{Raw: "null"}
		}
	}

	rawBytes, err := json.Marshal(current)
	if err != nil {
		return Result{Raw: "null"}
	}
	return Result{Raw: string(rawBytes)}
}

// splitPath memisahkan nested path dan array index
func splitPath(path string) []string {
	// Contoh: details.city -> ["details","city"]
	//          items[0].name -> ["items","0","name"]
	re := regexp.MustCompile(`\.|\[(\d+)\]`)
	matches := re.FindAllStringSubmatchIndex(path, -1)

	parts := []string{}
	last := 0
	for _, match := range matches {
		if match[0] > last {
			parts = append(parts, path[last:match[0]])
		}
		if len(match) > 2 && match[2] != -1 {
			parts = append(parts, path[match[2]:match[3]])
		}
		last = match[1]
	}
	if last < len(path) {
		parts = append(parts, path[last:])
	}
	return parts
}
