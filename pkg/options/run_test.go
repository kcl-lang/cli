package options

import (
	"bytes"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRunOptions_Run(t *testing.T) {
	options := NewRunOptions()
	options.Entries = []string{"./testdata/run/kubernetes.k"}

	// test yaml output
	var buf1 bytes.Buffer
	options.Writer = &buf1
	options.Format = Yaml
	options.SortKeys = true

	err := options.Run()
	if err != nil {
		t.Errorf("RunOptions.Run() failed: %v", err)
	}

	expectedOutput := `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:1.14.2
        name: nginx
        ports:
        - containerPort: 80
`

	if got := buf1.String(); got != expectedOutput {
		t.Errorf("\nexpected: %s\ngot: %s", expectedOutput, got)
	}

	// test json output
	var buf2 bytes.Buffer
	options.Writer = &buf2
	options.Format = Json
	options.SortKeys = true

	err = options.Run()
	if err != nil {
		t.Errorf("RunOptions.Run() failed: %v", err)
	}

	expectedOutput = `{
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
        "labels": {
            "app": "nginx"
        },
        "name": "nginx"
    },
    "spec": {
        "replicas": 3,
        "selector": {
            "matchLabels": {
                "app": "nginx"
            }
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "nginx"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "nginx:1.14.2",
                        "name": "nginx",
                        "ports": [
                            {
                                "containerPort": 80
                            }
                        ]
                    }
                ]
            }
        }
    }
}
`

	if got := buf2.String(); got != expectedOutput {
		t.Errorf("\nexpected: %s\ngot: %s", expectedOutput, got)
	}
}

func TestRunOptions_Complete(t *testing.T) {
	options := NewRunOptions()
	args := []string{"file1.k", "file2.k", "file3.k"}

	err := options.Complete(args)
	if err != nil {
		t.Errorf("RunOptions.Complete() failed: %v", err)
	}

	expectedEntries := []string{"file1.k", "file2.k", "file3.k"}

	if len(options.Entries) != len(expectedEntries) {
		t.Fatalf("unexpected number of entries:\nexpected: %d\ngot: %d", len(expectedEntries), len(options.Entries))
	}

	for i := range options.Entries {
		if options.Entries[i] != expectedEntries[i] {
			t.Errorf("unexpected entry at index %d:\nexpected: %s\ngot: %s", i, expectedEntries[i], options.Entries[i])
		}
	}
}

func TestRunOptions_Validate(t *testing.T) {
	options := NewRunOptions()
	options.Format = "invalid_format"

	err := options.Validate()
	if err == nil {
		t.Errorf("RunOptions.Validate() did not return an error")
	} else {
		expectedError := "invalid output format, expected [json yaml toml], got invalid_format"
		if err.Error() != expectedError {
			t.Errorf("unexpected error message:\nexpected: %s\ngot: %s", expectedError, err.Error())
		}
	}
}

func TestRunPkg(t *testing.T) {
	var buf bytes.Buffer
	options := NewRunOptions()
	options.Entries = []string{"./testdata/run_pkg"}
	options.Writer = &buf

	err := options.Run()
	assert.Equal(t, err, nil)
	resStr := strings.ReplaceAll(buf.String(), "\r\n", "")
	resStr = strings.ReplaceAll(resStr, "\n", "")

	assert.Equal(t, resStr, "The_first_kcl_program: Hello World!")
}
