package test

import (
	"fmt"
	"testing"

	"github.com/adipresto/jepego/configuration"
	"github.com/adipresto/jepego/utils/jsonParser"
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
	if string(res["action"].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res["action"].Data))
	}
	expected = "2015"
	if string(res["vehicle.year"].Data) != expected {
		t.Errorf("expected %s, got %s", expected, string(res["vehicle.year"].Data))
	}
}

// ==== BENCHMARK === //
// Perbedaan waktu eksekusi antara Raw = true dengan false. Mengambil kolom diakhir JSON
func BenchmarkGetSubselectorObject_v3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonParser.GetMany([]byte(configuration.JsonStr), []string{`action`, `serviceOrderJobs[1].price`, `vehicle.year`})
	}
}

const jsonArr = `
	[
		{
			"Name1": "adi"
		},
		{
			"Name1": "aya"
		}
	]
`

func TestGetArr(t *testing.T) {
	res := jsonParser.GetMany([]byte(jsonArr), []string{"[].Name1"})
	fmt.Printf("\n%s\n", res["[].Name1"].Data)
	if string(res["[].Name1"].Data) != "adi" {
		t.Errorf("Failed to get")
	}
	fmt.Printf("\nnangis\n")
}

func TestGetArrays(t *testing.T) {
	res := jsonParser.Get([]byte(jsonArr), "[0].Name1")

	if string(res.Data) != "adi" {
		t.Errorf("Failed to get")
	}
	fmt.Printf("\nnangis\n")
}

func TestGetObject(t *testing.T) {
	res := jsonParser.Get([]byte(configuration.JsonStr), "customer")

	fmt.Println(string(res.Data))
}

func TestGetAnotherObject(t *testing.T) {
	grjson := configuration.Grjson
	res := jsonParser.Get([]byte(grjson), "data.one_account")                              // object
	res1 := jsonParser.Get([]byte(grjson), "process")                                      // "service_progress_update"
	res2 := jsonParser.Get([]byte(grjson), "actual_mileage")                               // undefined
	res3 := jsonParser.Get([]byte(grjson), "data.one_account.email2")                      // null
	res4 := jsonParser.Get([]byte(grjson), "data.payment[1].payment_by")                   // "DEALERB"
	res4a := jsonParser.GetAll([]byte(grjson), "data.payment[].payment_by")                // "DEALER"
	res4many := jsonParser.GetMany([]byte(grjson), []string{"data.payment[1].payment_by"}) // "DEALERB"
	res5 := jsonParser.Get([]byte(grjson), "data.payment")                                 // Array of object

	fmt.Println(string(res.Data))
	fmt.Println(string(res1.Data))
	fmt.Println(string(res2.Data))
	fmt.Println(string(res3.Data))
	fmt.Println(string(res4.Data))
	fmt.Println(len(res4a))
	for _, v := range res4a {
		fmt.Println(string(v.Data)) // "DEALER", "DEALER"
	}
	for _, v := range res4many {
		fmt.Println(string(v.Data))
	}
	fmt.Println(string(res5.Data))

}

var jsonStrEnc string = `
	{
  "enc": "QWkTSCHD67OxxU6kbrMtvF1NwnWW4gtYBIYb7t8lSi26wCfMwlFxlSQZE94CikxJURjNYdwvjTgTM0YUwnrg9jlXcLDYOhuOpc6VTcJlqhuHxAnPqwq8G3GWVoVkKFmSW2D1J6eV7dvvE6HLstQBEsnRZPdOc5IhkZITMP4LlRSA3v3IvTqpFQlfylWfVdexdu1a+tmKjnGALNggluB/jwgDC7WH1qSqXoiWaJztjwLHBMLOBQmJEo9CNafYgWYgDuYCwlfo4IQsmc0VfUzbQkZ7KTDzlvlIv3gGBtHp2XKbuxPd9wYJLl1ie+u9tGtPL9V0YlLeEkJbL8KlSdh6qalN3szTk6e+FT1YCP2YKNoduE13b/HqmH+UGPm6j3NWY9FAx9883y9lFDKpgWJ02jmXDt87SbNurqbi8kPEHkco21tP4HTO4ViPyPEYKtzkBTQm26kPvl+wfW6p0GaAT+2TvxelDQI9ttsoNvtXqnXqkP7Nze/29zlcL3YsG+wIe5fTrCBALhkhk9LChIxT7aJTXmyllYdRC30Iw67LyTWrr/VpZpSsDqljr+ND4ryHUVtIKt4JarLT+SaaqJM+FbATRFNhEnQ1EQInPt+GhAeBF6/3egEar/8PObbW1zXbuUP4U9BzquK++wA+CLPIPBJJw7L09uwhAxooSFvhK6jLkvEzwwMtOrXAH8bqUgBrzNRCkZWuqysPCOrdp0+MB3VYEuWad8wXN0QODF42a09tMs54tp1UI3OAMvWjnSmlSYg3UrTRZppbprE+zjeAPxiwLBuDpdxXt46l2mDsiTHvhkCbZaVgNgAxJRgJfEyxojRBno0SbREwu/WfW2ZimUfB6OxbePcobHbQQRES//lWh4OKPOlFP2GAqpYkvw3U8js7xiIqgNiSzBK7UFFxOhG52ySr+I7OK1NV2teK3ddsKo7fyYTi6d1wSN6G+JDpyoqgb6+egavUTwumlI7rwTsTzyEFr8uWjHGlTRxoipdfexUmC+ITWZ/AX627KhxoqNOZ9eAorGFSnzskqy9ODj8jtyow7FxR39qmXXvHqDlsN5JFANY6h0hHzPVuWQ5XhNtZGcpQe8QdnKzU0qR1IIBLFs6LSvpho/90oCYOSHhnFAh6Ll1G2AkPBTTJSbs8azOhFNT6VotSrZWYfomWlm0n7sTZy9dHWAeR3ch2E02aTI+ShCv59AH0fenCBmDkt1SlX2uOzIBOAc0I+Ky9IUKwvrSg7Z8RRoYQPV62PzCPTLo1RPYKjPF437RG5vj2rRisZ7vtwbf8knpsEV9+QAuVLT6LjpK/eliPnNaVKnA="
}
`

func TestEnc(t *testing.T) {
	res := jsonParser.Get([]byte(jsonStrEnc), `enc`)
	fmt.Printf(string(res.Data))
}
