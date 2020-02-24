package model


type ResponseMessage struct {
	Oid      string `json:"oid"`
	Value    string `json:"value"`
	Ip       string `json:"ip"`
	DeviceId string `json:"device_id"`
}


const QUEUE_NAME   = "SNMP_QUEUE"
const SAVE_API_URL = "http://172.16.16.235:8080/data/save"
const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"
