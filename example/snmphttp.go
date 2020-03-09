//package main
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//	"time"
//
//	"github.com/aleasoluciones/gosnmpquerier"
//)
//
//// Example of execution
//// Run the http server
////		go run examples/snmphttpserver/snmphttpserver.go
//// curl a request
////		curl http://127.0.0.1:8080 -X PUT  -H 'Content-Type: multipart/form-data' -d '{"cmd":"walk", "destination":"MY_HOST_IP", "community":"MY_COMMUNITY", "oid":["AN_OID"]}'
//
//const (
//	CONTENTION = 4
//)
//
//func rootHandler(querier gosnmpquerier.SyncQuerier, w http.ResponseWriter, r *http.Request) {
//
//	cmd, _ := gosnmpquerier.ConvertCommand(r.FormValue("cmd"))
//	oid := r.FormValue("oid")
//	community := r.FormValue("community")
//	dest := r.FormValue("destination")
//
//
//	query := gosnmpquerier.Query{
//		Cmd:         cmd,
//		Community:   community,
//		Oids:        []string{oid},
//		Timeout:     time.Duration(10) * time.Second,
//		Retries:     1,
//		Destination: dest,
//	}
//
//	processed := querier.ExecuteQuery(query)
//	jsonProcessed, err := gosnmpquerier.ToJson(&processed)
//
//	if err != nil {
//		fmt.Fprint(w, err)
//	}
//
//	fmt.Fprint(w, jsonProcessed)
//
//}
//
//func main() {
//
//	querier := gosnmpquerier.NewSyncQuerier(CONTENTION, 3, 3*time.Second)
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		rootHandler(querier, w, r)
//	})
//
//	log.Fatal(http.ListenAndServe(":9091", nil))
//
//}






package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

// Example of execution
// go run examples/get/get.go -community MY_COMMUNITY -host MY_HOST_IP  AN_OID

func main() {

	community := flag.String("community", "public", "snmp v2 community")
	host      := flag.String("host", "192.168.2.148", "host")
	timeout   := flag.Duration("timeout", 1*time.Second, "Timeout (ms/s/m/h)")
	retries   := flag.Int("retries", 1, "Retries")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] [[oid] ...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	oids := []string{".1.3.6.1.2.1.1.9.1.3.4"}
	for _, oid := range flag.Args() {
		oids = append(oids, oid)
	}

	querier := gosnmpquerier.NewSyncQuerier(1, 3, 3*time.Second)
	result, err := querier.Get(*host, *community, oids, *timeout, *retries)
	fmt.Println(result, err)

}

