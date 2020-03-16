package snmp_handler

import (
	"fmt"
	"log"
	"os"
	"time"

	model "Snmp-Golang/internal/models"
	snmp "github.com/soniah/gosnmp"
)

type PduValue struct {
	name   string
	column string
	value  interface{}
}

func SnmpConfigInit(s model.SnmpSendParams) (*snmp.GoSNMP, error) {
	client := &snmp.GoSNMP{
		Target:    s.Ip,
		Port:      uint16(161),
		Community: s.Community,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(10) * time.Second,
		//Logger:    log.New(os.Stdout, "", 0),
		//Retries:   s.Retries,
	}
	err := client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return client, nil
}

func SnmpManagerStart(conf model.SnmpSendParams, fType string) (SnmpResultItems, error) {

	now := time.Now()

	client, err := SnmpConfigInit(conf)
	if err != nil {
		log.Fatal("Snmp Config Is Nil")
	}
	defer client.Conn.Close()

	snmpResult := SnmpResultItems{}
	snmpResult.DeviceId = conf.DeviceId
	snmpResult.Ip = conf.Ip
	oid := conf.Oid

	switch fType {
	case "GET":
		oids := []string{oid}
		pkt, err := client.Get(oids)
		if err != nil {
			return snmpResult, err
		}
		for _, pdu := range pkt.Variables {
			snmpResult.CollectValues(pdu)
		}

	case "BULK":
		err := client.BulkWalk(oid, snmpResult.CollectValues)
		if err != nil {
			return snmpResult, err
		}

	case "BULK_WALK_ALL":
		pkt, err := client.BulkWalkAll(oid)
		if err != nil {
			return snmpResult, err
		}
		for _, pdu := range pkt {
			snmpResult.CollectValues(pdu)
		}

	default:
		err := client.BulkWalk(oid, snmpResult.CollectValues)
		if err != nil {
			return snmpResult, err
		}
	}

	ch := len(snmpResult.Items)
	model.LogPrint(ch, "Snmp Result Oid Count:")
	fmt.Println("Now", now)

	return snmpResult, nil
}
