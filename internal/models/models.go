package model

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//const QUEUE_NAME   = "SNMP_QUEUE"
//const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"
//const SAVE_API_URL = "http://172.16.16.235:8080/data/multi_save"
//const SNMP_FUNC_TYPE = "GET"
//const AMQP_FUNC_TYPE = "get"
//const SNMP_FUNC_TYPE = "BULK"
//const AMQP_FUNC_TYPE = "consumer"

const LOGS_PATH   = "logs"
const CONFIG_FILE = "config.toml"
const CONFIG_PATH = "config"

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


type CommonInitParam struct {
	QueueName  string
	AmqpUrl     string
	SaveApiUrl  string
	AmqpFuncType string
	SnmpFuncType string
}


func GetCommonInitParam() CommonInitParam {

	configFile := CONFIG_FILE
	configPath := GetRootDir(CONFIG_PATH)
	path := configPath + "/" +configFile

	var config CommonInitParam
	if _, err := toml.DecodeFile(path, &config); err != nil {
		FailOnError(err, "FATAL ERROR : Failed to open config file")
	}


	return config
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

	logsPath := GetRootDir(LOGS_PATH)
    logFileName := logsPath + "/" +fileName+ ".log"
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


func EventLogSave(msg, eventName, fileName string) {

	if fileName == "" {
		fileName = "events"
	}

	if eventName == "" {
		eventName = "CommonEvent"
	}

	logsPath := GetRootDir(LOGS_PATH)
	logFileName := logsPath + "/" +fileName+ ".log"

	logfile, errLog := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errLog != nil {
		log.Fatalln(errLog)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Printf("[Name]:%s ___ [Msg]: %s", eventName, msg)

}


func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		LogSave(msg, err, "FatalErrors")
		log.Fatalf("%s: %s", msg, err)
	}
}

func FatalError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		LogSave(msg, err, "FatalErrors")
		log.Fatalf("%s: %s", msg, err)
	}
}



func GetRootDir(localDirName string) string {

	rootPath, errDir := filepath.Abs(".")
	if errDir != nil {
		fmt.Println(errDir, "Not Load Root Dir")
	}

	delimiter := "/"
	// delimiter := "/../"
	rootDir := rootPath + delimiter + localDirName
	fmt.Println("Root Dir:", rootDir)

	return rootDir
}



func DateTypeConvert(num int, str string) (int, string) {

	if num != 0 {
		s := strconv.Itoa(num)
		return num, s
	} else {

		if str == "" {
			return num, str
		}

		d, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println(err)
		}
		return d, str
	}
}