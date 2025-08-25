package benchmark

import (
	"poc_jsonParser/configuration"
	"poc_jsonParser/utils/jsonParser"
	"poc_jsonParser/utils/jsonParser_v2"
	"testing"
)

func BenchmarkGetField_v1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser.Get([]byte(configuration.JsonStrClean), "action")
	}
}
func BenchmarkGetFields_v1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// a :=
		jsonParser.GetMany([]byte(configuration.JsonStrClean), []string{"action", "branchCode", "vehicle.year"})
		// fmt.Printf("%s, %s, %s", a[`action`], a[`branchCode`], a[`vehicle.year`])
	}
}
func BenchmarkGetField_v2_Last(b *testing.B) {
	count := b.N
	for i := 0; i < count; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "action")
	}
}
func BenchmarkGetField_First(b *testing.B) {
	count := b.N
	for i := 0; i < count; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime", jsonParser_v2.ParserOption{IsRaw: false})
	}
}
func BenchmarkGetField_First_Once(b *testing.B) {
	jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime", jsonParser_v2.ParserOption{IsRaw: false})
}
func BenchmarkGetFields_Raw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[3].price,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: true})
	}
}
func BenchmarkGetFields_Collection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[3].price,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: false})
	}
}
func BenchmarkGetFields_Once(b *testing.B) {
	jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[3].price,bookingDateTime}", jsonParser_v2.ParserOption{IsRaw: false})
}
func BenchmarkGetNestedField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "customer.Name1")
	}
}

func BenchmarkGetArrayIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "serviceOrderJobs[1]")
	}
}

func BenchmarkGetSubselectorObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "{action,bookingDateTime}")
	}
}

func BenchmarkGetSubselectorArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "[0,2]")
	}
}

func BenchmarkGetModifierStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "@name")
		jsonParser_v2.Get(configuration.JsonStr, "!age")
	}
}
