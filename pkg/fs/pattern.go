package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExpandInputFiles returns all the filenames that match the input filename,
// filepath or patterns.
func ExpandInputFiles(files []string, recursive bool) ([]string, error) {
	var result []string
	for _, f := range files {
		expandFiles, err := ExpandIfFilePattern(f, recursive)
		if err != nil {
			return result, err
		}
		result = append(result, expandFiles...)
	}
	return result, nil
}

// ExpandIfFilePattern returns all the filenames that match the input pattern
// or the filename if it is a specific filename and not a pattern.
// If the input is a pattern and it yields no result it will result in an error.
func ExpandIfFilePattern(pattern string, recursive bool) ([]string, error) {
	if _, err := os.Stat(pattern); os.IsNotExist(err) {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) == 0 {
			return nil, fmt.Errorf("the path %q does not exist", pattern)
		}
		if err == filepath.ErrBadPattern {
			return nil, fmt.Errorf("pattern %q is not valid: %v", pattern, err)
		}
		return matches, err
	}
	if IsDir(pattern) {
		return GetAllFilesInFolder(pattern, recursive)
	}
	return []string{pattern}, nil
}

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
