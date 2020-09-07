package main

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"time"
)

func decrypt(key []byte, nonceString, encryptedString string) (url url.URL, expiry time.Time, headers map[string]string) {
	nonce, _ := b64.URLEncoding.DecodeString(nonceString)
	ciphertext, _ := b64.URLEncoding.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	checkError(err)

	//Create a new GCM
	aesGCM, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	checkError(err)

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	checkError(err)

	//Parse the URL and expiry
	payload := strings.Split(string(plaintext), ";")
	layout := "2006-01-02T15:04:05.000000"
	t, err := time.Parse(layout, payload[1])
	checkError(err)

	var headerJson map[string]string
	json.Unmarshal([]byte(payload[2]), &headerJson)

	u, err := url.Parse(payload[0])
	checkError(err)

	return *u, t, headerJson
}
