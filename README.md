# copy-to-pwd

Скрипт для поиска и копирования файлов определенного в conf.txt расширения файлов со всех дисков, кроме того с которого запущено приложение, в директорию рядом с приложением.
По умолчанию запускается с правами администратора и повышенными привелегиями доступа к директориям.

# сборка:

build:

  `GOOS=windows GOARCH=amd64 go build -o ./copytoflash.exe main.go`

build silent mode (without window):

  `GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ./copytoflash.exe main.go`
