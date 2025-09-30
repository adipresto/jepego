package test

import (
	"fmt"
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
	fmt.Printf("\n%s\n", res[1].Data)
	if string(res[1].Data) != "adi" {
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
	grjson := `
{
  "process": "service_progress_update",
  "event_ID": "0d5be854-7f4a-4ed0-be00-da098d3420e3",
  "timestamp": 1706726960,
  "data": {
    "one_account": {
      "one_account_ID": "GMO4GNYBSI0D85IP6K59OYGJZ6VOKW3Y",
      "dealer_customer_ID": "ASTVAJMF000552",
      "first_name": "Nkoc",
      "last_name": "Maf",
      "phone_number": "MFV8O1O4y46k",
      "email": "pgb.terfsxgfsnwu@v",
      "email2": null
    },
    "customer_vehicle": {
      "vin": "MKFKZE81SCJ115045",
      "katashiki_suffix": "NSP170R-MWYQKD02",
      "color_code": "3R6",
      "model": "Innova Zenix",
      "variant": "2.0 Q A/T",
      "color": "HITAM METALIK",
      "police_number": "V+096*XXP",
      "actual_mileage": 15000
    },
    "service_booking": {
      "booking_ID": "53A91DC8-04A0-4E8A-8B06-F8D04B3C941C",
      "booking_number": "AUT001301-03-20250327-nEZ",
      "created_datetime": 1709096400,
      "service_location": "DEALER",
      "service_category": "PERIODIC_MAINTENANCE",
      "slot_datetime_start": 1708398000,
      "slot_datetime_end": 1709085600,
      "service_sequence": 1,
      "outlet_ID": "AST001329",
      "outlet_name": "Astrido Toyota Bitung",
      "carrier_name": "RTsv vyOf",
      "carrier_phone_number": "TCLyY7ki67on",
      "vehicle_problem": "Wiper macet tidak bisa bergerak dan mesin tidak mulus",
      "booking_source": "MTOYOTA",
      "job_description": "kondisi mesin",
      "stall_or_squad_number": "A1"
    },
    "job": [
      {
        "job_ID": "08880CE605",
        "service_type": "RTJ",
        "job_name": "BATTERY_CHANGE",
        "labor_est_price": 500000
      },
      {
        "job_ID": "08880CE605",
        "service_type": "RTJ",
        "job_name": "OIL_CHANGE",
        "labor_est_price": 500000
      }
    ],
    "part": [
      {
        "part_type": "PACKAGE",
        "part_number": "P55B00KA0F",
        "part_name": "Lux Package",
        "part_quantity": 1,
        "package_parts": [
          {
            "part_number": "64102TA560",
            "part_name": "Black Outer Mirror Ornament"
          },
          {
            "part_number": "64102TA564",
            "part_name": "Side Body Moulding"
          }
        ],
        "part_size": "",
        "part_color": "",
        "part_est_price": 13690000,
        "flag_part_need_down_payment": true
      },
      {
        "part_type": "MERCHANDISE",
        "part_number": "64102TA560",
        "part_name": "Jaket",
        "part_quantity": 1,
        "package_parts": [],
        "part_size": "small",
        "part_color": "Black",
        "part_est_price": 200000,
        "flag_part_need_down_payment": false
      },
      {
        "part_type": "ACCESSORIES",
        "part_number": "64102TA560",
        "part_name": "White Outer Mirror Ornament",
        "part_quantity": 1,
        "package_parts": [],
        "part_size": "",
        "part_color": "Black",
        "part_est_price": 40000,
        "flag_part_need_down_payment": false
      }
    ],
    "estimation": {
      "service_estimated_delivery_datetime": 1709085600
    },
    "working_order": {
      "reception_datetime": 1709085600,
      "wo_number": "20302/SWO/25/03/00004",
      "wo_created_datetime": 1709085600,
      "bstk_created_datetime": 1709085600,
      "service_preference": "LEFT"
    },
    "production": {
      "clock_on_datetime": 1709085600,
      "technical_complete_datetime": 1709085600,
      "service_pause_flag": true,
      "service_pause_reason": "WAITING_FOR_INSURANCE_CONFIRMATION"
    },
    "payment": [
      {
        "invoice_number": "INV-000-0001",
        "payment_by": "DEALER",
        "payment_stage": [
          {
            "payment_status": "WAITING_FOR_BILLING",
            "status_datetime": 1709085600
          },
          {
            "payment_status": "INVOICE_RELEASED",
            "status_datetime": 1709085600
          },
          {
            "payment_status": "PAYMENT_COMPLETED",
            "status_datetime": 1709085600
          }
        ]
      },
      {
        "invoice_number": "INV-000-0002",
        "payment_by": "DEALER",
        "payment_stage": [
          {
            "payment_status": "WAITING_FOR_BILLING",
            "status_datetime": 1709085600
          },
          {
            "payment_status": "INVOICE_RELEASED",
            "status_datetime": 1709085600
          },
          {
            "payment_status": "PAYMENT_COMPLETED",
            "status_datetime": 1709085600
          }
        ]
      }
    ],
    "delivery": {
      "service_delivery_handover_datetime": 1709085600,
      "service_vehicle_received_by": "RTsv vyOf",
      "service_vehicle_receiver_phone_number": "TCLyY7ki67on"
    },
    "psfu": {
      "psfu_owner_flag": true
    }
  }
}

	`
	res := jsonParser.Get([]byte(grjson), "data.one_account")           // object
	res1 := jsonParser.Get([]byte(grjson), "process")                   // "service_progress_update"
	res2 := jsonParser.Get([]byte(grjson), "actual_mileage")            // undefined
	res3 := jsonParser.Get([]byte(grjson), "data.one_account.email2")   // null
	res4 := jsonParser.Get([]byte(grjson), "data.payment[].payment_by") // "DEALER"
	res5 := jsonParser.Get([]byte(grjson), "data.payment")              // Array of object

	fmt.Println(string(res.Data))
	fmt.Println(string(res1.Data))
	fmt.Println(string(res2.Data))
	fmt.Println(string(res3.Data))
	fmt.Println(string(res4.Data))
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
