package main

import (
	"fmt"
	"log"

	snmp "github.com/soniah/gosnmp"
)

func main() {

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	snmp.Default.Target = "190.169.1.8"
	err := snmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}

	defer snmp.Default.Conn.Close()

	oids := []string{".1.3.6.1.4.1.119.2.3.69.5.1.1.1.3.1"}

	// result, err2 := snmp.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	result, err2 := snmp.Default.GetBulk(oids, 0, 10)

	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for i, variable := range result.Variables {

		fmt.Printf("%d: oid: %s ", i, variable.Name)

		switch variable.Type {
		case snmp.OctetString:
			fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
		default:
			fmt.Printf("number: %d\n", snmp.ToBigInt(variable.Value))
		}
	}

}