package main

import (
	model "../internal/models"
	sn "../internal/snmp_handler"
	"fmt"
	"github.com/soniah/gosnmp"
)


func main() {


	param := model.SnmpSendParams{
		"190.169.1.5",
		".1.3.6.1.4.1.119.2.3.69.501.7.1.1",
		"public",
		"161",
		0,
	}


	sn.SnmpExecute(param)

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
	//
	////if len(flag.Args()) < 1 {
	////	flag.Usage()
	////	os.Exit(1)
	////}
	////
	////target := flag.Args()[0]
	////
	////var oid string
	////
	////if len(flag.Args()) > 1 {
	////	oid = flag.Args()[1]
	////}
	//
	//// .1.3.6.1.4.1.119.2.3.69.501.7.10.1.8.1
	//
	//
	//target := "190.169.1.5"
	//// oid := ".1.3.6.1.4.1"
	//// oid := ".1.3.6.1.4.1.119.2.3.69.501.7"
	//oid := ".1.3.6.1.4.1.119.2.3.69.501.7.1.1"
	//
	//gosnmp.Default.Target    = target
	//gosnmp.Default.Community = community
	//gosnmp.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
	//err := gosnmp.Default.Connect()
	//if err != nil {
	//	fmt.Printf("Connect err: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//defer gosnmp.Default.Conn.Close()
	//
	//sn := ResponseSnmpItems{}
	//
	//err = gosnmp.Default.BulkWalk(oid, sn.collectValues)
	//if err != nil {
	//	fmt.Printf("Walk Error: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//
	//sn.printValues()
	//// fmt.Println(r.RespItems)
}

func printValue(pdu gosnmp.SnmpPDU) error {

	fmt.Printf("%s = ", pdu.Name)

	// message := []model.ResponseMessage{}

	switch pdu.Type {
		case gosnmp.OctetString:
			b := pdu.Value.([]byte)
			fmt.Printf("STRING: %s\n", string(b))
		default:
			fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}

	return nil
}