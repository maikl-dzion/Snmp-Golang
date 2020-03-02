package main

import (
	// model "../internal/models"
	// sn "../internal/snmp_handler"
	rmq "../internal/amqp_handler"
)


func main() {


	//param := model.SnmpSendParams{
	//	"190.169.1.5",
	//	".1.3.6.1.4.1.119.2.3.69.501.7.1.1",
	//	"public",
	//	"161",
	//	0,
	//}
	//
	//sn.BulkRequestRun(param)
	//
	//
	//param2 := model.SnmpSendParams{
	//	"190.169.1.5",
	//	".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.4.17",
	//	"public",
	//	"161",
	//	0,
	//}
	//
	//sn.GetRequestRun(param2)


	rmq.RecevieMessagesFromQueue()


}
