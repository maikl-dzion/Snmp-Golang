package model

type ResponseMessage struct {
	Oid      string `json:"oid"`
	Value    string `json:"value"`
	Ip       string `json:"ip"`
	DeviceId string `json:"device_id"`
}

type ResponseJsonItems []ResponseMessage


//type SnmpResponseMessage struct {
//	Oid      string `json:"oid"`
//	ValueInt string `json:"value_int"`
//	ValueStr string `json:"value_str"`
//	Ip       string `json:"ip"`
//	DeviceId string `json:"device_id"`
//}
//
//
//
//type SnmpResponseMessagesList struct {
//	Items []SnmpResponseMessage
//}

//func NewSnmpResponseMessagesList(items []SnmpResponseMessage) *SnmpResponseMessagesList {
//	return &SnmpResponseMessagesList{Items: items}
//}

//func (resp *ResponseJsonItems) SetJsonItem(ip, oid, value, deviceId, dateType string)  (error) {
//
//	resp.Items = append(resp.Items, ResponseMessage{Oid: oid,Ip : ip,Value:value,DeviceId: deviceId})
//
//	return nil
//
//}

const QUEUE_NAME   = "SNMP_QUEUE"
const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"
const SAVE_API_URL = "http://172.16.16.235:8080/data/multi_save"

type SnmpSendParams struct {
	Ip string
	Oid string
	Community string
	Port string
	SelCount int
	DeviceId string
}