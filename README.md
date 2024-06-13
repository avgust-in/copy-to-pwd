# copy-to-pwd

Скрипт для копирования pdf со всех дисков кроме того с которого запущено приложение в директорию рядом с приложением.
По умолчанию запускается с правами администратора и повышенными привелегиями доступа к директориям.

При желании можно указать:
директорию для сканирования
директории для пропуска сканирования
расширение файлов

# сборка:

build:

GOOS=windows GOARCH=amd64 go build -o ./copytoflash.exe main.go

build silent mode (without window):

GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ./copytoflash.exe main.go
