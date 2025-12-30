// Copyright The KCL Authors. All rights reserved.

package json

import (
	"bytes"
	"encoding/json"

	"kcl-lang.io/cli/pkg/format/yaml"
	"kcl-lang.io/kcl-go/pkg/kcl"
)

// Single converts a single KCL result to JSON format.
func Single(result *kcl.KCLResultList) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(result.GetRawJsonResult()), "", "    ")
	if err != nil {
		return nil, err
	}
	return []byte(out.String() + "\n"), nil
}

// Stream converts a YAML Stream to JSON format.
func Stream(yamlResult string) ([]byte, error) {
	docs, err := yaml.ParseStream(yamlResult)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	for i, doc := range docs {
		jsonData, err := json.MarshalIndent(doc, "", "    ")
		if err != nil {
			return nil, err
		}
		out.Write(jsonData)
		if i < len(docs)-1 {
			out.WriteString(",\n")
		} else {
			out.WriteString("\n")
		}
	}
	return out.Bytes(), nil
}
