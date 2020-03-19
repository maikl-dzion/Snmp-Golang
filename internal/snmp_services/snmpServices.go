package snmp_services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	model "Snmp-Golang/internal/models"
	snmp "github.com/soniah/gosnmp"
)


////////////////////////////////////////
/**************************************
   СТРУКТУРА ДЛЯ СОХРАНЕНИЯ РЕЗУЛЬТАТОВ
           SNMP - ЗАПРОСА
*************************************/
type SnmpResultMessage struct {
	Oid      string `json:"oid"`
	Ip       string `json:"ip"`
	ValueInt int64  `json:"value_int"`
	ValueStr string `json:"value_str"`
	DeviceId string `json:"device_id"`
	DataType string `json:"data_type"`
}

type SnmpResultItems struct {
	Items    []SnmpResultMessage
	DeviceId string
	Ip string
}

func (r *SnmpResultItems) CollectValues(pdu snmp.SnmpPDU) error {

	item := SnmpResultDataConvert(pdu)
	item.DeviceId = r.DeviceId
	item.Ip = r.Ip
	r.Items = append(r.Items, item)

	return nil
}

func (items *SnmpResultItems) PrintValues() {
	for i, item := range items.Items {
		fmt.Println("[Ch]=", i,
			"[Oid]=", item.Oid,
			"[Ip]=", item.Ip,
			"[ValInt]=", item.ValueInt,
			"[ValStr]=", item.ValueStr,
			"[DateType]=", item.DataType,
			"[DeviceId]=", item.DeviceId)
	}
}

//**************************************
////////////////////////////////////////


func SnmpConfigInit(s model.SnmpSendParams) (*snmp.GoSNMP, error) {
	client := &snmp.GoSNMP{
		Target:    s.Ip,
		Port:      uint16(161),
		Community: s.Community,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(10) * time.Second,
		//Logger:    log.New(os.Stdout, "", 0),
		//Retries:   s.Retries,
	}
	err := client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return client, nil
}

func SnmpManagerStart(conf model.SnmpSendParams, fType string) (SnmpResultItems, error) {

	now := time.Now()

	client, err := SnmpConfigInit(conf)
	if err != nil {
		log.Fatal("Snmp Config Is Nil")
	}
	defer client.Conn.Close()

	snmpResult := SnmpResultItems{}
	snmpResult.DeviceId = conf.DeviceId
	snmpResult.Ip = conf.Ip
	oid := conf.Oid

	switch fType {
	case "GET":
		oids := []string{oid}
		pkt, err := client.Get(oids)
		if err != nil {
			return snmpResult, err
		}
		for _, pdu := range pkt.Variables {
			snmpResult.CollectValues(pdu)
		}

	case "BULK":
		err := client.BulkWalk(oid, snmpResult.CollectValues)
		if err != nil {
			return snmpResult, err
		}

	case "BULK_WALK_ALL":
		pkt, err := client.BulkWalkAll(oid)
		if err != nil {
			return snmpResult, err
		}
		for _, pdu := range pkt {
			snmpResult.CollectValues(pdu)
		}

	default:
		err := client.BulkWalk(oid, snmpResult.CollectValues)
		if err != nil {
			return snmpResult, err
		}
	}

	ch := len(snmpResult.Items)
	model.LogPrint(ch, "Snmp Result Oid Count:")
	fmt.Println("Now", now)

	return snmpResult, nil
}



func MakeJsonMultiRequest(apiUrl string, messages []SnmpResultMessage) (map[string]interface{}, error) {

	var jsonResultLog map[string]interface{}

	bytesRepresentation, err := json.Marshal(messages)
	if err != nil {
		model.WarnOnError(err, "MakeJsonMultiRequest::Json Marshal ERROR:")
		return jsonResultLog, err
	}

	resp, err := http.Post(apiUrl, "application/json",
		                   bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		model.WarnOnError(err, "MakeJsonMultiRequest::Http Post ERROR:")
		return jsonResultLog, err
	}

	jsonSaveError := json.NewDecoder(resp.Body).Decode(&jsonResultLog)

	return jsonResultLog, jsonSaveError

}

func SnmpResultDataConvert(pdu snmp.SnmpPDU) SnmpResultMessage {

	var valueStr string = ""
	var valueInt int64 = 0
	dataType := pdu.Type.String()
	oid := pdu.Name
	ip  := snmp.Default.Target
	deviceId := ""

	switch pdu.Type {
	case snmp.OctetString:
		valueStr = string(pdu.Value.([]byte))
	default:
		valueInt = snmp.ToBigInt(pdu.Value).Int64()
	}

	item := SnmpResultMessage{
		Ip:       ip,
		Oid:      oid,
		ValueInt: valueInt,
		ValueStr: valueStr,
		DataType: dataType,
		DeviceId: deviceId,
	}

	return item
}



func SnmpExceptionHandler(msg model.SnmpSendParams) error {
	fmt.Println(msg)
	return nil
}

