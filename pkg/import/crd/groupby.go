// Copy from https://github.com/DavidChevallier/CRDtoKCL/blob/main/main.go
package crd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
)

// knownAPIVersions is a slice of strings that contains the known API versions.
var knownAPIVersions = []string{
	"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10",
	"v1alpha1", "v1alpha2", "v1alpha3", "v1alpha4", "v1alpha5",
	"v2alpha1", "v2alpha2", "v2alpha3", "v2alpha4", "v2alpha5",
	"v3alpha1", "v3alpha2", "v3alpha3", "v3alpha4", "v3alpha5",
	"v1beta1", "v1beta2", "v1beta3", "v1beta4", "v1beta5",
	"v2beta1", "v2beta2", "v2beta3", "v2beta4", "v2beta5",
	"v3beta1", "v3beta2", "v3beta3", "v3beta4", "v3beta5",
}

const unknownVersion = "unknown"

// moveKclFiles moves files with the ".k" extension to a directory based on their API version.
// It walks through the specified base directory and for each file with the ".k" extension,
// it determines the API version based on the file name and moves the file to a subdirectory
// named after the API version. If the API version cannot be determined, the file is moved
// to a subdirectory named "unknown".
//
// Parameters:
// - baseDir: The base directory to search for files.
func GroupByKclFiles(baseDir string) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".k") {
			continue
		}

		apiVersion, err := extractAPIVersionFromName(entry.Name(), knownAPIVersions)
		if err != nil {
			apiVersion = unknownVersion
		}

		newDir := filepath.Join(baseDir, apiVersion)
		os.MkdirAll(newDir, os.ModePerm)

		oldPath := filepath.Join(baseDir, entry.Name())
		newPath := filepath.Join(newDir, entry.Name())

		err = os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
	}
	kclPkg := pkg.NewKclPkg(&opt.InitOptions{
		Name:     filepath.Base(baseDir),
		InitPath: baseDir,
	})
	err = kclPkg.ModFile.StoreModFile()
	if err != nil {
		return err
	}
	return removeEmptyDirs(baseDir)
}

// removeEmptyDirs removes all empty directories within the specified directory.
// It recursively walks through the directory and checks if each directory is empty.
// If an empty directory is found, it is removed. The function stops when there are no more empty directories left.
//
// Parameters:
// - dir: The directory path to start the search from.
//
// Example usage:
// removeEmptyDirs("/path/to/directory")
//
// Note: This function does not handle errors related to directory traversal or removal.
func removeEmptyDirs(dir string) error {
	for {
		var emptyDirs []string

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				empty, err := fs.IsEmptyDir(path)
				if err != nil {
					return err
				}
				if empty {
					emptyDirs = append(emptyDirs, path)
				}
			}
			return nil
		})

		if len(emptyDirs) == 0 {
			break
		}

		for i := len(emptyDirs) - 1; i >= 0; i-- {
			err := os.Remove(emptyDirs[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// extractAPIVersionFromName extracts the API version from the given name string.
// It searches for a pattern in the name string and returns the extracted API version.
// If the API version is found and is known (present in the knownAPIVersions slice),
// it returns the API version string. Otherwise, it returns "unknown".
//
// Parameters:
// - name: The name string to extract the API version from.
// - knownAPIVersions: A slice of known API versions.
//
// Returns:
// - string: The extracted API version or "unknown" if not found.
// - error: An error if the extraction fails.
//
// Example usage:
// version, err := extractAPIVersionFromName("example_v1alpha1", knownAPIVersions)
func extractAPIVersionFromName(name string, knownAPIVersions []string) (string, error) {
	re := regexp.MustCompile(`_v([0-9a-zA-Z]+)`)
	matches := re.FindStringSubmatch(name)
	if len(matches) > 1 {
		apiVersion := "v" + matches[1]
		for _, v := range knownAPIVersions {
			if apiVersion == v {
				return apiVersion, nil
			}
		}
	}
	return unknownVersion, nil
}
