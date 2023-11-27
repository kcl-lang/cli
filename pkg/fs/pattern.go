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
