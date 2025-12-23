# JEPEGO Traversal Architecture

Library ini menggunakan **byte-level traversal engine** tanpa decoding ke AST atau struct.

---

## High-Level Traversal Flow

```
┌────────────────────────────┐
│  Public API (Get / Delete) │
└──────────────┬─────────────┘
               │
               ▼
┌────────────────────────────┐
│ removeCommentsBytes        │
│ - strip // outside string │
└──────────────┬─────────────┘
               │
               ▼
┌────────────────────────────┐
│ splitPathBytes             │
│ "a.b[3].c" → pathToken[]  │
└──────────────┬─────────────┘
               │
               ▼
┌────────────────────────────┐
│ Traversal Engine           │
│ getNestedValue / Values   │
└──────────────┬─────────────┘
               │
               ▼
┌────────────────────────────┐
│ extractValue               │
│ returns raw JSON slice     │
└────────────────────────────┘
```

### **Catatan penting**

* Tidak ada `map[string]interface{}`
* Tidak ada `encoding/json`
* Semua operasi berbasis `[]byte`

---

## Traversal Engine (Path-driven)

Traversal dilakukan **token per token**, bukan recursive parsing penuh.

### Path Token Model

```
"a.b[3].c"

↓ splitPathBytes

[
  { key: "a" },
  { key: "b" },
  { isIdx: true, index: 3 },
  { key: "c" }
]
```

---

### Traversal Decision Tree

```
getNestedValue(json, tokens)
        │
        ▼
┌────────────────────────────┐
│ cur := trimSpace(json)     │
└──────────────┬─────────────┘
               │
               ▼
     ┌────────────────────┐
     │ cur[0] == '{' ?    │
     └───────┬────────────┘
             │yes
             ▼
   ┌─────────────────────────┐
   │ getTopLevelKey          │
   │ - scan object keys     │
   │ - compare bytes        │
   └──────────┬─────────────┘
              │
              ▼
     next token → loop
```

```
     ┌────────────────────┐
     │ cur[0] == '[' ?    │
     └───────┬────────────┘
             │yes
             ▼
   ┌─────────────────────────┐
   │ getArrayIndex           │
   │ OR wildcard []          │
   └──────────┬─────────────┘
              │
              ▼
     next token → loop
```

 **Traversal berhenti lebih awal** jika:

* Type tidak cocok (`object vs array`)
* Key tidak ditemukan
* Index out of bounds

---

## Wildcard Traversal (`[]`) Expansion

Wildcard tidak mengubah traversal engine, hanya **menggandakan jalur**.

```
users[].name
```

### Flow

```
Array [
  { name: "A" },
  { name: "B" }
]

↓ wildcard

getNestedValues
│
├── traverse element[0] → name
│
└── traverse element[1] → name
```

### Diagram

```
┌───────────────┐
│ users[]       │
└───────┬───────┘
        │
 ┌──────┴───────┐
 │              │
 ▼              ▼
user[0]       user[1]
 │              │
 ▼              ▼
name           name
```

### Implementasi:

```go
getNestedValues
extractValue
```

---

## Low-Level JSON Value Scanner

Semua traversal **bergantung pada satu primitive inti**:

```
extractValue
```

### Prinsip Kerja

```
Input:  { "a": [1, {"b":2}] }
        ^
        pointer
```

```
┌───────────────┐
│ detect type   │
│ '{' '[' '"'  │
└───────┬───────┘
        │
        ▼
┌──────────────────────────┐
│ depth counter            │
│ inString / escaped flag │
└──────────┬──────────────┘
           │
           ▼
┌──────────────────────────┐
│ return raw slice         │
│ without allocation      │
└──────────────────────────┘
```

### Karakteristik:

* Depth-based scanning
* Aman terhadap nested object/array
* String-aware (`\"`, `\\`)
* Return **subslice**, bukan copy

---

## Delete & Upsert Reuse Traversal

Delete dan Upsert **tidak punya traversal sendiri**.

```
DeleteField / Upsert
        │
        ▼
getTopLevelKey
        │
        ▼
extractValue
        │
        ▼
rebuild JSON fragment
```

### Ini memastikan:

* Satu sumber kebenaran traversal
* Konsistensi behavior
* Mudah diuji

---

## Mental Model (Ringkas)

> “JSON adalah byte stream.
> Path adalah navigator.
> extractValue adalah pisau bedah.”

---

##  Summary

* Traversal **token-driven**
* Parsing **on-demand**
* Tidak ada global parse
* Tidak ada intermediate tree
* Semua operasi **O(n)** terhadap fragment yang disentuh
