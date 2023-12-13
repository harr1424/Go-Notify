package gonotify

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

//TODO Use CSP solution to store device tokens and associated locations

/*
This file contains logic required to register device tokens when the register endpoint is contacted,
a helper function to remove duplicae tokens,
logic to maintain a file describing all tokens,
and logic to read this file during application startup.

In the future, this file will be updated to store tokens and locations in a CSP database offering
encryption.
*/

// Represents a token corresponding to an iOS device and associated locations
type token struct {
	ID string
}

// A slice of all tokens
var tokens []token

// Helper function to remove duplicate tokens
func removeDuplicateTokens(tokenSlice []token) []token {
	keys := make(map[token]bool)
	var set []token
	for _, token := range tokenSlice {
		if _, value := keys[token]; !value {
			keys[token] = true
			set = append(set, token)
		}
	}
	return set
}

// Called when the register endpoint is contacted
// Expects to receive POST data describing an iOS device token
func createNewToken(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/register" {
		http.NotFound(res, req)
		return
	}

	var newToken token

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&newToken); err != nil {
		log.Println("Could not create new token from request body: ", err)
		return
	}

	tokens = append(tokens, newToken)
	tokens = removeDuplicateTokens(tokens)
	createOrUpdateTokenFile(newToken)
	res.WriteHeader(http.StatusCreated)
	fmt.Println("All tokens (in memory):", tokens)
}

/*
Called by createNewToken() in order to append the new token to a CSV list of all
tokens stored on local filesystem. This function checks if the file exists,
and if it does not, it is created.
*/
func createOrUpdateTokenFile(t token) {
	file, err := os.OpenFile("tokens.data", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		log.Println("Error accessing file:", err)
	}

	defer file.Close()

	log.Println("Writing token to file...")
	encoder := csv.NewWriter(file)

	encrypted := encryptToken(t)
	encoder.Write([]string{encrypted})
	encoder.Flush()
	err = encoder.Error()
	if err != nil {
		log.Println(err)
	}

}

/*
Called at application startup. If no token file is present, it will be crated.
Otherwise, the contents of the file are read to the tokens slice in memory.
*/
func readTokensFromFile() {
	if _, err := os.Stat("tokens.data"); errors.Is(err, os.ErrNotExist) {
		log.Println("Token file not found... Creating one...")
		file, err := os.Create("tokens.data")
		if err != nil {
			log.Println("Error creating token file:", err)
		}
		defer file.Close()
	} else {
		file, err := os.Open("tokens.data")
		if err != nil {
			log.Println("Error opening file:", err)
		}

		defer file.Close()

		reader := csv.NewReader(file)
		for err == nil {
			var s []string

			s, err = reader.Read()
			if len(s) > 0 {
				decrypted := decryptToken(s[0])
				tokens = append(tokens, token{ID: decrypted})
			}
		}

		tokens = removeDuplicateTokens(tokens)
	}
}
