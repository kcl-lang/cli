package fs

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// IsURL checks if the given path is an HTTP or HTTPS URL.
func IsURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// GenTempFileFromURL fetches content from a URL and saves it to a temp file.
// It returns the path to the temp file or an error if the fetch fails.
// The caller is responsible for removing the temp file after use.
func GenTempFileFromURL(urlStr string) (string, error) {
	// Parse the URL to extract the file extension
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL %q: %w", urlStr, err)
	}

	// Extract file extension from URL path for temp file naming
	ext := filepath.Ext(parsedURL.Path)
	if ext == "" {
		ext = ".tmp"
	}

	// Create temp file with appropriate extension
	tempFile, err := os.CreateTemp("", "kcl-import-*"+ext)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Fetch content from URL
	resp, err := http.Get(urlStr)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to fetch URL %q: %w", urlStr, err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to fetch URL %q: HTTP %d %s", urlStr, resp.StatusCode, resp.Status)
	}

	// Copy response body to temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to save content from URL %q: %w", urlStr, err)
	}

	return tempFile.Name(), nil
}

func GenTempFileFromStdin() (string, error) {
	tempFile, err := os.CreateTemp("", "stdin-*.k")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(tempFile, os.Stdin)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
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

// IsEmptyDir checks if a directory is empty.
// It takes a string parameter `name` representing the directory path.
// It returns a boolean value indicating whether the directory is empty or not,
// and an error if any occurred during the process.
//
// Parameters:
// - name: The path of the directory to check.
//
// Returns:
// - bool: True if the directory is empty, false otherwise.
// - error: An error if the directory cannot be read.
//
// Example usage:
// empty, err := IsEmptyDir("/path/to/directory")
func IsEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
