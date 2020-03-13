package main

import (
	"Snmp-Golang/example/snmp"
	"fmt"
	//"net"
	//"time"
)

func main() {

	ip := "190.169.1.5"
	// ip := "192.168.2.184"

	sn, err := snmp.New(ip, "public")


	fmt.Println(err)
	fmt.Println(sn)

	//ipAddr := "192.168.2.184"
	//oid := "1.3.6.1.2.1.1.1.0"
	//
	//var resp interface{}
	//
	//if err := snmp.Get(ipAddr, "public", oid, &resp); err != nil {
	//	fmt.Println("Error:", err, "OID:", oid)
	//}
	//
	//have := fmt.Sprintf("%T", resp)
	//fmt.Println(resp, have)

}
