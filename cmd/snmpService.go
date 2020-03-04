package main

import (
	// "fmt"
	// snmp_serv "../internal/snmp_handler"
	rabbitmq "../internal/amqp_handler"
	model "../internal/models"
)


func main() {


	queueName := model.QUEUE_NAME
	amqpUrl   := model.AMQP_API_URL

	rabbitmq.RecevieMessagesListFromQueue(amqpUrl, queueName)




	//param := model.SnmpSendParams{
	//	"190.169.1.5",
	//	".1.3.6.1.4.1.119.2.3.69.501.7.1.1",
	//	"public",
	//	"161",
	//	0,
	//}
	//
	//snmp_serv.BulkRequestRun(param)

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



	// saveApiUrl := model.SAVE_API_URL

	//messages := []model.ResponseMessage{}
	//
	//messages = append(messages, model.ResponseMessage{
	//	                         Oid:".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17",
	//	                         Ip: "192.168.10.12",
	//	                         Value:"MG-It",
	//	                         DeviceId: "234",
	//                            })
	//
	//
	//rabbitmq.SendCurlExec(saveUrl, messages, "10")


	// saveApiUrl = "http://localhost/snmp-url/"
	// saveApiUrl = "http://localhost/snmp-url/"

	// messages := model.ResponseJsonItems{}

	//messages.SetJsonItem("192.168.10.12", "Тест 100", ".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17", "234", "INTEGER")
	//messages.SetJsonItem("192.168.10.12", "Тест 200", ".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17", "235", "INTEGER")
	//messages.SetJsonItem("192.168.10.12", "Тест 300", ".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17", "236", "INTEGER")
	//
	//fmt.Println(messages.Items)

	//messages = append(messages, model.ResponseMessage{
	//								 Oid:".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17",
	//								 Ip: "192.168.10.12 cvbb",
	//								 Value:"Тест 100",
	//								 DeviceId: "234",
	//	                        })
	//
	//
	//messages = append(messages,model.ResponseMessage{
	//								Oid:".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17",
	//								Ip: "192.168.10.12 dff",
	//								Value:"Тест 20456",
	//								DeviceId: "234",
	//                           })
	//
	//fmt.Println(messages)
	//
	//
	//rabbitmq.MakeJsonRequest(saveApiUrl, messages)


}



//func (r ResponseJsonItems) FomattedJson(d model.ResponseMessage) {
//	r = append(r, []byte(d))
//}