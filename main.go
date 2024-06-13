package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var filetypeFind = ".pdf"
var whereFind = [...]string{"C", "E"}
var distFolder = "copy"
var logFile = "copy.log"

// Эта функция проверяет, нужно ли пропустить указанную директорию
func shouldSkipDirectory(path string) bool {
	// Список директорий, которые необходимо пропустить
	skipDirs := []string{"ProgramData", "Windows"}

	// Проверяем, содержится ли имя директории в списке для пропуска
	for _, dir := range skipDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

func findAndCopy(folderToSearch, fileExt, whereCopy string) error {
	err := filepath.WalkDir(folderToSearch, func(s string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Skipping directory due to permission error: %s\n", s)
				time.Sleep(8 * time.Second)
				return nil
			}
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == fileExt {
			// Пропускаем копирование, если директория должна быть пропущена
			if shouldSkipDirectory(filepath.Dir(s)) {
				fmt.Printf("Skipping directory: %s\n", filepath.Dir(s))
				return nil
			}
			err := copyFile(s, filepath.Join(whereCopy, distFolder, d.Name()))
			if err != nil {
				log.Printf("Error copying file %s: %v\n", s, err)
			} else {
				log.Printf("Copied: %s\n", s)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking directory %s: %v\n", folderToSearch, err)
	}
	return nil
}

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

	// Используем буферизацию для улучшения производительности
	buf := make([]byte, 1024)
	_, err = io.CopyBuffer(dstFile, srcFile, buf)
	return err
}

func main() {
	logf, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logf.Close()
	log.SetOutput(logf)

	whereCopy, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(whereCopy, distFolder), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create destination directory: %v", err)
	}

	for _, drive := range whereFind {
		searchPath := filepath.Join(drive + ":\\Users")
		log.Printf("Scanning %s for %s files...\n", searchPath, filetypeFind)
		err := findAndCopy(searchPath, filetypeFind, whereCopy)
		if err != nil {
			log.Printf("Error scanning and copying files: %v\n", err)
		}
	}
}
