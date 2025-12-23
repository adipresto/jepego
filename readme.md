# JEPEGO

---
# JEPEGO (byte-level JSON traversal for Go)

> **Fast, zero-allocation JSON traversal, mutation, and transformation — without decoding.**

`jepego` adalah library Go untuk **membaca, memodifikasi, dan mentransform JSON langsung dari `[]byte`**, tanpa:

* `encoding/json`
* struct binding
* reflection
* AST

Dirancang untuk **performance-critical services**, API gateway, dan middleware.

## Installation

```bash
go get github.com/yourname/jsonparser
```

---

## Quick Start

### Get value by path

```go
v, ok := jsonparser.Get(data, "user.profile.name")
if ok {
    fmt.Println(string(v))
}
```

---

### Get all values (wildcard)

```go
res := jsonparser.GetAll(data, "users[].id")
for _, r := range res {
    fmt.Println(string(r.Data))
}
```

---

### Delete field

```go
out := jsonparser.DeleteField(data, "user.password")
```

---

### Upsert field

```go
out := jsonparser.Upsert(data, "user.role", []byte(`"admin"`))
```

---

### Transform key case

```go
out, _ := jsonparser.TransformCaseJSON(data, CamelCase)
```

---

# Main Features

Library ini bekerja langsung di level **`[]byte` JSON** tanpa `encoding/json`.

## Comment Stripping (`//`) – JSON Superset

Library ini mendukung JSON dengan komentar `//`, mirip JSONC.

### Contoh

```json
{
  // user identity
  "user": {
    "name": "Setyo", // display name
    "age": 30
  }
}
```

Semua API publik **secara otomatis**:

* Menghapus komentar
* Tetap aman terhadap string literal (`"http://..."` tidak rusak)

### Cara Kerja

* Single pass
* In-place rewrite (tanpa alokasi buffer besar)
* Komentar hanya dihapus **di luar string**

### Implementasi inti:

```go
removeCommentsBytes
```

---

## Wildcard Traversal (`[]`) – Expand Array Values

Path mendukung wildcard `[]` untuk mengambil **semua elemen array**.

### Contoh JSON

```json
{
  "users": [
    { "id": 1, "name": "Adi" },
    { "id": 2, "name": "Prasetyo" }
  ]
}
```

### Ambil semua nama user

```go
res := jsonparser.GetAll(data, "users[].name")
```

### Hasil

```go
[]Result{
  { Key: "users[].name", Data: []byte("Adi"), OK: true },
  { Key: "users[].name", Data: []byte("Prasetyo"), OK: true },
}
```

### Karakteristik

* Tidak perlu loop manual
* Tidak decode array ke struct
* Depth-first traversal
* Aman untuk nested wildcard

### Implementasi inti:

```go
getNestedValues
```

---

## Streaming Case Transform (Zero Decode)

Transformasi gaya penamaan key **tanpa parsing ke struct atau map**.

### Contoh

```json
{
  "user_name": "setyo",
  "user_profile": {
    "first_name": "Setyo",
    "last_name": "Hadee"
  }
}
```

### Convert ke `camelCase`

```go
out, _ := jsonparser.TransformCaseJSON(data, CamelCase)
```

### Output

```json
{
  "userName": "setyo",
  "userProfile": {
    "firstName": "Setyo",
    "lastName": "Hadee"
  }
}
```

### Karakteristik

* Streaming rewrite
* Tidak membangun AST
* Key di-unescape → di-transform → di-escape ulang
* Value **tidak disentuh**

### Cocok untuk:

* API gateway
* Middleware
* Payload normalization

### Implementasi inti:

```go
transformObject
transformArray
```

---

## Zero-Allocation Heavy Path Traversal

Traversal JSON dilakukan **tanpa alokasi string / map / struct**.

### Prinsip

* Path dipecah ke `[]pathToken`
* Token `key` adalah **subslice langsung dari path**
* JSON value direferensikan sebagai subslice dari input

```go
type pathToken struct {
  key        []byte
  index      int
  isIdx      bool
  isWildcard bool
}
```

### Contoh

```go
r := jsonparser.Get(data, "a.b[3].c")
```

Yang terjadi:

* ❌ Tidak ada `map[string]interface{}`
* ❌ Tidak ada `string(path)`
* ❌ Tidak ada reflect
* ✅ Pure byte scanning

### Dampak

* Sangat cepat untuk payload besar
* GC pressure minimal
* Cocok untuk high-throughput service

### Implementasi inti:

```go
splitPathBytes
getNestedValue
extractValue
```

---

## Feature Comparison

| Feature                   | encoding/json | gjson | jsonparser |
| ------------------------- | ------------- | ----- | ---------- |
| Comment support (`//`)    | ❌             | ❌     | ✅          |
| Wildcard array traversal  | ❌             | ✅     | ✅          |
| Streaming key transform   | ❌             | ❌     | ✅          |
| Zero-allocation traversal | ❌             | ✅     | ✅          |
| Upsert by path            | ❌             | ❌     | ✅          |

---

##  Design Philosophy

* **Read & write JSON as bytes**
* **No reflection**
* **No struct binding**
* **Streaming-first**
* **Composable primitives**

Library ini ditujukan untuk:

* API Gateway
* ETL / JSON reshaping
* Middleware
* Performance-critical services

---


