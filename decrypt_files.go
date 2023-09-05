// decrypt_files.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	fmt.Print("Enter the secret key: ")
	var secretKey string
	fmt.Scanln(&secretKey)

	if len(secretKey) == 0 {
		fmt.Println("The secret key cannot be empty.")
		return
	}

	userDirectories := getUserDirectories()

	for _, directory := range userDirectories {
		err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
			return visit(path, f, err, []byte(secretKey), false)
		})

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func visit(path string, f os.FileInfo, err error, key []byte, encrypt bool) error {
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	if !f.IsDir() {
		if !encrypt && strings.HasSuffix(path, ".enc") {
			fmt.Println("Decrypting file:", path)
			err := decryptFile(path, key)
			if err != nil {
				fmt.Println("Decryption error:", err)
			} else {
				//Delete .enc file after successful decryption
				os.Remove(path)
			}
		}
	}

	return nil
}

func getUserDirectories() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"C:\\Users\\"}
	case "linux":
		return []string{"/home/"}
	case "darwin":
		return []string{"/Users/"}
	default:
		return []string{}
	}
}

func decryptFile(inputFile string, key []byte) error {
	ciphertext, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	newPath := strings.TrimSuffix(inputFile, ".enc")
	err = ioutil.WriteFile(newPath, plainText, 0777)
	if err != nil {
		return err
	}

	return nil
}
