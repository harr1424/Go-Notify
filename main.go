package main

import (
	"fmt"
	"harr1424/go_notify/gonotify"
	"log"
	"net/http"
)

func main() {
	gonotify.HandleCrypto(&gonotify.Key)

	gonotify.ReadTokensFromFile()
	fmt.Println("All tokens (from file):", gonotify.Tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", gonotify.CreateNewToken)

	log.Fatal(http.ListenAndServe("0.0.0.0:5050", nil))
	//log.Fatal(http.ListenAndServeTLS(":5050", "localhost.crt", "localhost.key", nil)) // support TLS when available

}
