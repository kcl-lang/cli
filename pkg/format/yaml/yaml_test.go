// Copyright The KCL Authors. All rights reserved.

package yaml

import (
	"testing"
)

func TestIsStream(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected bool
	}{
		{
			name:     "YAML Stream with multiple documents",
			yaml:     "---\na: 1\n---\nb: 2\n",
			expected: true,
		},
		{
			name:     "YAML Stream at beginning",
			yaml:     "---\nname: test\nvalue: 123\n",
			expected: true,
		},
		{
			name:     "Single YAML document",
			yaml:     "name: test\nvalue: 123\n",
			expected: false,
		},
		{
			name:     "Empty string",
			yaml:     "",
			expected: false,
		},
		{
			name:     "YAML with nested structure",
			yaml:     "config:\n  name: test\n  value: 123\n",
			expected: false,
		},
		{
			name:     "Multiple documents with proper separators",
			yaml:     "---\ndoc1: value1\n---\ndoc2: value2\n---\ndoc3: value3\n",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStream(tt.yaml)
			if result != tt.expected {
				t.Errorf("IsStream() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseStream(t *testing.T) {
	tests := []struct {
		name         string
		yaml         string
		wantDocCount int
		wantErr      bool
	}{
		{
			name:         "Single document",
			yaml:         "name: test\nvalue: 123\n",
			wantDocCount: 1,
			wantErr:      false,
		},
		{
			name:         "YAML Stream with 2 documents",
			yaml:         "---\na: 1\n---\nb: 2\n",
			wantDocCount: 2,
			wantErr:      false,
		},
		{
			name:         "YAML Stream with 3 documents",
			yaml:         "---\nname: First\n---\nname: Second\n---\nname: Third\n",
			wantDocCount: 3,
			wantErr:      false,
		},
		{
			name:         "Complex nested structures",
			yaml:         "---\nconfig:\n  name: test\n  value: 123\n---\nconfig:\n  name: test2\n  value: 456\n",
			wantDocCount: 2,
			wantErr:      false,
		},
		{
			name:         "Empty YAML",
			yaml:         "",
			wantDocCount: 1,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, err := ParseStream(tt.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(docs) != tt.wantDocCount {
				t.Errorf("ParseStream() returned %d documents, want %d", len(docs), tt.wantDocCount)
			}
		})
	}
}

func TestParseStreamDocumentContent(t *testing.T) {
	yamlStream := "---\nname: First\nvalue: 1\n---\nname: Second\nvalue: 2\n"

	docs, err := ParseStream(yamlStream)
	if err != nil {
		t.Fatalf("ParseStream() error = %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("ParseStream() returned %d documents, want 2", len(docs))
	}

	// Check first document
	doc1, ok1 := docs[0].(map[string]interface{})
	if !ok1 {
		t.Fatalf("First document is not a map")
	}
	if doc1["name"] != "First" {
		t.Errorf("First document name = %v, want 'First'", doc1["name"])
	}
	// Check value - accept both int and uint types
	value1, err1 := getAsInt64(doc1["value"])
	if !err1 || value1 != 1 {
		t.Errorf("First document value = %v (type %T), want 1", doc1["value"], doc1["value"])
	}

	// Check second document
	doc2, ok2 := docs[1].(map[string]interface{})
	if !ok2 {
		t.Fatalf("Second document is not a map")
	}
	if doc2["name"] != "Second" {
		t.Errorf("Second document name = %v, want 'Second'", doc2["name"])
	}
	value2, err2 := getAsInt64(doc2["value"])
	if !err2 || value2 != 2 {
		t.Errorf("Second document value = %v (type %T), want 2", doc2["value"], doc2["value"])
	}
}

// Helper function to extract int64 from various numeric types
func getAsInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case uint:
		return int64(val), true
	case uint64:
		return int64(val), true
	case float64:
		return int64(val), true
	default:
		return 0, false
	}
}
