package main

import (
	"harr1424/go_notify/gonotify"
	"log"
	"net/http"
	"time"
)

func main() {

	err := gonotify.ReadRemoteTableContents()
	if err != nil {
		log.Fatal("Failed to read remote table contents: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", gonotify.RegisterToken)
	mux.HandleFunc("/add_location", gonotify.HandleLocationAdd)
	mux.HandleFunc("/remove_location", gonotify.HandleLocationRemove)
	mux.HandleFunc("/test_dynamo", testDynamoDBInsert)

	go func() {
		for {
			gonotify.CheckAllLocationsForFrost()
			gonotify.UpdateTokenLocationMap(gonotify.TokenLocationMap)
			time.Sleep(12 * time.Hour)
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5050", mux))
	//log.Fatal(http.ListenAndServeTLS(":5050", "localhost.crt", "localhost.key", nil)) // support TLS when available
}

func testDynamoDBInsert(res http.ResponseWriter, req *http.Request) {
	gonotify.UpdateTokenLocationMap(gonotify.TokenLocationMap)
}
