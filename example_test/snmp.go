package main

import (
	"fmt"

	model "../internal/models"
	snmp_service "../internal/snmp_handler"
)


func main() {

	realQuery := true

	param := model.SnmpSendParams{
		Ip:"192.168.2.184",
		// Oid:".1.3.6.1",
		Oid:".1.3.6.1.2.1.1.9.1",
		Community:"public",
		Port:"161",
		DeviceId:"",
		SelCount:0,
	}

	if realQuery {

		param.Ip   = "190.169.1.5"
		// param.Oid = ".1.3.6.1.4.1.119.2.3.69.501.7"
		param.Oid = ".1.3.6.1.4"

	}

	res, err := snmp_service.BulkRequestRun(param)

	fmt.Println(err)
	res.PrintValues()
	// fmt.Println(res.Items)


	//snmp_serv.GetRequestRun(param)

}

