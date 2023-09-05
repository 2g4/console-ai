// AES-128 to encrypt and decrypt messages.
// The key is a 16 byte array.
// Primarly used to encrypt the API key.

package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func EncryptMessage(key []byte, message string) (string, error) {
	// Convert the message string to a byte slice
	byteMsg := []byte(message)

	// Create a new AES cipher using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	// Prepare the slice to hold the initialization vector (IV) and encrypted data
	cipherText := make([]byte, aes.BlockSize+len(byteMsg))

	// Extract the slice for the initialization vector (IV)
	iv := cipherText[:aes.BlockSize]

	// Populate the IV with random data
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not generate IV: %v", err)
	}

	// Encrypt the message
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	// Convert the cipherText to a base64 encoded string and return
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(key []byte, encodedMessage string) (string, error) {
	// Decode the base64 encoded message to get the ciphertext
	cipherText, err := base64.StdEncoding.DecodeString(encodedMessage)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	// Create a new AES cipher using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	// Check if the ciphertext is of valid size
	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	// Separate the initialization vector (IV) and the actual encrypted data
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
