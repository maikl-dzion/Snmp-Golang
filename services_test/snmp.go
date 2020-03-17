package main

import (
	model "Snmp-Golang/internal/models"
	snmp_service "Snmp-Golang/internal/snmp_handler"
	"fmt"
)


func main() {

	saveApiUrl := model.SAVE_API_URL
	realQuery  := false

	fmt.Println(saveApiUrl)

	param := model.SnmpSendParams{
		Ip: "192.168.2.184",
		Oid:".1.3.6.1",
		// Oid:       ".1.3.6.1.2.1.1.9.1",
		Community: "public",
		Port:      "161",
		DeviceId:  "",
		SelCount:  0,
	}

	if realQuery {

		param.Ip = "190.169.1.5"
		param.Oid = ".1.3.6.1.4.1.119.2.3.69.501.7"
		// param.Oid = ".1.3.6.1.4.1.119.2.3.69"

	}


	forever := make(chan bool)

	for i := 0; i < 100; i++ {

		go func() {

			forever <-true

			snmpItems, err := snmp_service.BulkRequestRun(param)

			if err != nil {
				panic("Snmp Send Error")
			}

			fmt.Println(snmpItems.Items)

		}()

	}

	<-forever

	// forever := make(chan bool)

	// go func(param model.SnmpSendParams, saveApiUrl string) {

	   // for i := 0; i < 5000; i++ {
	   //for {
	   //
		//	//res, err := snmp_service.BulkRequestRun(param)
		//	//fmt.Println(err)
		//	//res.PrintValues()
	   //
		//	// forever := make(chan bool)
	   //
		//	//go func(param model.SnmpSendParams, saveApiUrl string) {
	   //
		//		go snmp_service.SnmpBulkRequestSend(param, saveApiUrl)
		//        time.Sleep(4 * time.Second)
	   //
		//	// }(param, saveApiUrl)
	   //
		//	// <-forever
		//}

	// }(param, saveApiUrl)


	// <-forever


	//waitGroup := sync.WaitGroup{}
	//waitGroup.Add(5000)

	//forever := make(chan bool)
	//
	////go func() {
	//
	//	for i := 0; i < 5000; i++ {
	//
	//
	//			snmp_service.SnmpBulkRequestSend(param, saveApiUrl)
	//
	//			//sendError := snmp_service.SnmpBulkRequestSend(param, saveApiUrl)
	//			//fmt.Println("Error:", sendError)
	//			fmt.Println("Number :", i)
	//
	//			//waitGroup.Done()
	//
	//	}
	//
	////}()
	//
	//<-forever
	//waitGroup.Wait()


	urls := []string{}

	for i := 0; i < 1000; i++ {
		urls = append(urls, "546353")
	}

	//jsonResponses := make(chan string)
	//
	//var wg sync.WaitGroup
	//
	//wg.Add(len(urls))
	//
	//for _, url := range urls {
	//	go func(url string) {
	//		defer wg.Done()
	//		err := snmp_service.SnmpBulkRequestSend(param, saveApiUrl)
	//
	//		if err != nil {
	//			log.Fatal(err)
	//		} else {
	//			jsonResponses <- string("tryryy")
	//		}
	//
	//	}(url)
	//}
	//
	//go func() {
	//	for response := range jsonResponses {
	//		fmt.Println(response)
	//	}
	//}()
	//
	//wg.Wait()



	//var wg sync.WaitGroup
	//// Create bufferized channel with size 5
	//// results := make(chan error)
	//goroutines := make(chan struct{}, 5)
	//// Read data from input channel
	//for data := range urls {
	//	// 1 struct{}{} - 1 goroutine
	//	goroutines <- struct{}{}
	//	wg.Add(1)
	//	go func(goroutines <-chan struct{}, data int, wg *sync.WaitGroup) {
	//		// Process data
	//		snmp_service.SnmpBulkRequestSend(param, saveApiUrl)
	//		// Read from "goroutines" channel to free space for new goroutine
	//		<-goroutines
	//		wg.Done()
	//	}(goroutines, data, &wg)
	//}
	//
	//wg.Wait()
	//close(goroutines)



}

