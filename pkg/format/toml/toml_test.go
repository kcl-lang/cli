// Copyright The KCL Authors. All rights reserved.

package toml

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestSingle(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		sortKeys bool
		wantErr  bool
		contains []string
	}{
		{
			name:     "Simple object without sorting",
			yaml:     "name: test\nvalue: 123\n",
			sortKeys: false,
			wantErr:  false,
			contains: []string{"name", "test", "value", "123"},
		},
		{
			name:     "Simple object with sorting",
			yaml:     "name: test\nvalue: 123\n",
			sortKeys: true,
			wantErr:  false,
			contains: []string{"name", "test", "value", "123"},
		},
		{
			name:     "Nested structure",
			yaml:     "config:\n  name: test\n  value: 123\n",
			sortKeys: false,
			wantErr:  false,
			contains: []string{"[config]", "name", "value"},
		},
		{
			name:     "Array-like structure",
			yaml:     "items:\n  - name: first\n  - name: second\n",
			sortKeys: false,
			wantErr:  false,
			contains: []string{"items", "name", "first", "second"},
		},
		{
			name:     "Complex nested structure",
			yaml:     "server:\n  host: localhost\n  port: 8080\ndatabase:\n  host: db.example.com\n  port: 5432\n",
			sortKeys: false,
			wantErr:  false,
			contains: []string{"[server]", "[database]", "host", "port"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Single(tt.yaml, tt.sortKeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Single() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				resultStr := string(result)
				// Check that all expected strings are in the result
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("Single() result = %v, want to contain %v", resultStr, expected)
					}
				}
				// Check that result ends with newline
				if !strings.HasSuffix(resultStr, "\n") {
					t.Errorf("Single() result should end with newline")
				}
				// Verify it's valid TOML by parsing it
				var parsed interface{}
				if _, err := toml.Decode(resultStr, &parsed); err != nil {
					t.Errorf("Single() result is not valid TOML: %v", err)
				}
			}
		})
	}
}

func TestStream(t *testing.T) {
	tests := []struct {
		name       string
		yamlStream string
		sortKeys   bool
		wantErr    bool
		docCount   int
		contains   []string
	}{
		{
			name:       "YAML Stream with 2 documents",
			yamlStream: "---\nname: First\nvalue: 1\n---\nname: Second\nvalue: 2\n",
			sortKeys:   false,
			wantErr:    false,
			docCount:   2,
			contains:   []string{"name", "First", "Second", "# --- Document separator ---"},
		},
		{
			name:       "YAML Stream with 3 documents",
			yamlStream: "---\na: 1\n---\nb: 2\n---\nc: 3\n",
			sortKeys:   false,
			wantErr:    false,
			docCount:   3,
			contains:   []string{"a = 1", "b = 2", "c = 3", "# --- Document separator ---"},
		},
		{
			name:       "YAML Stream with nested structures",
			yamlStream: "---\nconfig:\n  name: test1\n---\nconfig:\n  name: test2\n",
			sortKeys:   false,
			wantErr:    false,
			docCount:   2,
			contains:   []string{"[config]", "test1", "test2", "# --- Document separator ---"},
		},
		{
			name:       "Single document (no stream)",
			yamlStream: "name: test\nvalue: 123\n",
			sortKeys:   false,
			wantErr:    false,
			docCount:   1,
			contains:   []string{"name", "test", "value", "123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Stream(tt.yamlStream, tt.sortKeys)
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
				// For stream with multiple documents, check separator count
				if tt.docCount > 1 {
					sepCount := strings.Count(resultStr, "# --- Document separator ---")
					expectedSepCount := tt.docCount - 1
					if sepCount != expectedSepCount {
						t.Errorf("Stream() should contain %d separators, got %d", expectedSepCount, sepCount)
					}
				}
			}
		})
	}
}

func TestStreamSortedKeys(t *testing.T) {
	yamlStream := "---\nz: 1\na: 2\nm: 3\n---\nz: 4\na: 5\nm: 6\n"

	// Test with sorted keys - note that TOML stream output has document separators
	// Each document should be valid TOML independently
	result, err := Stream(yamlStream, true)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	resultStr := string(result)

	// Check that document separators are present
	if !strings.Contains(resultStr, "# --- Document separator ---") {
		t.Errorf("Stream() result should contain document separators")
	}

	// Verify each document segment is valid TOML by splitting and checking
	docs := strings.Split(resultStr, "# --- Document separator ---")
	for i, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}
		var parsed interface{}
		if _, err := toml.Decode(doc, &parsed); err != nil {
			t.Logf("Document %d: %s", i, doc)
			t.Errorf("Stream() document %d is not valid TOML: %v", i, err)
		}
	}
}

func TestSingleIntegration(t *testing.T) {
	yaml := "name: test\nvalue: 123\nenabled: true\n"

	result, err := Single(yaml, false)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}

	// Parse the result to verify it's valid TOML
	var parsed map[string]interface{}
	if _, err := toml.Decode(string(result), &parsed); err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	// Verify values
	if parsed["name"] != "test" {
		t.Errorf("name = %v, want 'test'", parsed["name"])
	}
	if parsed["value"] != int64(123) {
		t.Errorf("value = %v, want 123", parsed["value"])
	}
	if parsed["enabled"] != true {
		t.Errorf("enabled = %v, want true", parsed["enabled"])
	}
}
