package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	// "strconv"
	"time"

	"github.com/soniah/gosnmp"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s [-community=<community>] host [oid]\n", filepath.Base(os.Args[0]))
		fmt.Printf("     host      - the host to walk/scan\n")
		fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
		flag.PrintDefaults()
	}

	var community string
	flag.StringVar(&community, "community", "public", "the community string for device")

	flag.Parse()


	target := "192.168.2.184"
	oid    := ".1.3.6.1"
	oid    = ".1.3.6.1.2.1.4.30.1.5.2"

	gosnmp.Default.Target    = target
	gosnmp.Default.Community = community
	gosnmp.Default.Timeout = time.Duration(10 * time.Second)
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer gosnmp.Default.Conn.Close()

	err = gosnmp.Default.BulkWalk(oid, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
}

func printValue(pdu gosnmp.SnmpPDU) error {

	fmt.Printf("%s = ", pdu.Name)
	fmt.Println("TypeName:", pdu.Type)
	fmt.Println("Ip:", gosnmp.Default.Target)
	// fmt.Println("Logger:", pdu.Logger.Printf)

	var value    string = ""
	var typeName string = pdu.Type.String()

	switch pdu.Type {
		case gosnmp.OctetString:
			value = string(pdu.Value.([]byte))
		default:
			value = gosnmp.ToBigInt(pdu.Value).String()
			// fmt.Printf("myType %s: %s\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}

	fmt.Println("Value:", value, "Type:", typeName)

	//switch pdu.Type {
	//case gosnmp.OctetString:
	//	b := pdu.Value.([]byte)
	//	fmt.Printf("STRING: %s\n", string(b))
	//default:
	//	fmt.Printf("INT %d\n", gosnmp.ToBigInt(pdu.Value))
	//}

	return nil

}