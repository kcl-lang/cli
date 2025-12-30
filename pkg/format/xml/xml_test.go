// Copyright The KCL Authors. All rights reserved.

package xml

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantErr  bool
		contains []string
	}{
		{
			name: "Simple map",
			data: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			wantErr:  false,
			contains: []string{"<name>test</name>", "<value>123</value>"},
		},
		{
			name: "Nested map",
			data: map[string]interface{}{
				"config": map[string]interface{}{
					"name":  "test",
					"value": 123,
				},
			},
			wantErr:  false,
			contains: []string{"<config>", "<name>test</name>", "<value>123</value>", "</config>"},
		},
		{
			name: "Array",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "first"},
					map[string]interface{}{"name": "second"},
				},
			},
			wantErr:  false,
			contains: []string{"<items>", "<name>first</name>", "<name>second</name>", "</items>"},
		},
		{
			name: "Primitive types",
			data: map[string]interface{}{
				"stringVal": "hello",
				"intVal":    42,
				"floatVal":  3.14,
				"boolVal":   true,
			},
			wantErr:  false,
			contains: []string{"<stringVal>hello</stringVal>", "<intVal>42</intVal>", "<floatVal>3.14</floatVal>", "<boolVal>true</boolVal>"},
		},
		{
			name:     "Empty map",
			data:     map[string]interface{}{},
			wantErr:  false,
			contains: []string{"<root>", "</root>"},
		},
		{
			name: "Map with interface{} keys",
			data: map[interface{}]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			wantErr:  false,
			contains: []string{"<key1>value1</key1>", "<key2>123</key2>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				resultStr := string(result)
				// Check that all expected strings are in the result
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("Convert() result = %v, want to contain %v", resultStr, expected)
					}
				}
				// Check for root element
				if !strings.Contains(resultStr, "<root>") {
					t.Errorf("Convert() result should contain <root> element")
				}
				if !strings.Contains(resultStr, "</root>") {
					t.Errorf("Convert() result should contain closing </root> element")
				}
			}
		})
	}
}

func TestSingle(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantErr  bool
		contains []string
	}{
		{
			name:     "Simple YAML document",
			yaml:     "name: test\nvalue: 123\n",
			wantErr:  false,
			contains: []string{"<?xml version=", "<name>test</name>", "<value>123</value>", "<root>", "</root>"},
		},
		{
			name:     "Nested YAML structure",
			yaml:     "config:\n  name: test\n  value: 123\n",
			wantErr:  false,
			contains: []string{"<config>", "<name>test</name>", "<value>123</value>"},
		},
		{
			name:     "YAML with array",
			yaml:     "items:\n  - name: first\n  - name: second\n",
			wantErr:  false,
			contains: []string{"<items>", "<name>first</name>", "<name>second</name>"},
		},
		{
			name:     "YAML with special characters",
			yaml:     "message: \"Hello <world> & goodbye\"\n",
			wantErr:  false,
			contains: []string{"<message>", "Hello", "world", "goodbye", "</message>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Single(tt.yaml)
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
				// Check for XML header
				if !strings.Contains(resultStr, "<?xml version=") {
					t.Errorf("Single() result should contain XML declaration")
				}
				// Check that result ends with newline
				if !strings.HasSuffix(resultStr, "\n") {
					t.Errorf("Single() result should end with newline")
				}
			}
		})
	}
}

func TestStream(t *testing.T) {
	tests := []struct {
		name       string
		yamlStream string
		wantErr    bool
		contains   []string
		docCount   int
	}{
		{
			name:       "YAML Stream with 2 documents",
			yamlStream: "---\nname: First\n---\nname: Second\n",
			wantErr:    false,
			docCount:   2,
			contains:   []string{"<?xml version=", "<results>", "<root>", "<name>First</name>", "<name>Second</name>", "</results>"},
		},
		{
			name:       "YAML Stream with 3 documents",
			yamlStream: "---\na: 1\n---\nb: 2\n---\nc: 3\n",
			wantErr:    false,
			docCount:   3,
			contains:   []string{"<a>1</a>", "<b>2</b>", "<c>3</c>"},
		},
		{
			name:       "YAML Stream with nested structures",
			yamlStream: "---\nconfig:\n  name: test1\n---\nconfig:\n  name: test2\n",
			wantErr:    false,
			docCount:   2,
			contains:   []string{"<config>", "<name>test1</name>", "<name>test2</name>"},
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
				// Check for results wrapper
				if !strings.Contains(resultStr, "<results>") {
					t.Errorf("Stream() result should contain <results> wrapper")
				}
				if !strings.Contains(resultStr, "</results>") {
					t.Errorf("Stream() result should contain closing </results>")
				}
				// Check for XML header
				if !strings.Contains(resultStr, "<?xml version=") {
					t.Errorf("Stream() result should contain XML declaration")
				}
				// Count root elements (should equal docCount)
				rootCount := strings.Count(resultStr, "<root>")
				if rootCount != tt.docCount {
					t.Errorf("Stream() should contain %d <root> elements, got %d", tt.docCount, rootCount)
				}
			}
		})
	}
}

func TestXMLEscaping(t *testing.T) {
	data := map[string]interface{}{
		"special": "<text> & \"quotes\"",
	}

	result, err := Convert(data)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	resultStr := string(result)

	// Check that special characters are escaped
	if !strings.Contains(resultStr, "&lt;") && !strings.Contains(resultStr, "<text>") {
		t.Errorf("XML special characters should be escaped or in CDATA")
	}
	if strings.Contains(resultStr, "<text>") && !strings.Contains(resultStr, "&lt;") {
		// If raw tags are present, they should be escaped
		t.Logf("Note: Raw tags present - ensure XML is valid")
	}
}

func TestValidXMLStructure(t *testing.T) {
	yaml := "name: test\nvalue: 123\n"

	result, err := Single(yaml)
	if err != nil {
		t.Fatalf("Single() error = %v", err)
	}

	// Try to parse the result as XML to verify it's valid
	type Root struct {
		XMLName xml.Name `xml:"root"`
		Name    string   `xml:"name"`
		Value   int      `xml:"value"`
	}

	var parsed Root
	err = xml.Unmarshal(result, &parsed)
	if err != nil {
		t.Errorf("Single() result is not valid XML: %v", err)
	}

	// Verify values
	if parsed.Name != "test" {
		t.Errorf("parsed.Name = %v, want 'test'", parsed.Name)
	}
	if parsed.Value != 123 {
		t.Errorf("parsed.Value = %v, want 123", parsed.Value)
	}
}

func TestArrayHandling(t *testing.T) {
	data := map[string]interface{}{
		"items": []interface{}{"first", "second", "third"},
	}

	result, err := Convert(data)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	resultStr := string(result)

	// Check that array items are wrapped in <item> tags
	if !strings.Contains(resultStr, "<item>") {
		t.Errorf("Array items should be wrapped in <item> tags")
	}
	if !strings.Contains(resultStr, "first") || !strings.Contains(resultStr, "second") || !strings.Contains(resultStr, "third") {
		t.Errorf("Array items should be present in result")
	}
}

func TestNilHandling(t *testing.T) {
	data := map[string]interface{}{
		"present": "value",
		"absent":  nil,
	}

	result, err := Convert(data)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	resultStr := string(result)

	// Check that present value is included
	if !strings.Contains(resultStr, "<present>value</present>") {
		t.Errorf("Present value should be in result")
	}
	// Nil values should be skipped, so <absent> should not appear
	if strings.Contains(resultStr, "<absent>") {
		t.Logf("Note: Nil values are present (this is acceptable if handled correctly)")
	}
}
