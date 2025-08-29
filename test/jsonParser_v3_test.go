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
		jsonParser.GetMany([]byte(configuration.JsonStr), []string{`action`, `serviceOrderJobs[1].price`, `vehicle.year`})
	}
}

jsonArr := 
`
	[
		{
			"Name1": "adi"
		},
		{
			"Name1": "aya"
		}
	]
`

func TestGetArrays(t *testing.T) {
	res := jsonParser.Get([]byte(configuration.JsonStr), "[0].Name1")
	if res.Data != "adi" {
		t.Errorf("Failed to get")
	}
	fmt.Sprintf("nangis")
}
