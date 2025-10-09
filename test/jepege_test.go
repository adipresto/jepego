package tester

import (
	"fmt"
	"testing"

	"github.com/adipresto/jepego/apis"
	"github.com/adipresto/jepego/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetField(t *testing.T) {
	// test ambil key/kolom diakhir JSON
	res := apis.Get([]byte(jsonStr), "bookingDateTime")
	if string(res.Data) != `2025-01-23 10:00:00` {
		t.Errorf("expected \"2025-01-23 10:00:00\", got %s", string(res.Data))
	}
	if res.Key != "bookingDateTime" {
		t.Errorf("expected \"bookingDateTime\", got %s", res.Key)
	}
}
func TestGetNestedField(t *testing.T) {
	res := apis.Get([]byte(jsonStr), "customer.Name1")
	if string(res.Data) != `"1dxmcz"` {
		t.Errorf("expected \"1dxmcz\", got %s", string(res.Data))
	}
}

func TestGetArrayIndex(t *testing.T) {
	res := apis.Get([]byte(jsonStr), "serviceOrderJobs[1].price")
	if string(res.Data) != `383000.0` {
		t.Errorf("expected \"383000.0\", got %s", string(res.Data))
	}
}

func TestGetSubselectorObject(t *testing.T) {
	res := apis.GetMany([]byte(jsonStr), []string{`action`, `vehicle.year`, `bookingdateTime`, `event_id`})

	expected := `GR`
	// field is Action, we fetch action
	if string(res[`action`].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res[`action`].Data))
	}
	expected = `2015`
	if string(res[`vehicle.year`].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res[`vehicle.year`].Data))
	}

	// field is BOOKingDateTime, we fetch bookingdateTime
	expected = `2025-01-23 10:00:00`
	if string(res[`bookingdateTime`].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res[`bookingdatetime`].Data))
	}

	expected = `123abcdef7890`
	if string(res[`event_id`].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res[`event_id`].Data))
	}
}

func BenchmarkGetSubselectorObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		apis.GetMany([]byte(jsonStr), []string{`action`, `branchCode`})
	}
}

// Menentukan jenis Getter dari koleksi field_tokenisasi
func TestSpecifyGetterFields(t *testing.T) {
	field_detokenizes := []string{
		"data.payment[].payment_by",                           //
		"data.payment[].payment_stage[].payment_status",       // siksa lagi
		"data.delivery.service_delivery_handover_datetime",    //
		"data.delivery.service_vehicle_receiver_phone_number", //
	}

	// cukup panggil GetManyAll sekali
	results := apis.GetManyAll([]byte(jsonStrArr), field_detokenizes)

	for i, r := range results {
		fmt.Printf("%s: \n", i)
		for _, v := range r {
			fmt.Printf("key=%s, val=%s\n", v.Key, v.Data)
		}
	}
}
func TestUpsert(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		path     string
		value    []byte
		wantJSON string
		wantType utils.DataType
	}{
		{
			name:     "Add new string field",
			input:    []byte(`{"a":1}`),
			path:     "b",
			value:    []byte(`hello`),
			wantJSON: `{"a":1,"b":"hello"}`,
			wantType: utils.TypeString,
		},
		{
			name:     "Overwrite number field",
			input:    []byte(`{"a":1}`),
			path:     "a",
			value:    []byte(`42`),
			wantJSON: `{"a":42}`,
			wantType: utils.TypeNumber,
		},
		{
			name:     "Add bool field",
			input:    []byte(`{"obj":{}}`),
			path:     "obj.flag",
			value:    []byte(`true`),
			wantJSON: `{"obj":{"flag":true}}`,
			wantType: utils.TypeBool,
		},
		{
			name:     "Add null field",
			input:    []byte(`{"obj":{}}`),
			path:     "obj.none",
			value:    []byte(`null`),
			wantJSON: `{"obj":{"none":null}}`,
			wantType: utils.TypeNull,
		},
		{
			name:     "Add object field",
			input:    []byte(`{"a":{}}`),
			path:     "a.obj",
			value:    []byte(`{"x":1}`),
			wantJSON: `{"a":{"obj":{"x":1}}}`,
			wantType: utils.TypeObject,
		},
		{
			name:     "Add array field",
			input:    []byte(`{"a":{}}`),
			path:     "a.arr",
			value:    []byte(`[1,2,3]`),
			wantJSON: `{"a":{"arr":[1,2,3]}}`,
			wantType: utils.TypeArray,
		},
		{
			name:     "Overwrite array element",
			input:    []byte(`{"arr":[{"x":1},{"x":2}]}`),
			path:     "arr[1].x",
			value:    []byte(`99`),
			wantJSON: `{"arr":[{"x":1},{"x":99}]}`,
			wantType: utils.TypeNumber,
		},
		{
			name:     "Nested array element",
			input:    []byte(`{"a":[[{"b":1}]]}`),
			path:     "a[0][0].b",
			value:    []byte(`10`),
			wantJSON: `{"a":[[{"b":10}]]}`,
			wantType: utils.TypeNumber,
		},
		{
			name:     "Pad array and insert",
			input:    []byte(`{"arr":[]}`),
			path:     "arr[2]",
			value:    []byte(`pad`),
			wantJSON: `{"arr":[null,null,"pad"]}`,
			wantType: utils.TypeString,
		},
		{
			name:     "Insert into deeply nested array",
			input:    []byte(`{"a":[[],[]]}`),
			path:     "a[1][0]",
			value:    []byte(`123`),
			wantJSON: `{"a":[[],[123]]}`,
			wantType: utils.TypeNumber,
		},
		{
			name:     "Overwrite object with array",
			input:    []byte(`{"a":{"b":1}}`),
			path:     "a",
			value:    []byte(`[1,2,3]`),
			wantJSON: `{"a":[1,2,3]}`,
			wantType: utils.TypeArray,
		},
		{
			name:     "Overwrite array with object",
			input:    []byte(`{"a":[1,2,3]}`),
			path:     "a",
			value:    []byte(`{"x":9}`),
			wantJSON: `{"a":{"x":9}}`,
			wantType: utils.TypeObject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := apis.Upsert(tt.input, tt.path, tt.value, tt.wantType)

			// cek JSON string (tidak peduli urutan key)
			assert.JSONEq(t, tt.wantJSON, string(got))

			// cek DataType
			val := apis.Get(got, tt.path)
			assert.True(t, val.OK, "expected value to exist after upsert")
			assert.Equal(t, tt.wantType, val.DataType)
		})
	}
}
