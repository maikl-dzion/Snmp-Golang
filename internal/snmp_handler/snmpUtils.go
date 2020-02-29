package snmp_handler

import (
	"fmt"
	"math/big"
	"os"
	"time"
	"github.com/soniah/gosnmp"
	model "../models"
)


type SnmpResultMessage struct{
	Oid string
	ValueStr string
	ValueInt *big.Int
	Type gosnmp.Asn1BER
}

type ResponseSnmpItems struct{
	Items []SnmpResultMessage
}


func (sn *ResponseSnmpItems) CollectValues(pdu gosnmp.SnmpPDU)  error {

	var valueStr string = ""
	var valueInt *big.Int

	switch pdu.Type {
		case gosnmp.OctetString:
			valueStr = string(pdu.Value.([]byte))
		default:
			valueInt = gosnmp.ToBigInt(pdu.Value)
	}


	item := SnmpResultMessage{
		pdu.Name,
		valueStr,
		valueInt,
		pdu.Type,
	}

	sn.Items = append(sn.Items, item)

	return nil

}


func (sn *ResponseSnmpItems) PrintValues()  {

	for i , item := range sn.Items {
		fmt.Println("I:", i,
			           "Oid:", item.Oid,
			           "Value str:", item.ValueStr,
			           "Value int:", item.ValueInt,
			           "Type:", item.Type,
		            )
	}

}


func SnmpBulkExecute(sendParams model.SnmpSendParams) (ResponseSnmpItems, error) {

	oid := sendParams.Oid
	gosnmp.Default.Target    = sendParams.Ip
	gosnmp.Default.Community = sendParams.Community
	gosnmp.Default.Timeout   = time.Duration(10 * time.Second) // Timeout better suited to walking
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}

	defer gosnmp.Default.Conn.Close()

	sn := ResponseSnmpItems{}

	err = gosnmp.Default.BulkWalk(oid, sn.CollectValues)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	return sn, err
}
