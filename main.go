package main

import "harr1424/go_notify/gonotify"

func main() {
	// handleCrypto(&key)

	// readTokensFromFile()
	// fmt.Println("All tokens (from file):", tokens)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/register", createNewToken)

	// log.Fatal(http.ListenAndServe("0.0.0.0:5050", nil))
	//log.Fatal(http.ListenAndServeTLS(":5050", "localhost.crt", "localhost.key", nil)) // support TLS when available

	gonotify.GetForecast()
}