package main

import (
	model "Snmp-Golang/internal/models"
	snmp_serv "Snmp-Golang/internal/snmp_handler"
	"fmt"
	"log"
)

func main() {

	realIpAddr := false
	saveApiUrl := model.SAVE_API_URL
	fmt.Println(saveApiUrl)

	param := model.SnmpSendParams{
		Ip: "192.168.2.184",
		// Oid:".1.3.6.1",
		Oid:".1.3.6.1.2.1.1.9.1",
		// Oid:"1.3.6.1.2.1.1.1.0",
		Community: "public",
		Port:      "161",
		DeviceId:  "Printer-Id-100",
		SelCount:  0,
	}

	if realIpAddr {
		param.Ip = "190.169.1.5"
		param.Oid = ".1.3.6.1.4.1.119.2.3.69.501.7"
		// param.Oid = ".1.3.6.1.4.1.119.2.3.69"
	}

	// funcType := "walk_all"
	// funcType := "bulk"
	// funcType := "bulk_walk_all"
	// fType := "BULK_WALK_ALL"
	fType := "BULK"

	//for i := 0; i < 3; i++ {
	//
	//	// resp , err := snmp_serv.SnmpWalkAllExecute(param)
	//	// resp , err := snmp_serv.SnmpBulkExecute(param)
	//	// resp , err := snmp_serv.SnmpGetExecute(param)
	//	response , err := snmp_serv.SnmpRequestRun(param, funcType)
	//	if err != nil {
	//		fmt.Println("Error: WalkGetAllRun function", err)
	//		os.Exit(1)
	//	}
	//
	//	snmp_serv.MakeJsonMultiRequest(saveApiUrl, response.Items)
	//	print_(response, i)
	//
	//}


	//for i := 0; i < 3; i++ {
	//	saveError := snmp_serv.SnmpStart(param, saveApiUrl, funcType)
	//	if saveError != nil {
	//		fmt.Println("Error: SnmpStart function", saveError)
	//	}
	//}

    r, err := snmp_serv.SnmpManagerStart(param, fType)
    if err != nil {
    	fmt.Println(err)
		log.Fatal("Snmp Request Func Error")
	}

    r.PrintValues()

}

func print_(r snmp_serv.SnmpResultItems, i int) {

	r.PrintValues()
	// fmt.Println("-Items:", r.Items)
	fmt.Println("-Queue Message Num:", i)
	model.DatetimePrint()
}

