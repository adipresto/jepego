package test

import (
	"poc_jsonParser/configuration"
	"poc_jsonParser/utils/jsonParser_v2"
	"testing"
)

func TestGetField_First(t *testing.T) {
	// test ambil key/kolom diakhir JSON
	res := jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Raw != `"2025-01-23 10:00:00"` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", res.Raw)
	}
}
func TestGetField_Last(t *testing.T) {
	// test ambil key/kolom diakhir JSON
	res := jsonParser_v2.Get(configuration.JsonStr, "action", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Raw != `"GR"` {
		t.Errorf("expected \"GR\", got %s", res.Raw)
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

func TestGetSubselectorObject(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "{action,vehicle.year,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: false})
	expected := `{"action":"GR","vehicle.year":2015,"bookingDateTime":"2025-01-23 10:00:00"}`
	if res.Raw != expected {
		t.Errorf("expected %s, got %s", expected, res.Raw)
	}
	res = jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[1].price,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: false})
	expected = `{"action":"GR","serviceOrderJobs[1].price":"383000.0","bookingDateTime":"2025-01-23 10:00:00"}`
	if res.Raw != expected {
		t.Errorf("expected %s, got %s", expected, res.Raw)
	}
}

func TestGetSubselectorObject_false(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs.year,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: false})
	expected := `{"action":"GR","bookingDateTime":"2025-01-23 10:00:00"}`
	if res.Raw != expected {
		t.Errorf("expected %s, got %s", expected, res.Raw)
	}
}

func TestGetSubselectorArray(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "[0,2]", jsonParser_v2.ParserOption{IsRaw: false})
	expected := `[10,30]`
	if res.Raw != expected {
		t.Errorf("expected %s, got %s", expected, res.Raw)
	}
}

func TestGetModifierStatic(t *testing.T) {
	res := jsonParser_v2.Get(configuration.JsonStr, "@name", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Raw != `"Alice"` {
		t.Errorf("expected \"Alice\", got %s", res.Raw)
	}

	res = jsonParser_v2.Get(configuration.JsonStr, "!age", jsonParser_v2.ParserOption{IsRaw: false})
	if res.Raw != "30" {
		t.Errorf("expected 30, got %s", res.Raw)
	}
}
