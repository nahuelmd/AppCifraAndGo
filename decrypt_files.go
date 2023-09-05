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

	userDirectories, err := getUserDirectories()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, directory := range userDirectories {
		err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
			return decryptFileAtPath(path, f, err, []byte(secretKey))
		})

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func decryptFileAtPath(path string, f os.FileInfo, err error, key []byte) error {
	if err != nil {
		return err
	}

	if !f.IsDir() && strings.HasSuffix(path, ".enc") {
		if err := decryptAndReplaceFile(path, key); err != nil {
			return fmt.Errorf("Decryption error: %v", err)
		}
	}
	return nil
}

func getUserDirectories() ([]string, error) {
	switch runtime.GOOS {
	case "windows":
		return []string{"C:\\Users\\"}, nil
	case "linux":
		return []string{"/home/"}, nil
	case "darwin":
		return []string{"/Users/"}, nil
	default:
		return nil, fmt.Errorf("unsupported OS")
	}
}

func decryptAndReplaceFile(inputFile string, key []byte) error {
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
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	newPath := strings.TrimSuffix(inputFile, ".enc")
	if err = ioutil.WriteFile(newPath, plainText, 0777); err != nil {
		return err
	}

	if err = os.Remove(inputFile); err != nil {
		return err
	}

	return nil
}
