package cmd

import (
	"bytes"
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kpm/pkg/client"
)

func TestLoginCmdWithSkipTlsVerify(t *testing.T) {
	var buf bytes.Buffer

	mux := gohttp.NewServeMux()
	mux.HandleFunc("/", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		buf.WriteString("Called Success\n")
		fmt.Fprintln(w, "Hello, client")
	})

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	fmt.Printf("ts.URL: %v\n", ts.URL)
	turl, err := url.Parse(ts.URL)
	assert.Equal(t, err, nil)
	turl.Path = filepath.Join(turl.Path, "subpath")
	fmt.Printf("turl.String(): %v\n", turl.String())

	cli, err := client.NewKpmClient()
	assert.Equal(t, err, nil)
	cmd := NewRegistryLoginCmd(cli)
	cmd.SetArgs([]string{turl.String(), "--username=test-user", "--password=test-pass", "--insecure-skip-tls-verify"})
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, buf.String(), "Called Success\nCalled Success\n")
}
