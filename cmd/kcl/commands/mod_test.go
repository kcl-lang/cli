package cmd

import (
	"bytes"
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kpm/pkg/client"
)

func TestModCmdWithSkipTlsVerify(t *testing.T) {
	var buf bytes.Buffer

	mux := gohttp.NewServeMux()
	mux.HandleFunc("/", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		buf.WriteString("Called Success\n")
		fmt.Fprintln(w, "Hello, client")
	})

	mux.HandleFunc("/subpath/tags/list", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		buf.WriteString("Called Success\n")
		fmt.Fprintln(w, "Hello, client")
	})

	mux.HandleFunc("/subpath", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		fmt.Fprintln(w, "Hello from subpath")
	})

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	fmt.Printf("ts.URL: %v\n", ts.URL)
	turl, err := url.Parse(ts.URL)
	assert.Equal(t, err, nil)
	turl.Scheme = "oci"
	turl.Path = filepath.Join(turl.Path, "subpath")
	fmt.Printf("turl.String(): %v\n", turl.String())

	kpmcli, err := client.NewKpmClient()
	assert.Equal(t, err, nil)

	originalDir, err := os.Getwd()
	assert.NoError(t, err)
	testRootDir := filepath.Join(originalDir, "test_data")

	runTest := func(testDir string, testFunc func(), beforeTestFuncs []func()) {
		if testDir != "" {
			err = os.Chdir(filepath.Join(testRootDir, testDir))
			assert.NoError(t, err)
		}

		for _, beforeTestFunc := range beforeTestFuncs {
			beforeTestFunc()
		}

		testFunc()
		assert.Equal(t, buf.String(), "Called Success\n")
		buf.Reset()
		defer func() {
			err := os.Chdir(originalDir)
			assert.NoError(t, err)
		}()
	}

	genKclModWithDep := func(depUrl string) {
		fmt.Println("Executing extra function for test_mod_graph")
		kclModContent := fmt.Sprintf(`[package]
name = "test_mod"
edition = "v0.10.0"
version = "0.0.1"

[dependencies]
dep1 = { oci = "%s"}
	`, depUrl)
		err := os.WriteFile("kcl.mod", []byte(kclModContent), 0644)
		assert.NoError(t, err)
	}

	runTest("", func() {
		fmt.Println("test_mod_pull")
		cmd := NewModPullCmd(kpmcli)
		cmd.SetArgs([]string{turl.String(), "--insecure-skip-tls-verify"})
		_ = cmd.Execute()
	}, []func(){})

	runTest("test_mod_push", func() {
		fmt.Println("test_mod_push")
		cmd := NewModPushCmd(kpmcli)
		cmd.SetArgs([]string{turl.String(), "--insecure-skip-tls-verify"})
		_ = cmd.Execute()
	}, []func(){})

	runTest("test_mod_add", func() {
		fmt.Println("test_mod_add")
		cmd := NewModAddCmd(kpmcli)
		cmd.SetArgs([]string{turl.String(), "--insecure-skip-tls-verify"})
		_ = cmd.Execute()
	}, []func(){})

	runTest("test_mod_graph", func() {
		fmt.Println("test_mod_graph")
		cmd := NewModGraphCmd(kpmcli)
		cmd.SetArgs([]string{"--insecure-skip-tls-verify"})
		_ = cmd.Execute()
	}, []func(){
		func() { genKclModWithDep(turl.String()) },
	})

	runTest("test_mod_metadata", func() {
		fmt.Println("test_mod_metadata")
		cmd := NewModMetadataCmd(kpmcli)
		cmd.SetArgs([]string{"--update", "--insecure-skip-tls-verify"})
		_ = cmd.Execute()
	}, []func(){
		func() { genKclModWithDep(turl.String()) },
	})

	runTest("test_mod_update", func() {
		fmt.Println("test_mod_update")
		cmd := NewModUpdateCmd(kpmcli)
		cmd.SetArgs([]string{"--insecure-skip-tls-verify"})
		_ = cmd.Execute()
	}, []func(){
		func() { genKclModWithDep(turl.String()) },
	})
}
