package main

import (
	model "../internal/models"
	sn "../internal/snmp_handler"
)


func main() {


	param := model.SnmpSendParams{
		"190.169.1.5",
		".1.3.6.1.4.1.119.2.3.69.501.7.1.1",
		"public",
		"161",
		0,
	}


	response, err := sn.SnmpBulkExecute(param)

	if err != nil {
		panic("Snmp Send Error")
	}


	response.PrintValues()

}
