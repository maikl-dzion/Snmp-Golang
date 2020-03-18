package model

import (
	"fmt"
	"log"
	"os"
	"time"
)

//type ResponseMessage struct {
//	Oid      string `json:"oid"`
//	Value    string `json:"value"`
//	Ip       string `json:"ip"`
//	DeviceId string `json:"device_id"`
//}
//
//type ResponseJsonItems []ResponseMessage

const QUEUE_NAME   = "SNMP_QUEUE"
const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"
const SAVE_API_URL = "http://172.16.16.235:8080/data/multi_save"
const SNMP_FUNC_TYPE = "GET"
const AMQP_FUNC_TYPE = "get"
const LOGS_PATH = "logs"

type SnmpSendParams struct {
	Id        string
	Ip        string
	Oid       string
	Type      string
	Port      string
	Retry     int
	DeviceId  string
	Community string
	SelCount  int
}

/////////////////////////////////////
// ---- Вспомогательные функции ----
func DatetimePrint() {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d__%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println("DateTime->", formatted)
}

func LogPrint(data interface{}, name string) {
	if name != "" {
		fmt.Println(name, data)
	} else {
		fmt.Println(data)
	}
}


func LogSave(msg string, err error, fileName string) {

	if fileName == "" {
		fileName = "log"
	}
    logFileName := LOGS_PATH + "/" +fileName+ ".log"
	logfile, errLog := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errLog != nil {
		log.Fatalln(errLog)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Printf("[Msg]:%s -||||- [Error]: %s", msg, err)

}


func WarnOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		LogSave(msg, err, "WarningErrors")
	}
}
