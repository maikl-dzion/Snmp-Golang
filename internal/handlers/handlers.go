package handler

import (
	"../models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k-sone/snmpgo"
	"log"
	"net/http"
	"sync"
)


func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}



func SnmpGet(ip string, oids snmpgo.Oids, port string) {


	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   ip + ":" + port,
		Retries:   1,
		Community: "public",
	})


	if err != nil {
		fmt.Printf("[%s] : construct error - %s\n", ip, err)
		return
	}

	if err = snmp.Open(); err != nil {
		fmt.Printf("[%s] : open error - %s\n", ip, err)
		return
	}

	defer snmp.Close()

	pdu, err := snmp.GetRequest(oids)
	if err != nil {
		fmt.Printf("[%s] : get error - %s\n", ip, err)
		return
	}

	if pdu.ErrorStatus() != snmpgo.NoError {
		fmt.Printf("[%s] : error status - %s at %d\n",
			ip, pdu.ErrorStatus(), pdu.ErrorIndex())
	}


	resp := pdu.VarBinds()[0];

	// fmt.Printf("[%s] : %s\n", ip, pdu.VarBinds())
	fmt.Println(resp)
	fmt.Println(resp.Oid.Value)
	fmt.Println(resp.Variable.String())
	fmt.Println(resp.Variable.Type())


	message := model.ResponseMessage{}
	message.Oid = resp.Oid.String()
	message.Value = resp.Variable.String()
	message.Ip = ip


	fmt.Println(message);

	CurlRun(message)

}


func CurlRun(message model.ResponseMessage) {

	payloadBytes, err := json.Marshal(message)
	if err != nil {
		// handle err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", model.SAVE_API_URL, body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer b7d03a6947b217efb6f3ec3bd3504582")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	defer resp.Body.Close()

}



func SnmpRequestInit(ipAddress string, oid string, port string ) {

	oidList, err := snmpgo.NewOids([]string{oid})

	if err != nil {
		fmt.Println(err)
		return
	}


	var wg sync.WaitGroup

	wg.Add(1)

	go func(ip string, p string) {

		defer wg.Done()

		SnmpGet(ip, oidList, p)

	}(ipAddress, port)

	wg.Wait()
}

