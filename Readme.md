My first Malware with GO!!

With this code you can generate a Ransomware that will encrypt all files in the OS user folder

Works for Windows, MacOS and Linux

This software was generated for educational purposes. NEVER use this software for malicious purposes.

##### Never run this code on a production machine

Usage:

To compile the executable file you must pass the secret key as an environment variable. If you compile without passing the secretKey variable, the program will automatically assign "MyVerySecretKeyForEncryption1234" as the secret key.

Compile for Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.secretKey=MyVerySecretKeyForEncryption1234'" -o encrypt_files_linux encrypt_files.go
GOOS=linux GOARCH=amd64 go build -o decrypt_files_linux decrypt_files.go

Compile for Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.secretKey=MyVerySecretKeyForEncryption1234'" -o encrypt_files_win.exe encrypt_files.go
GOOS=windows GOARCH=amd64 go build -o decrypt_files_win.exe decrypt_files.go

Compile for MacOs
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.secretKey=MyVerySecretKeyForEncryption1234'" -o encrypt_files_mac encrypt_files.go
GOOS=darwin GOARCH=amd64 go build -o decrypt_files_mac decrypt_files.go
