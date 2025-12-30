// Copyright The KCL Authors. All rights reserved.

package toml

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	yamlformat "kcl-lang.io/cli/pkg/format/yaml"
	"kcl-lang.io/kcl-go/pkg/3rdparty/toml"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
)

// Single converts a single KCL result to TOML format.
func Single(yamlResult string, sortKeys bool) ([]byte, error) {
	var out []byte
	var err error
	if sortKeys {
		yamlData := make(map[string]any)
		if err := yaml.UnmarshalWithOptions([]byte(yamlResult), &yamlData); err != nil {
			return nil, err
		}
		out, err = toml.Marshal(&yamlData)
	} else {
		yamlData := &yaml.MapSlice{}
		if err := yaml.UnmarshalWithOptions([]byte(yamlResult), yamlData, yaml.UseOrderedMap()); err != nil {
			return nil, err
		}
		out, err = gen.MarshalTOML(yamlData)
	}
	if err != nil {
		return nil, err
	}
	return []byte(string(out) + "\n"), nil
}

// Stream converts a YAML Stream to TOML format with document separators.
func Stream(yamlResult string, sortKeys bool) ([]byte, error) {
	docs, err := yamlformat.ParseStream(yamlResult)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	for i, doc := range docs {
		var tomlData []byte
		var err error
		if sortKeys {
			if data, ok := doc.(map[string]any); ok {
				tomlData, err = toml.Marshal(&data)
			} else {
				return nil, fmt.Errorf("document %d is not a map", i+1)
			}
		} else {
			// Convert to MapSlice for ordered output
			yamlBytes, err := yaml.Marshal(doc)
			if err != nil {
				return nil, err
			}
			yamlData := &yaml.MapSlice{}
			if err := yaml.UnmarshalWithOptions(yamlBytes, yamlData, yaml.UseOrderedMap()); err != nil {
				return nil, err
			}
			tomlData, err = gen.MarshalTOML(yamlData)
		}
		if err != nil {
			return nil, err
		}
		out.Write(tomlData)
		if i < len(docs)-1 {
			out.WriteString("\n# --- Document separator ---\n\n")
		}
	}
	return out.Bytes(), nil
}
