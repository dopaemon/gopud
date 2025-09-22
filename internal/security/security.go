package security

import (
	"bytes"
	"errors"
	"fmt"
	_ "io"

	"github.com/minio/sio"
)

func EncryptData(data, key []byte) ([]byte, error) {
	var buf bytes.Buffer
	_, err := sio.Encrypt(&buf, bytes.NewReader(data), sio.Config{Key: key})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecryptData(encrypted, key []byte) ([]byte, error) {
	var buf bytes.Buffer
	_, err := sio.Decrypt(&buf, bytes.NewReader(encrypted), sio.Config{Key: key})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetAPIKey() {
	key := []byte("12345678901234567890123456789012")
	message := []byte("Hello, secure world!")

	encrypted, _ := EncryptData(message, key)
	fmt.Println("Encrypted:", encrypted)

	decrypted, err := DecryptData(encrypted, key)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Decrypted:", string(decrypted))
	}

	wrongKey := []byte("00000000000000000000000000000000")
	decrypted, err = DecryptData(encrypted, wrongKey)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Decrypted:", string(decrypted))
	}
}

