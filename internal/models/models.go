package model

import (
	"fmt"
	"time"
)

//type ResponseMessage struct {
//	Oid      string `json:"oid"`
//	Value    string `json:"value"`
//	Ip       string `json:"ip"`
//	DeviceId string `json:"device_id"`
//}
//
//type ResponseJsonItems []ResponseMessage


const QUEUE_NAME   = "SNMP_QUEUE"
const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"
const SAVE_API_URL = "http://172.16.16.235:8080/data/multi_save"

type SnmpSendParams struct {
	Id string
	Ip string
	Oid string
	Community string
	Port string
	SelCount int
	DeviceId string
}


/////////////////////////////////////
// ---- Вспомогательные функции ----

func DatetimePrint() {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d__%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println("DateTime->", formatted)
}


func LogPrint(data interface{}, name string) {
	if name != ""{
		fmt.Println(name, data)
	} else {
		fmt.Println(data)
	}
}