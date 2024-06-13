# copy-to-pwd
build:

GOOS=windows GOARCH=amd64 go build -o ./copytoflash.exe main.go

build silent mode (without window):

GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ./copytoflash.exe main.go
