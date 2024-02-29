package e2e

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gexec"
	"kcl-lang.io/kpm/pkg/reporter"
)

const TEST_WORKSPACE = "test_kpm_workspace"

// CreateTestWorkspace create an empty dir "($pwd)/test_kpm_workspace" for test kpm cli.
func CreateTestWorkspace() string {
	workspacePath := filepath.Join(GetWorkDir(), TEST_WORKSPACE)
	err := os.MkdirAll(workspacePath, 0755)
	if err != nil {
		reporter.ExitWithReport("failed to create workspace.")
	}
	return workspacePath
}

// CleanUpTestWorkspace will do 'rm -rf' to the "($pwd)/test_kpm_workspace".
func CleanUpTestWorkspace() string {
	workspacePath := filepath.Join(GetWorkDir(), TEST_WORKSPACE)
	err := os.RemoveAll(workspacePath)
	if err != nil {
		reporter.ExitWithReport("failed to clean up workspace.")
	}
	return workspacePath
}

// GetWorkspace return the path of test workspace.
func GetWorkspace() string {
	return filepath.Join(GetWorkDir(), TEST_WORKSPACE)
}

// GetWrokDir return work directory
func GetWorkDir() string {
	dir, err := os.Getwd()
	if err != nil {
		reporter.ExitWithReport("failed to load workspace.")
	}
	return dir
}

// GetKpmCLIBin return kusion binary path in e2e test
func GetKpmCLIBin() string {
	dir, _ := os.Getwd()
	binPath := filepath.Join(dir, "../..", "bin")
	return binPath
}

// Exec execute common command
func Exec(cli string) (string, string, error) {
	var output []byte
	c := strings.Fields(cli)
	command := exec.Command(c[0], c[1:]...)
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return string(output), string(output), err
	}
	s := session.Wait(300 * time.Second)
	return string(s.Out.Contents()), string(s.Err.Contents()), nil
}

// Exec execute common command
func ExecWithWorkDir(cli, dir string) (string, error) {
	var output []byte
	c := strings.Fields(cli)
	command := exec.Command(c[0], c[1:]...)
	command.Dir = dir
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return string(output), err
	}
	s := session.Wait(300 * time.Second)
	return string(s.Out.Contents()) + string(s.Err.Contents()), nil
}

// ExecKpm executes kusion command
func ExecKpm(cli string) (string, error) {
	var output []byte
	c := strings.Fields(cli)
	commandName := filepath.Join(GetKpmCLIBin(), c[0])
	command := exec.Command(commandName, c[1:]...)
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return string(output), err
	}
	s := session.Wait(300 * time.Second)
	return string(s.Out.Contents()) + string(s.Err.Contents()), nil
}

// ExecKpmWithInDirWithEnv executes kpm command in work directory with env.
func ExecKpmWithInDirWithEnv(cli, dir, env string) (string, string, error) {
	var output []byte
	c := strings.Fields(cli)
	commandName := filepath.Join(GetKpmCLIBin(), c[0])
	command := exec.Command(commandName, c[1:]...)
	command.Dir = dir

	// Set the env
	envLines := strings.FieldsFunc(env, func(c rune) bool { return c == '\n' })
	command.Env = append(command.Env, envLines...)
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return string(output), string(output), err
	}
	s := session.Wait(300 * time.Second)
	return string(s.Out.Contents()), string(s.Err.Contents()), nil
}

// ExecKpmWithWorkDir executes kpm command in work directory
func ExecKpmWithWorkDir(cli, dir string) (string, string, error) {
	var output []byte
	c := SplitCommand(cli)
	commandName := filepath.Join(GetKpmCLIBin(), c[0])
	command := exec.Command(commandName, c[1:]...)
	command.Dir = dir
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return string(output), string(output), err
	}
	s := session.Wait(300 * time.Second)
	return string(s.Out.Contents()), string(s.Err.Contents()), nil
}

// ExecKpmWithStdin executes kpm command in work directory with stdin
func ExecKpmWithStdin(cli, dir, input string) (string, error) {
	var output []byte
	c := strings.Fields(cli)
	commandName := filepath.Join(GetKpmCLIBin(), c[0])
	command := exec.Command(commandName, c[1:]...)
	command.Dir = dir
	subStdin, err := command.StdinPipe()
	if err != nil {
		return string(output), err
	}
	_, err = io.WriteString(subStdin, input)
	if err != nil {
		return string(output), err
	}
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return string(output), err
	}
	s := session.Wait(300 * time.Second)
	return string(s.Out.Contents()) + string(s.Err.Contents()), nil
}
