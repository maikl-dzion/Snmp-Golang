package snmp_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	model "github.com/maikl-dzion/Snmp-Golang/internal/models"
	snmp "github.com/soniah/gosnmp"
)

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
}

func (items *SnmpResultItems) CollectValues(pdu snmp.SnmpPDU) error {

	item := FormSnmpResultItem(pdu)
	item.DeviceId = items.DeviceId
	items.Items = append(items.Items, item)

	return nil

}

func (sn *SnmpResultItems) PrintValues() {
	for i, item := range sn.Items {
		fmt.Println("Ch::", i,
			"Oid::", item.Oid,
			"Ip::", item.Ip,
			"ValInt::", item.ValueInt,
			"ValStr::", item.ValueStr,
			"DateType::", item.DataType,
			"DeviceId::", item.DeviceId,
		)
	}
}

func SnmpBulkExecute(params model.SnmpSendParams) (SnmpResultItems, error) {

	//oid := params.Oid
	oid := ".1.3.6.1.4.1.119.2.3.69.501.7"
	snmp.Default.Target = params.Ip
	snmp.Default.Community = params.Community
	snmp.Default.Timeout = time.Duration(10 * time.Second)
	err := snmp.Default.Connect()

	if err != nil {
		fmt.Printf("Snmp Connect error: %v %s", err)
		os.Exit(1)
	}

	defer snmp.Default.Conn.Close()

	snmpResult := SnmpResultItems{}
	snmpResult.DeviceId = params.DeviceId

	err = snmp.Default.BulkWalk(oid, snmpResult.CollectValues)
	if err != nil {
		fmt.Printf("Walk Error: %v == %s\n", err, oid)
		os.Exit(1)
	}

	return snmpResult, err
}

func FormSnmpResultItem(pdu snmp.SnmpPDU) SnmpResultMessage {

	var valueStr string = ""
	var valueInt int64 = 0
	dataType := pdu.Type.String()
	oid := pdu.Name
	ip := snmp.Default.Target
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

func BulkRequestRun(params model.SnmpSendParams) (SnmpResultItems, error) {

	resultItems, err := SnmpBulkExecute(params)
	if err != nil {
		panic("Snmp Send Error")
	}


	_len := len(resultItems.Items)
	fmt.Println("Len snmp items", _len)
	// resultItems.PrintValues()

	return resultItems, nil
}

func MakeJsonMultiRequest(apiUrl string, messages []SnmpResultMessage) error {

	bytesRepresentation, err := json.Marshal(messages)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(apiUrl, "application/json",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var _result map[string]interface{}

	sendError := json.NewDecoder(resp.Body).Decode(&_result)

	//log.Println(sendError)
	//log.Println(_result)

	return sendError

}

func SnmpBulkRequestSend(params model.SnmpSendParams, saveApiUrl string) error {

	snmpItems, err := BulkRequestRun(params)

	if err != nil {
		panic("Snmp Send Error")
	}

	sendError := MakeJsonMultiRequest(saveApiUrl, snmpItems.Items)

	fmt.Println(sendError)

	return sendError

}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////// GET REQUEST ////////////////////////////////////////

func SnmpGetExecute(send model.SnmpSendParams) (SnmpResultItems, error) {

	port, _ := strconv.ParseUint(send.Port, 10, 16)

	params := &snmp.GoSNMP{
		Target:    send.Ip,
		Port:      uint16(port),
		Community: send.Community,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		// Logger:    log.New(os.Stdout, "", 0),
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}

	defer params.Conn.Close()

	oids := []string{send.Oid}
	result, errRequest := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if errRequest != nil {
		log.Fatalf("Get() err: %v", errRequest)
	}

	sn := SnmpResultItems{}

	for _, variable := range result.Variables {
		message := FormSnmpResultItem(variable)
		sn.Items = append(sn.Items, message)
	}

	return sn, err
}

func GetRequestRun(params model.SnmpSendParams) {

	items, errGet := SnmpGetExecute(params)

	if errGet != nil {
		panic("Snmp Send Error")
	}

	items.PrintValues()

}
