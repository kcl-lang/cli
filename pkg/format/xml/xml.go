// Copyright The KCL Authors. All rights reserved.

package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/goccy/go-yaml"
	yamlformat "kcl-lang.io/cli/pkg/format/yaml"
)

// Convert converts arbitrary data structures to XML format with a root element.
func Convert(data any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("<root>")
	if err := encode(&buf, data, ""); err != nil {
		return nil, err
	}
	buf.WriteString("</root>")
	return buf.Bytes(), nil
}

// Single converts a single YAML result to XML format.
func Single(yamlResult string) ([]byte, error) {
	var yamlData any
	if err := yaml.UnmarshalWithOptions([]byte(yamlResult), &yamlData); err != nil {
		return nil, err
	}
	out, err := Convert(yamlData)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(out) + "\n"), nil
}

// Stream converts a YAML Stream to XML format with multiple root elements.
func Stream(yamlResult string) ([]byte, error) {
	docs, err := yamlformat.ParseStream(yamlResult)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	out.WriteString(xml.Header)
	out.WriteString("<results>\n")
	for _, doc := range docs {
		xmlData, err := Convert(doc)
		if err != nil {
			return nil, err
		}
		// Remove the XML header and wrap root element
		xmlStr := string(xmlData)
		out.WriteString("  ")
		out.WriteString(xmlStr)
		out.WriteString("\n")
	}
	out.WriteString("</results>\n")
	return out.Bytes(), nil
}

// encode recursively encodes data structures to XML.
func encode(buf *bytes.Buffer, data any, defaultKey string) error {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			if err := encodeElement(buf, key, value); err != nil {
				return err
			}
		}
	case map[any]any:
		for key, value := range v {
			keyStr := fmt.Sprintf("%v", key)
			if err := encodeElement(buf, keyStr, value); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range v {
			if defaultKey == "" {
				defaultKey = "item"
			}
			if err := encodeElement(buf, defaultKey, item); err != nil {
				return err
			}
		}
	case string:
		buf.WriteString(escapeString(v))
	case int, int64, float64, bool:
		buf.WriteString(fmt.Sprintf("%v", v))
	case nil:
		// Skip nil values
	default:
		buf.WriteString(fmt.Sprintf("%v", v))
	}
	return nil
}

// encodeElement encodes a single XML element.
func encodeElement(buf *bytes.Buffer, key string, value any) error {
	buf.WriteString("<")
	buf.WriteString(key)
	buf.WriteString(">")

	if err := encode(buf, value, ""); err != nil {
		return err
	}

	buf.WriteString("</")
	buf.WriteString(key)
	buf.WriteString(">")
	return nil
}

// escapeString escapes special XML characters in a string.
func escapeString(s string) string {
	var buf bytes.Buffer
	xml.Escape(&buf, []byte(s))
	return buf.String()
}
