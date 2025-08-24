# POC JsonParser
Bertujuan untuk efisiensi akses data JSON, memotong jalur `unmarshall` mengabaikan tokenisasi hanya berfokus pada pengambilan data dari key(s). _V3 bertujuan untuk mengurangi alokasi memori untuk menekan beban garbage collector (GC)

# Cepat, Hemat, Stabil, Praktis
1. Ambil JSON langsung ke Key, tidak `unmarshal` dan dijadikan ke objek. Langsung mengembalikan value
2. Enggak bikin boros RAM
3. Bikin aplikasi tetap stabil
4. Praktis, siapin `Result` tinggal `Get(JsonString, Key)` atau `GetMany(JsonString, Keys)`

# untuk kamu yang nerdy

## Perubahan Versi
1. v2
- Mengandalkan operasi `string`
- Alokasi memoriyang masih banyak
- Beban GC tinggi

2. v3.1
- Menekan drastis alokasi memori
- Menghilangkan operasi `string`

3. v3.2
- Optimasi spesifik `getMany`

## Tabel Optimasi
| Versi    | Test    | Ops (N) | ns/op | B/op  | allocs/op |
| -------- | ------- | ------- | ----- | ----- | --------- |
| **v2**   | get     | 382,465 | 3,205 | 3,856 | 39        |
|          | getmany | 165,984 | 6,839 | 4,496 | 99        |
| **v3.1** | get     | 383,431 | 2,982 | 1,224 | 5         |
|          | getmany | 244,556 | 5,031 | 3,192 | 30        |
| **v3.2** | get     | 338,244 | 3,186 | 1,216 | 4         |
|          | getmany | 325,588 | 3,772 | 2,568 | 10        |

## Catatan Perubahan
### v2 → v3.1

get
```
ns/op: turun dari 3,205 → 2,982 (~7% lebih cepat)
B/op: turun drastis 3,856 → 1,224 (-68%)
allocs/op: 39 → 5 (-87%)
```

getmany
```
ns/op: 6,839 → 5,031 (~26% lebih cepat)
B/op: 4,496 → 3,192 (-29%)
allocs/op: 99 → 30 (-70%)
```
✅ Lonjakan besar dalam pengurangan alokasi memori dan beban GC.

### v3.1 → v3.2

get
```
ns/op: stabil (2,982 → 3,186, perbedaan kecil karena noise)
B/op: hampir sama (1,224 → 1,216)
allocs/op: 5 → 4 (-20%)
```
getmany
```
ns/op: 5,031 → 3,772 (~25% lebih cepat)
B/op: 3,192 → 2,568 (-20%)
allocs/op: 30 → 10 (-67%)
```
✅ Optimasi batch (getmany) berhasil: jauh lebih cepat dan jauh lebih hemat alokasi.

## Kesimpulan

- v2 → v3.1: Lompatan terbesar di efisiensi memori dan penurunan allocs/op → dampak langsung ke stabilitas performa.

- v3.1 → v3.2: Fokus ke getmany: throughput batch naik ~25%, allocs/op turun 3×, B/op juga lebih hemat.

## Secara keseluruhan:
Implementasi v3.2 jauh lebih scalable untuk high-throughput service dibanding v2.