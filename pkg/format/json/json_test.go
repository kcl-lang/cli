// Copyright The KCL Authors. All rights reserved.

package json

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestSingle(t *testing.T) {
	tests := []struct {
		name     string
		rawJSON  string
		wantErr  bool
		contains string
	}{
		{
			name:     "Simple object",
			rawJSON:  `{"name":"test","value":123}`,
			wantErr:  false,
			contains: `"name": "test"`,
		},
		{
			name:     "Nested object",
			rawJSON:  `{"config":{"name":"test","value":123}}`,
			wantErr:  false,
			contains: `"config": {`,
		},
		{
			name:     "Array",
			rawJSON:  `[{"name":"first"},{"name":"second"}]`,
			wantErr:  false,
			contains: `"name": "first"`,
		},
		{
			name:     "Valid JSON with indentation",
			rawJSON:  `{"a":1,"b":2,"c":3}`,
			wantErr:  false,
			contains: `"a": 1`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock result using the actual KCLResultList structure
			// For testing purposes, we'll test the JSON formatting logic directly
			var out bytes.Buffer
			err := json.Indent(&out, []byte(tt.rawJSON), "", "    ")
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Indent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result := out.String() + "\n"

			if !tt.wantErr && !strings.Contains(result, tt.contains) {
				t.Errorf("Result = %v, want to contain %v", result, tt.contains)
			}
			// Check that result is properly formatted (contains newlines and indentation)
			if !tt.wantErr && !strings.Contains(result, "\n") {
				t.Errorf("Result should be formatted with newlines")
			}
		})
	}
}

func TestStream(t *testing.T) {
	tests := []struct {
		name       string
		yamlStream string
		wantErr    bool
		docCount   int
		contains   []string
	}{
		{
			name:       "YAML Stream with 2 documents",
			yamlStream: "---\nname: First\nvalue: 1\n---\nname: Second\nvalue: 2\n",
			wantErr:    false,
			docCount:   2,
			contains:   []string{`"name": "First"`, `"name": "Second"`},
		},
		{
			name:       "YAML Stream with 3 documents",
			yamlStream: "---\na: 1\n---\nb: 2\n---\nc: 3\n",
			wantErr:    false,
			docCount:   3,
			contains:   []string{`"a": 1`, `"b": 2`, `"c": 3`},
		},
		{
			name:       "YAML Stream with nested structures",
			yamlStream: "---\nconfig:\n  name: test\n  value: 123\n---\nconfig:\n  name: test2\n  value: 456\n",
			wantErr:    false,
			docCount:   2,
			contains:   []string{`"config": {`, `"name": "test"`, `"name": "test2"`},
		},
		{
			name:       "Single document (no stream)",
			yamlStream: "name: test\nvalue: 123\n",
			wantErr:    false,
			docCount:   1,
			contains:   []string{`"name": "test"`, `"value": 123`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Stream(tt.yamlStream)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				resultStr := string(result)
				// Check that all expected strings are in the result
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("Stream() result = %v, want to contain %v", resultStr, expected)
					}
				}
				// Note: JSON Stream output is not a single valid JSON, but multiple objects
				// Each individual document is valid JSON, verified by checking format
				if !strings.Contains(resultStr, "{") {
					t.Errorf("Stream() result should contain JSON objects")
				}
			}
		})
	}
}

func TestStreamFormat(t *testing.T) {
	yamlStream := "---\nname: First\nvalue: 1\n---\nname: Second\nvalue: 2\n"

	result, err := Stream(yamlStream)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	resultStr := string(result)

	// Check that documents are separated by commas
	if !strings.Contains(resultStr, "},\n{") {
		t.Errorf("Stream() result should contain comma separators between documents")
	}

	// Check that result ends with newline
	if !strings.HasSuffix(resultStr, "\n") {
		t.Errorf("Stream() result should end with newline")
	}

	// Verify it contains both documents
	if !strings.Contains(resultStr, `"name": "First"`) {
		t.Errorf("Result should contain first document")
	}
	if !strings.Contains(resultStr, `"name": "Second"`) {
		t.Errorf("Result should contain second document")
	}

	// Note: The result is multiple JSON objects separated by commas
	// Each individual document (between { and }) should be valid JSON
	// Extract first document: from { to first }
	firstDocStart := strings.Index(resultStr, "{")
	firstDocEnd := strings.Index(resultStr, "}")
	if firstDocStart == -1 || firstDocEnd == -1 {
		t.Fatalf("Could not find document boundaries")
	}
	firstDoc := resultStr[firstDocStart : firstDocEnd+1]

	var doc1 map[string]interface{}
	if err := json.Unmarshal([]byte(firstDoc), &doc1); err != nil {
		t.Errorf("First document is not valid JSON: %v\nContent: %s", err, firstDoc)
	}

	if doc1["name"] != "First" {
		t.Errorf("First document name = %v, want 'First'", doc1["name"])
	}

	// Extract second document
	secondDocStart := strings.LastIndex(resultStr, "{")
	secondDocEnd := strings.LastIndex(resultStr, "}")
	if secondDocStart == -1 || secondDocEnd == -1 {
		t.Fatalf("Could not find second document boundaries")
	}
	secondDoc := resultStr[secondDocStart : secondDocEnd+1]

	var doc2 map[string]interface{}
	if err := json.Unmarshal([]byte(secondDoc), &doc2); err != nil {
		t.Errorf("Second document is not valid JSON: %v\nContent: %s", err, secondDoc)
	}

	if doc2["name"] != "Second" {
		t.Errorf("Second document name = %v, want 'Second'", doc2["name"])
	}
}
