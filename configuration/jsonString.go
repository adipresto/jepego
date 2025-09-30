package configuration

var JsonStrClean string = `

{

"bookingDateTime": "2025-01-23 10:00:00", 
"branchCode": "T002", 
"vehicle": { 
"policeRegNo": "R01AND", 
"vin": "EQUIMENT001",
"year": 2015,
"vehicleModel": "AVANZA" 

},

"customer": { 
"Name1": "1dxmcz", 
"Name2": "iP4ep8", 
"Phone1": "STincMrSl9VQ", 
"Email1": "rDHJzJii@4mg5B.2jF", 
"Address1": "Swptz1Tx sLi4O2M, 5tpZFK 2452u, xCj4gO7WJ IeAiE6G, dvpisuda TJ1wJb, pRkQ1gHkP"

},

"serviceOrderJobs": [ 

{

"jobCode": "1000", 
"jobDescription": "Servis Berkala 1.000 KM", 
"price": "0.0", 
"duration": 5.0 

},

{

"jobCode": "OIL",
"jobDescription": "Engine Oil",
"price": "383000.0",
"duration": 5.0

},

{

"jobCode": "OTHER",
"jobDescription": "Other Service",
"price": "0.0",
"duration": 1.0

}

],
"symptom": "ga ada",
"totalDuration": 11.0, 
"orderDateStartTime": "2025-01-23 10:00:00", 
"city": "Jakarta Utara",
"region": "DKI Jakarta",
"action": "GR" 
}

`

var JsonStr string = `
	{

"bookingDateTime": "2025-01-23 10:00:00", //mandatory
"branchCode": "T002", //mandatory
"vehicle": { //mandatory
"policeRegNo": "R01AND", //mandatory
"vin": "EQUIMENT001",
"year": 2015,
"vehicleModel": "AVANZA" //mandatory

},

"customer": { //mandatory
"Name1": "1dxmcz", //mandatory
"Name2": "iP4ep8", //mandatory
"Phone1": "STincMrSl9VQ", //mandatory
"Email1": "rDHJzJii@4mg5B.2jF", //mandatory
"Address1": "Swptz1Tx sLi4O2M, 5tpZFK 2452u, xCj4gO7WJ IeAiE6G, dvpisuda TJ1wJb, pRkQ1gHkP"

},

"serviceOrderJobs": [ //optional

{

"jobCode": "1000", //mandatory
"jobDescription": "Servis Berkala 1.000 KM", //mandatory
"price": "0.0", //mandatory
"duration": 5.0 //mandatory

},

{

"jobCode": "OIL",
"jobDescription": "Engine Oil",
"price": "383000.0",
"duration": 5.0

},

{

"jobCode": "OTHER",
"jobDescription": "Other Service",
"price": "0.0",
"duration": 1.0

}

],
"symptom": "ga ada",
"totalDuration": 11.0, //mandatory
"orderDateStartTime": "2025-01-23 10:00:00", //mandatory
"city": "Jakarta Utara",
"region": "DKI Jakarta",
"action": "GR" //mandatory
}
`
