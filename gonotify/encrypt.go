package gonotify

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
)

/*
This file contains logic to encrypt a token before it is stored on the local filesystem,
and decrypt tokens read from a file. A key and nonce are defined in order to perform AES-GCM
authenticated encryption.

IMPORTANT: The key is written to a local file.
Whenever possible, it should be stored in a more secure location,
separate from the data that it encrypts.

Even better, CSPs offer database solutions better suited to storing devicd tokens. 
In the future, one of these solutions will be used insteead.
*/

var key = make([]byte, 32)

func encryptToken(t token) string {
	original := t.ID // ID is string member of token
	var nonce = make([]byte, 12)

	// read random bytes into nonce
	_, err := rand.Read(nonce)
	if err != nil {
		log.Println("Error reading random bytes into nonce:", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Error creating cipher during encrypt:", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Error creating GCM during encrypt:", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(original), nil)

	// prepend the ciphertext with the nonce
	out := append(nonce, ciphertext...)

	return hex.EncodeToString(out)
}

func decryptToken(s string) string {
	// read hex string describing nonce and ciphertext
	enc, err := hex.DecodeString(s)
	if err != nil {
		log.Println("Error decoding string from hex:", err)
	}

	// separate ciphertext from nonce
	nonce := enc[0:12]
	ciphertext := enc[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Error creating cipher during decrypt:", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Error creating GCM during decrypt:", err)
	}

	original, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Error decrypting to string:", err)
	}
	originalAsString := string(original)

	return originalAsString
}

func handleCrypto(key *[]byte) {
	if _, err := os.Stat("key.key"); errors.Is(err, os.ErrNotExist) {
		log.Println("Key file not found. Creating one...")
		key_file, err := os.OpenFile("key.key", os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Println("Error creating file:", err)
		}

		defer key_file.Close()

		_, err = rand.Read(*key)
		if err != nil {
			log.Println("Error creating key:", err)
		}

		_, err = key_file.Write(*key)
		if err != nil {
			log.Println("Error writing key to file:", err)
		}

	} else {
		log.Println("Crypto files found, attempting to read...")
		key_file, err := os.OpenFile("key.key", os.O_RDONLY, 0644)
		if err != nil {
			log.Println("Error accessing key file:", err)
		}

		defer key_file.Close()

		_, err = key_file.Read(*key)
		if err != nil {
			log.Println("Error reading key from file:", err)
		}
	}

}
