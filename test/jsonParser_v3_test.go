package test

import (
	"poc_jsonParser/configuration"
	jsonParser "poc_jsonParser/utils/jsonParser_v3"
	"testing"
)

// Mau ambil kolom diawal []byte/string atau diakhir. Tetap waktu eksekusi adalah sama

// mode raw = mengembalikan nilai dalam bentuk string.
// Kekurangannya terdapat operasi append string
func TestGetField_v3(t *testing.T) {
	// test ambil key/kolom diakhir JSON
	res := jsonParser.Get([]byte(configuration.JsonStr), "bookingDateTime")
	if string(res.Data) != `"2025-01-23 10:00:00"` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", string(res.Data))
	}
	if res.Key != "bookingDateTime" {
		t.Errorf("expected \"bookingDateTime\", got %s", res.Key)
	}
}
func TestGetNestedField_v3(t *testing.T) {
	res := jsonParser.Get([]byte(configuration.JsonStr), "customer.Name1")
	if string(res.Data) != `1dxmcz` {
		t.Errorf("expected \"1dxmcz\", got %s", string(res.Data))
	}
	res = jsonParser.Get([]byte(configuration.JsonStr), "vehicle.year")
	if string(res.Data) != `2015` {
		t.Errorf("expected 2015, got %s", string(res.Data))
	}
}

func TestGetArrayIndex_v3(t *testing.T) {
	res := jsonParser.Get([]byte(configuration.JsonStr), "serviceOrderJobs[1].price")
	if string(res.Data) != `"383000.0"` {
		t.Errorf("expected \"383000.0\", got %s", string(res.Data))
	}
}

func TestGetSubselectorObject_v3(t *testing.T) {
	res := jsonParser.GetMany([]byte(configuration.JsonStr), []string{`action`, `vehicle.year`, `bookingDateTime`})

	expected := `"GR"`
	if string(res[0].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res[0].Data))
	}
	expected = "2015"
	if string(res[1].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res[1].Data))
	}
}

// ==== BENCHMARK === //
// Perbedaan waktu eksekusi antara Raw = true dengan false. Mengambil kolom diakhir JSON
func BenchmarkGetSubselectorObject_v3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// jsonParser.GetMany([]byte(configuration.JsonStr), []string{`action`, `serviceOrderJobs[1].price`, `vehicle.year`})
		jsonParser.GetMany([]byte(configuration.JsonStr), []string{`action`, `branchCode`})
	}
}

// func TestGetFields_IsRawFalse(t *testing.T) {
// 	// test ambil key/kolom di akhir dan diawal JSON
// 	// contoh penggunaan yang tidak efisien
// 	res := jsonParser.Get(configuration.JsonStr, "action", jsonParser.ParserOption{IsRaw: false})
// 	if res.Collection[0] != `"GR"` {
// 		t.Errorf("expected \"GR\", got %s", res.Collection[0])
// 	}

// 	res = jsonParser.Get(configuration.JsonStr, "bookingDateTime", jsonParser.ParserOption{IsRaw: false})
// 	if res.Collection[0] != `"2025-01-23 10:00:00"` {
// 		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Collection[0])
// 	}
// }
// func TestGetFields(t *testing.T) {
// 	// test ambil key/kolom di akhir dan diawal JSON
// 	// contoh penggunaan yang tidak efisien
// 	res := jsonParser.Get(configuration.JsonStr, "action", jsonParser.ParserOption{IsRaw: true})
// 	if res.Raw != `"GR"` {
// 		t.Errorf("expected \"GR\", got %s", res.Raw)
// 	}

// 	res = jsonParser.Get(configuration.JsonStr, "bookingDateTime", jsonParser.ParserOption{IsRaw: true})
// 	if res.Raw != `"2025-01-23 10:00:00"` {
// 		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Raw)
// 	}
// }

// func TestGetSubselectorArray(t *testing.T) {
// 	res := jsonParser.Get(configuration.JsonStr, "[action,vehicle.year]", jsonParser.ParserOption{IsRaw: false})
// 	expected := []string{`"GR"`, `2015`}
// 	if strings.Join(res.Collection, ",") != strings.Join(expected, ",") {
// 		t.Errorf("expected %s, got %s", expected, res.Collection)
// 	}
// }

// func TestGetModifierStatic(t *testing.T) {
// 	res := jsonParser.Get(configuration.JsonStr, "@action", jsonParser.ParserOption{IsRaw: true})
// 	if res.Raw != `"GR"` {
// 		t.Errorf("expected \"GR\", got %s", res.Collection[0])
// 	}

// 	res = jsonParser.Get(configuration.JsonStr, "!vehicle.year", jsonParser.ParserOption{IsRaw: false})
// 	if res.Collection[0] != "2015" {
// 		t.Errorf("expected 2015, got %s", res.Collection[0])
// 	}
// }
