package main

import (
	"harr1424/go_notify/gonotify"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", gonotify.RegisterToken)
	mux.HandleFunc("/add_location", gonotify.HandleLocationAdd)

	go func() {
		for {
			gonotify.CheckAllLocationsForFrost()
			time.Sleep(24 * time.Hour)
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5050", mux))
	//log.Fatal(http.ListenAndServeTLS(":5050", "localhost.crt", "localhost.key", nil)) // support TLS when available
}
