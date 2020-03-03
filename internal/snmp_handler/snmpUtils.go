package snmp_handler

import (
	model "../models"
	"fmt"
	snmp "github.com/soniah/gosnmp"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)


type SnmpResultMessage struct{
	Oid string
	ValueStr string
	ValueInt *big.Int
	Type snmp.Asn1BER
}

type SnmpResultItems struct{
	Items []SnmpResultMessage
}


func (sn *SnmpResultItems) CollectValues(pdu snmp.SnmpPDU)  error {

	//var valueStr string = ""
	//var valueInt *big.Int
	//
	//switch pdu.Type {
	//	case snmp.OctetString:
	//		valueStr = string(pdu.Value.([]byte))
	//	default:
	//		valueInt = snmp.ToBigInt(pdu.Value)
	//}
	//
	//
	//item := SnmpResultMessage{
	//	pdu.Name,
	//	valueStr,
	//	valueInt,
	//	pdu.Type,
	//}

	message := FormSnmpItem(pdu)

	sn.Items = append(sn.Items, message)

	return nil

}


func (sn *SnmpResultItems) PrintValues()  {

	for i , item := range sn.Items {
		fmt.Println("I:", i,
			           "Oid:", item.Oid,
			           "Value str:", item.ValueStr,
			           "Value int:", item.ValueInt,
			           "Type:", item.Type,
		            )
	}

}


func SnmpBulkExecute(sendParams model.SnmpSendParams) (SnmpResultItems, error) {

	oid := sendParams.Oid
	snmp.Default.Target    = sendParams.Ip
	snmp.Default.Community = sendParams.Community
	snmp.Default.Timeout   = time.Duration(10 * time.Second) // Timeout better suited to walking
	err := snmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}

	defer snmp.Default.Conn.Close()

	sn := SnmpResultItems{}

	err = snmp.Default.BulkWalk(oid, sn.CollectValues)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	return sn, err
}




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
	result, errRequest := params.Get(oids)  // Get() accepts up to g.MAX_OIDS
	if errRequest != nil {
		log.Fatalf("Get() err: %v", errRequest)
	}

	sn := SnmpResultItems{}

	for _, variable := range result.Variables {
		message := FormSnmpItem(variable)
		sn.Items = append(sn.Items, message)
	}

	return sn, err
}



func FormSnmpItem(pdu snmp.SnmpPDU) SnmpResultMessage {

	var valueStr string = ""
	var valueInt *big.Int

	switch pdu.Type {
	case snmp.OctetString:
		valueStr = string(pdu.Value.([]byte))
	default:
		valueInt = snmp.ToBigInt(pdu.Value)
	}


	item := SnmpResultMessage{
		pdu.Name,
		valueStr,
		valueInt,
		pdu.Type,
	}

	return item
}


func BulkRequestRun(sendParams model.SnmpSendParams) {

	resultItems, err := SnmpBulkExecute(sendParams)
	if err != nil {
		panic("Snmp Send Error")
	}

	resultItems.PrintValues()
}


func GetRequestRun(sendParams model.SnmpSendParams) {

	items, errGet := SnmpGetExecute(sendParams)

	if errGet != nil {
		panic("Snmp Send Error")
	}

	items.PrintValues()

}



