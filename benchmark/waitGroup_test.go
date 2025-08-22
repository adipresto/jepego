package benchmark

import (
	"poc_jsonParser/configuration"
	"poc_jsonParser/utils/jsonParser_v2"
	"sync"
	"testing"
)

func BenchmarkWaitGroup(b *testing.B) {
	var wg sync.WaitGroup
	res := make([]jsonParser_v2.Result, 3)
	for i := 0; i < b.N; i++ {
		// mau ambil 3 nilai dari kata-kunci dengan sub task
		// tambahkan subtasks dulu
		wg.Add(3)
		//
		go func() {
			defer wg.Done()
			res[0] = jsonParser_v2.Get(configuration.JsonStr, "action")
		}()
		go func() {
			defer wg.Done()
			res[1] = jsonParser_v2.Get(configuration.JsonStr, "serviceOrderJobs[1].price")
		}()
		go func() {
			defer wg.Done()
			res[2] = jsonParser_v2.Get(configuration.JsonStr, "vehicle.year")
		}()

		// tunggu sampai semua goroutine selesai
		wg.Wait()
	}
	// cleanup
	res = nil
}
func BenchmarkNonWaitGroup(b *testing.B) {
	res := make([]jsonParser_v2.Result, 3)
	for i := 0; i < b.N; i++ {
		res[0] = jsonParser_v2.Get(configuration.JsonStr, "action")
		res[1] = jsonParser_v2.Get(configuration.JsonStr, "serviceOrderJobs[1].price")
		res[2] = jsonParser_v2.Get(configuration.JsonStr, "vehicle.year")
	}
	// cleanup
	res = nil
}
