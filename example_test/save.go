package main

import (
	mq "../internal/amqp_handler"
	model "../internal/models"
	snmp_service "../internal/snmp_handler"
)


func main() {

	saveApiUrl := model.SAVE_API_URL

	resultMessages := snmp_service.SnmpResultItems{}

	elem := snmp_service.SnmpResultMessage{
		Ip : "192.168.2.184",
		Oid: ".1.3.6.1.2.1.1.9.1.2.2",
		ValueInt:0,
		ValueStr:"PRINTER-HP-6789",
		DataType:"ObjectIdentifier",
		DeviceId:"Printer-6789-model",
	}
	resultMessages.Items = append(resultMessages.Items, elem)

	elem = snmp_service.SnmpResultMessage{
		Ip : "192.168.2.184",
		Oid: ".1.3.6.1.2.1.1.9.1.2.4",
		ValueInt:5789,
		ValueStr:"",
		DataType:"INTEGER",
		DeviceId:"Printer-6789-model",
	}
	resultMessages.Items = append(resultMessages.Items, elem)


	mq.MakeJsonRequest(saveApiUrl, resultMessages)

}
