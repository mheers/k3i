package helpers

import (
	"os"
	"strings"
)

func getDir(path string) string {
	index := strings.LastIndex(path, "/")
	if index == -1 {
		return "./"
	}
	return path[:index]
}

func CreateFileWithDirs(path string) (*os.File, error) {
	// create file with dirs
	err := os.MkdirAll(getDir(path), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
