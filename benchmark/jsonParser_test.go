package benchmark

import (
	"poc_jsonParser/configuration"
	"poc_jsonParser/utils/jsonParser_v2"
	"testing"
)

func BenchmarkGetField_Last(b *testing.B) {
	count := b.N
	for i := 0; i < count; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "action")
	}
}
func BenchmarkGetField_First(b *testing.B) {
	count := b.N
	for i := 0; i < count; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime")
	}
}
func BenchmarkGetField_First_Once(b *testing.B) {
	jsonParser_v2.Get(configuration.JsonStr, "bookingDateTime")
}
func BenchmarkGetFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[3].price,bookingDateTime}")
	}
}
func BenchmarkGetFields_Once(b *testing.B) {
	jsonParser_v2.Get(configuration.JsonStr, "{action,serviceOrderJobs[3].price,bookingDateTime}")
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
