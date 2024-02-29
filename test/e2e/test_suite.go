package e2e

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/utils"
)

const TEST_SUITES_DIR = "test_suites"

const ENV = "env"
const CONF = "conf.json"
const STDOUT = "stdout"
const STDERR = "stderr"
const INPUT = "input"
const IGNORE = "ignore"

type TestConf struct {
	Cwd string
}

type TestSuite struct {
	Name         string
	Envs         string
	Input        string
	ExpectStdout string
	ExpectStderr string
	ignore       bool
	Conf         TestConf
}

// LoadTestSuite load test suite from 'getWorkDir()/test_suites/name'.
func LoadTestSuite(testSuitePath, name string) *TestSuite {
	reporter.Report("testing: ", name)

	confPath := filepath.Join(testSuitePath, CONF)
	var conf TestConf

	if utils.DirExists(confPath) {
		err := json.Unmarshal([]byte(LoadFile(confPath)), &conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	if conf.Cwd == "" {
		conf.Cwd = filepath.Join(testSuitePath, "test_space")
	} else {
		conf.Cwd = filepath.Join(testSuitePath, conf.Cwd)
	}

	ts := TestSuite{
		Name:         name,
		ExpectStdout: LoadFile(filepath.Join(testSuitePath, STDOUT)),
		ExpectStderr: LoadFile(filepath.Join(testSuitePath, STDERR)),
		Input:        LoadFile(filepath.Join(testSuitePath, INPUT)),
		// Envs:         LoadFile(filepath.Join(testSuitePath, ENV)),
		ignore: utils.DirExists(filepath.Join(testSuitePath, IGNORE)),
		Conf:   conf,
	}
	return ts.ReplaceStringVar(testSuitePath)
}

func (ts *TestSuite) ReplaceStringVar(testSuitePath string) *TestSuite {
	ts.ExpectStdout = strings.ReplaceAll(ts.ExpectStdout, "<workspace>", testSuitePath)
	ts.ExpectStderr = strings.ReplaceAll(ts.ExpectStderr, "<workspace>", testSuitePath)
	ts.Input = strings.ReplaceAll(ts.Input, "<workspace>", testSuitePath)
	return ts
}

// LoadFile will read the file from 'path' and return the content.
func LoadFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

// LoadAllTestSuites load all test suites from 'getWorkDir()/test_suites'.
func LoadAllTestSuites(testSuitesDir string) []TestSuite {
	testSuites := make([]TestSuite, 0)
	files, err := os.ReadDir(testSuitesDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		ts := LoadTestSuite(filepath.Join(testSuitesDir, file.Name()), file.Name())
		if file.IsDir() && ts != nil {
			testSuites = append(
				testSuites,
				*ts,
			)
		}
	}

	return testSuites
}

// GetTestSuiteInfo return a info for a test suite "<name>:<info>:<env>"
func (ts *TestSuite) GetTestSuiteInfo() string {
	return fmt.Sprintf("%s:%s", ts.Name, strings.ReplaceAll(ts.Envs, "\n", ":"))
}
