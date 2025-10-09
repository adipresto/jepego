package configuration

var JsonExample string = `
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
        "payment_by": "DEALERB",
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
