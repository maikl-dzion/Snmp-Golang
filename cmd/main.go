package main

import (

	common_serv "Snmp-Golang/internal/common_services"
	model "Snmp-Golang/internal/models"
)


func main() {

	commonParams := model.GetCommonInitParam()

	sendParams := model.SnmpSendParams{
		Ip:  "192.168.2.184",
		Oid: ".1.3.6.1",
		// Oid:".1.3.6.1.2.1.1.9.1",
		Community: "public",
		Port:      "161",
		DeviceId:  "TestDeviceId-777-999",
		SelCount:  0,
	}

	common_serv.GetQueueMessagesInizialize(commonParams, sendParams)

}