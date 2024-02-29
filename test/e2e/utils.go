package e2e

import (
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/otiai10/copy"
	"github.com/thoas/go-funk"
	"kcl-lang.io/kpm/pkg/reporter"
)

// Copy will copy file from 'srcPath' to 'dstPath'.
func Copy(srcPath, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}

// CopyDir will copy dir from 'srcDir' to 'dstDir'.
func CopyDir(srcDir, dstDir string) {
	err := copy.Copy(srcDir, dstDir)
	if err != nil {
		reporter.ExitWithReport("failed to copy dir.")
	}
}

var KEYS = []string{"<workspace>", "<ignore>", "<un_ordered>", "<user_home>"}

// IsIgnore will reture whether the expected result in 'expectedStr' should be ignored.
func IsIgnore(expectedStr string) bool {
	return strings.Contains(expectedStr, "<ignore>")
}

// ReplaceAllKeyByValue will replace all 'key's by 'value' in 'originStr'.
func ReplaceAllKeyByValue(originStr, key, value string) string {
	if !funk.Contains(KEYS, key) {
		reporter.ExitWithReport("unknown key.", key)
	} else {
		return strings.ReplaceAll(originStr, key, value)
	}

	return originStr
}

// SplitCommand will spilt command string into []string,
// but the string in quotes will not be cut.
// If 'command' is 'aaa bbb "ccc ddd"', SplitCommand will return ["aaa", "bbb", "ccc ddd"].
func SplitCommand(command string) []string {
	var args []string
	var currentArg string
	inQuotes := false
	for _, char := range command {
		if char == '"' {
			inQuotes = !inQuotes
			continue
		}
		if char == ' ' && !inQuotes {
			args = append(args, currentArg)
			currentArg = ""
			continue
		}
		currentArg += string(char)
	}
	if currentArg != "" {
		args = append(args, currentArg)
	}
	return args
}

// RemoveLineOrder will remove the line order in 'str1'.
func RemoveLineOrder(str1 string) string {
	// Split the strings into slices of lines
	lines1 := strings.Split(str1, "\n")

	// Sort the slices of lines
	sort.Strings(lines1)

	// Compare the sorted slices of lines
	return strings.Join(lines1, "\n")
}
