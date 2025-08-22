package test

import (
	"poc_jsonParser/configuration"
	"poc_jsonParser/utils/jsonParser_v2"
	"sync"
	"testing"
)

func Test_waitGroup(t *testing.T) {
	var wg sync.WaitGroup
	res := make([]jsonParser_v2.Result, 3)

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

	expected := `"GR"`
	if res[0].Raw != expected {
		t.Errorf("expected %s, got %s", expected, res[0].Raw)
	}
	expected = `"383000.0"`
	if res[1].Raw != expected {
		t.Errorf("expected %s, got %s", expected, res[1].Raw)
	}
	expected = `2015`
	if res[2].Raw != expected {
		t.Errorf("expected %s, got %s", expected, res[2].Raw)
	}
}
