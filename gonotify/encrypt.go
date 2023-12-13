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
and decrypt tokens read from a file. A Key and nonce are defined in order to perform AES-GCM
authenticated encryption.

IMPORTANT: The Key is written to a local file.
Whenever possible, it should be stored in a more secure location,
separate from the data that it encrypts.

Even better, CSPs offer database solutions better suited to storing devicd tokens.
In the future, one of these solutions will be used insteead.
*/

var Key = make([]byte, 32)

func encryptToken(t token) string {
	original := t.ID // ID is string member of token
	var nonce = make([]byte, 12)

	// read random bytes into nonce
	_, err := rand.Read(nonce)
	if err != nil {
		log.Println("Error reading random bytes into nonce:", err)
	}

	block, err := aes.NewCipher(Key)
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

	block, err := aes.NewCipher(Key)
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

func HandleCrypto(Key *[]byte) {
	if _, err := os.Stat("Key.Key"); errors.Is(err, os.ErrNotExist) {
		log.Println("Key file not found. Creating one...")
		Key_file, err := os.OpenFile("Key.Key", os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Println("Error creating file:", err)
		}

		defer Key_file.Close()

		_, err = rand.Read(*Key)
		if err != nil {
			log.Println("Error creating Key:", err)
		}

		_, err = Key_file.Write(*Key)
		if err != nil {
			log.Println("Error writing Key to file:", err)
		}

	} else {
		log.Println("Crypto files found, attempting to read...")
		Key_file, err := os.OpenFile("Key.Key", os.O_RDONLY, 0644)
		if err != nil {
			log.Println("Error accessing Key file:", err)
		}

		defer Key_file.Close()

		_, err = Key_file.Read(*Key)
		if err != nil {
			log.Println("Error reading Key from file:", err)
		}
	}

}
