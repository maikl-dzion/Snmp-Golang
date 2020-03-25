package main

import (
	common_serv "Snmp-Golang/internal/common_services"
	model "Snmp-Golang/internal/models"
	"fmt"
	"sync"
	"time"
)


func main() {

	// snmpbulkget -v2c -Cn0 -Cr5 -c public 190.169.1.5 .1.3.6.1.4.1.119.2.3.69.501.7
	// model.LogSave("test1 test2 test3", "")

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

	// common_serv.GetQueueMessagesInizialize(commonParams, sendParams)

	var wg sync.WaitGroup

	//messChannel := make(chan model.SnmpSendParams)
	//
	//ch := 0
	//stateCh := 0
	//limit := 100
	//
	//for {
	//
	//	//wg.Add(2)
	//
    //    ch++
	//	stateCh++
	//
	//	_printCount(ch)
	//
	//	go common_serv.GetTasksFromQueue(commonParams, sendParams, messChannel)
	//	go common_serv.SendRequestToDevice(messChannel, commonParams)
	//
    //    if stateCh > limit {
	//		stateCh = 0
	//		time.Sleep(1 * time.Second)
	//	}
	//
	//
	//}
	//
	//fmt.Println("End Request:")

	// messChannel := make(chan model.SnmpSendParams, 20)

	ch      := 0
	stateCh := 0
	limit   := 30

	// for {
	for i := 0; i < 5000; i++ {

		wg.Add(2)

		ch++
		stateCh++

		messChannel := make(chan model.SnmpSendParams)

		// go common_serv.GetTasksFromQueueNew(commonParams, sendParams, &wg)

		go common_serv.GetTasksFromQueue(messChannel, commonParams, sendParams, &wg, ch)
		go common_serv.SendRequestToDevice(messChannel, commonParams, &wg, ch)

		// _printCh(ch)

		if stateCh > limit {
			stateCh = 0
			time.Sleep(1 * time.Second)
		}

	}


	fmt.Println("Befory Ok:")

	for i := 0; i < 500; i++ {
		time.Sleep(5 * time.Second)
	}

	fmt.Println("End Request Ok:")

	wg.Wait()


}


func _printCh(ch int) {
	fmt.Println("Кол. полученных заданий", ch)
}
