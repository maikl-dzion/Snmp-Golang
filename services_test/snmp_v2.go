package main

import (
	model "Snmp-Golang/internal/models"
	snmp_serv "Snmp-Golang/internal/snmp_handler"
	"fmt"
)

func main() {

	realQuery  := false
	saveApiUrl := model.SAVE_API_URL
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

	for i := 0; i < 5000; i++ {

		resp , err := snmp_serv.WalkGetAllRun(param)

		if err != nil {
			fmt.Println("Error: WalkGetAllRun function", err)
		}

		// snmp_serv.MakeJsonMultiRequest(saveApiUrl, resp.Items)

		fmt.Println("Items : ", resp.Items)
		fmt.Println("Amqp Message Number : ", i)
		model.DatetimePrint()
	}


}


//func WalkAllRun(param model.SnmpSendParams) {
//
//
//	g.Default.Target    = param.Ip
//	g.Default.Community = param.Community
//	oid := param.Oid
//
//	err := g.Default.Connect()
//	if err != nil {
//		fmt.Printf("Snmp Connect Error: %v", err)
//	}
//
//	defer g.Default.Conn.Close()
//
//
//	snmpResults, err := g.Default.WalkAll(oid)
//
//	if err != nil {
//		fmt.Printf("Snmp Walk Error: %v", err)
//	}
//
//	ch := 0
//	for _, v := range snmpResults {
//
//		ch++
//		model.LogPrint(ch, "Number->")
//
//		model.DatetimePrint()
//
//		fmt.Printf("oid: %s, value: ", v.Name)
//		switch v.Type{
//		case g.OctetString:
//			fmt.Printf("%s\n", string(v.Value.([]byte)))
//		default:
//			fmt.Printf("%d\n", g.ToBigInt(v.Value))
//		}
//	}
//
//}