package test

import (
	"poc_jsonParser/configuration"
	"poc_jsonParser/utils/jsonParser_v2"
	"strings"
	"testing"
)

// Mau ambil kolom diawal []byte/string atau diakhir. Tetap waktu eksekusi adalah sama

// ambil nilai berdasarkan key/kolom,
// default mengirimkan koleksi dan berguna untuk mengambil lebih dari satu key
func TestGetField(t *testing.T) {
	// test ambil key/kolom diakhir JSON
	res := jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime")
	if res.Collection[0] != `"2025-01-23 10:00:00"` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Collection[0])
	}
}

// mode raw = mengembalikan nilai dalam bentuk string.
// Kekurangannya terdapat operasi append string
func TestGetField_Raw(t *testing.T) {
	// test ambil key/kolom diakhir JSON
	res := jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime", jsonParser_v2.ParserOption{IsRaw: true})
	if res.Raw != `"2025-01-23 10:00:00"` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Raw)
	}
}

// Perbedaan waktu eksekusi antara Raw = true dengan false. Mengambil kolom diakhir JSON

func BenchmarkTestGetField(b *testing.B) {
	// b.N = 1_000_000
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "action")
	}
}

func BenchmarkTestGetField_Raw(b *testing.B) {
	// b.N = 1_000_000
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "action", jsonParser_v2.ParserOption{IsRaw: true})
	}
}

func TestGetFields_IsRawFalse(t *testing.T) {
	// test ambil key/kolom di akhir dan diawal JSON
	// contoh penggunaan yang tidak efisien
	res := jsonParser_v2.Get(configuration.JsonStr, "action", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Collection[0] != `"GR"` {
		t.Errorf("expected \"GR\", got %s", res.Collection[0])
	}

	res = jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Collection[0] != `"2025-01-23 10:00:00"` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Collection[0])
	}
}
func TestGetFields(t *testing.T) {
	// test ambil key/kolom di akhir dan diawal JSON
	// contoh penggunaan yang tidak efisien
	res := jsonParser_v2.Get(configuration.JsonStr, "action", jsonParser_v2.ParserOption{IsRaw: true})
	if res.Raw != `"GR"` {
		t.Errorf("expected \"GR\", got %s", res.Raw)
	}

	res = jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime", jsonParser_v2.ParserOption{IsRaw: true})
	if res.Raw != `"2025-01-23 10:00:00"` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Raw)
	}
}

func TestGetNestedField(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "customer.Name1", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Raw != `"1dxmcz"` {
		t.Errorf("expected \"1dxmcz\", got %s", res.Raw)
	}
}

func TestGetArrayIndex(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "serviceOrderJobs[1]", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Raw != "{\"duration\":5,\"jobCode\":\"OIL\",\"jobDescription\":\"Engine Oil\",\"price\":\"383000.0\"}" {
		t.Errorf("expected \"{\"duration\":5,\"jobCode\":\"OIL\",\"jobDescription\":\"Engine Oil\",\"price\":\"383000.0\"}\", got %s", res.Raw)
	}
}

// Much easier when ParserOption.IsRaw is set to false
// get the value based on key position index
func TestGetSubselectorObject(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "{action,vehicle.year,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: true})

	expected := `{"action":"GR","vehicle.year":2015,"bookingDateTime":"2025-01-23 10:00:00"}`
	if res.Raw != expected {
		t.Errorf("expected %s, got %s", expected, res.Raw)
	}
	res = jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[1].price,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: false})
	// 0. action = "GR"
	// 1. serviceOrderJobs[1].price = "383000.0"
	// 2. bookingDateTime = "2025-01-23 10:00:00"
	expected_collection := []string{`"GR"`, `"383000.0"`, `"2025-01-23 10:00:00"`}
	if strings.Join(res.Collection, ",") != strings.Join(expected_collection, ",") {
		t.Errorf("expected %s, got %s", expected_collection, res.Collection)
	}
}

func TestGetSubselectorArray(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "[action,vehicle.year]", jsonParser_v2.ParserOption{IsRaw: false})
	expected := []string{`"GR"`, `2015`}
	if strings.Join(res.Collection, ",") != strings.Join(expected, ",") {
		t.Errorf("expected %s, got %s", expected, res.Collection)
	}
}

func TestGetModifierStatic(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "@action", jsonParser_v2.ParserOption{IsRaw: true})
	if res.Raw != `"GR"` {
		t.Errorf("expected \"GR\", got %s", res.Collection[0])
	}

	res = jsonParser_v2.Get(configuration.JsonStr, "!vehicle.year", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Collection[0] != "2015" {
		t.Errorf("expected 2015, got %s", res.Collection[0])
	}
}
