package snmp_handler

import (
	//"flag"
	"fmt"
	"os"
	//"path/filepath"
	"time"

	"github.com/soniah/gosnmp"

	model "../models"
)



type SnmpMessage struct{
	Oid string
}

type ResponseSnmpItems struct{
	Items []SnmpMessage
}


func (sn *ResponseSnmpItems) collectValues(pdu gosnmp.SnmpPDU)  error {

	item := SnmpMessage{ Oid: pdu.Name}

	sn.Items = append(sn.Items, item)

	return nil

}


func (sn *ResponseSnmpItems) printValues()  {

	for i , item := range sn.Items {
		fmt.Println("i:", i, "oid:", item.Oid)
	}

}


func SnmpExecute(sendParams model.SnmpSendParams) {

	//flag.Usage = func() {
	//	fmt.Printf("Usage:\n")
	//	fmt.Printf("   %s [-community=<community>] host [oid]\n", filepath.Base(os.Args[0]))
	//	fmt.Printf("     host      - the host to walk/scan\n")
	//	fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
	//	flag.PrintDefaults()
	//}
	//
	//var community string
	//
	//flag.StringVar(&community, "community", "public", "the community string for device")
	//
	//flag.Parse()

	//if len(flag.Args()) < 1 {
	//	flag.Usage()
	//	os.Exit(1)
	//}
	//
	//target := flag.Args()[0]
	//
	//var oid string
	//
	//if len(flag.Args()) > 1 {
	//	oid = flag.Args()[1]
	//}

	// .1.3.6.1.4.1.119.2.3.69.501.7.10.1.8.1


	//target := "190.169.1.5"
	//// oid := ".1.3.6.1.4.1"
	//// oid := ".1.3.6.1.4.1.119.2.3.69.501.7"
	//oid := ".1.3.6.1.4.1.119.2.3.69.501.7.1.1"

	oid := sendParams.Oid
	gosnmp.Default.Target    = sendParams.Ip
	gosnmp.Default.Community = sendParams.Community
	gosnmp.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}

	defer gosnmp.Default.Conn.Close()

	sn := ResponseSnmpItems{}

	err = gosnmp.Default.BulkWalk(oid, sn.collectValues)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	sn.printValues()
	// fmt.Println(r.RespItems)
}