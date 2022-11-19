package files

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type DirStructure struct {
	Files       []string
	Directories []string
}

func CopyFile(destFile string, sourceFile string) error {
	file, createErr := os.Create(destFile)
	src, openErr := os.Open(sourceFile)
	if createErr != nil {
		return createErr
	}
	if openErr != nil {
		return openErr
	}
	_, copyErr := io.Copy(file, src)
	if copyErr != copyErr {
		return copyErr
	}
	return nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ReadDirStructure(target string) (DirStructure, error) {
	dirStruct := DirStructure{}
	filepath.WalkDir(target, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		path, err = filepath.Rel(target, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirStruct.Directories = append(dirStruct.Directories, path)
		} else {
			dirStruct.Files = append(dirStruct.Files, path)
		}
		return nil
	})
	return dirStruct, nil
}

func IsRequiredOverwirte(path string) bool {
	for {
		var ans string
		fmt.Println("\"", path, "\"", "is already exist. Overwrite? [y/n]")
		fmt.Scan(&ans)
		if ans == "y" {
			return true
		} else if ans == "n" {
			return false
		} else {
			continue
		}
	}
}

func Create(path string) error {
	if mkdirErr := os.MkdirAll(filepath.Dir(".\\"+path), 0666); mkdirErr != nil {
		return mkdirErr
	}
	f, openErr := os.OpenFile(".\\"+path, os.O_CREATE, 0666)
	if openErr != nil {
		return openErr
	}
	f.Write([]byte{})

	return nil
}
