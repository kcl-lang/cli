package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAllFilesInFolder(folderPath string, recursive bool) ([]string, error) {
	var fileList []string

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() && recursive {
			subFolderFiles, err := GetAllFilesInFolder(filepath.Join(folderPath, file.Name()), recursive)
			if err != nil {
				return fileList, fmt.Errorf("error while reading files from subfolder: %s", err)
			}
			fileList = append(fileList, subFolderFiles...)
		} else if !file.IsDir() {
			fileList = append(fileList, filepath.Join(folderPath, file.Name()))
		}
	}
	return fileList, nil
}

func IgnoreFile(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return false
	}
	ext := filepath.Ext(path)
	for _, s := range extensions {
		if s == ext {
			return false
		}
	}
	return true
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func FileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil || fi.IsDir() {
		return false
	}
	return true
}
