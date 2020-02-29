package snmp_handler

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	st "strings"
	"time"

	model "../models"
	snmp "github.com/soniah/gosnmp"
)

func SnmpBulkExec(target string,
	           oid string,
	           selectCount int,
	           community string,
	           port string)  ([]model.ResponseMessage, bool, string) {


	if len(target) <= 0 {
		log.Fatalf("environment variable not set: GOSNMP_TARGET")
	}
	if len(port) <= 0 {
		log.Fatalf("environment variable not set: GOSNMP_PORT")
	}

	p, _ := strconv.ParseUint(port, 10, 16)


	params := &snmp.GoSNMP{
		Target:    target,
		Port:      uint16(p),
		Community: community,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		// Logger:    log.New(os.Stdout, "", 0),
	}

	err := params.Connect()

	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}

	defer params.Conn.Close()

	oids := []string{oid}

	fmt.Println(reflect.TypeOf(params))


	result, err2 := params.GetBulk(oids, 0, uint8(selectCount))

	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}


	messages := []model.ResponseMessage{}


	var endState bool = false
	var endOid string = ""



	for _ , variable := range result.Variables {

		name := variable.Name
		endOid = name

		// fmt.Println(variable.Value.())

		// fmt.Println(variable.Value)


		if !st.HasPrefix(name, oid) {
			endState = true
			return messages, endState, endOid
		}

		//fmt.Println(_state)
		var value string = ""

		switch variable.Type {
			case snmp.OctetString:
				value = string(variable.Value.([]byte))
				//fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			default:
				value = string(variable.Value.([]byte))
				//fmt.Printf("number: %d\n", snmp.ToBigInt(variable.Value))
		}

		message := model.ResponseMessage{}
		message.Oid   = name
		message.Ip    = target
		message.Value = value
		messages = append(messages, message)
	}


	if oid == ".1.3.6.1.4.1.119.2.3.69.501.7.1.2.1.6.23" {

		fmt.Println("oid:", oid)
		fmt.Println("oidEnd:", endOid)

		fmt.Println("message:", messages)
	}

	return messages, endState, endOid
}





func BulkTreeRecursive(ip string, oid string, selectCount int, community string, port string) ([]model.ResponseMessage){

	response, endState, endOid := SnmpBulkExec(ip, oid, selectCount, community, port)

	if !endState {
		resp := BulkTreeRecursive(ip, endOid, selectCount, community, port)
        // fmt.Println("oid:", endOid)
		for _, item := range resp {
			fmt.Println("ert:", "gfhr")
			response = append(response, item)
		}
		// result := append(response, resp...)
		// return response
	}

	return response
}