package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var filetypeFind = ".exe"
var whereFind = "C:\\tmp\\MDF to ISO"
var distFolder = "/copy/"

func find(root, ext string) (filenameAndFilepath map[string]string) {
	filenameAndFilepath = make(map[string]string)
	filepath.WalkDir(root, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			filenameAndFilepath[d.Name()] = s
		}
		return nil

	})
	return filenameAndFilepath
}

func copy(src, dst string) {
	r, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	w, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	w.ReadFrom(r)
}

func main() {
	// ========= узнать pwd
	whereCopy, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	// =========

	filenameAndFilepath := find(whereFind, filetypeFind)
	err = os.MkdirAll(whereCopy+distFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for name, path := range filenameAndFilepath {
		copy(path, whereCopy+distFolder+name)
	}

}
