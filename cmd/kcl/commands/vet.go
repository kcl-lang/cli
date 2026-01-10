// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
	"kcl-lang.io/kcl-go/pkg/tools/validate"
)

const (
	vetDesc = `This command validates the data file using the kcl code.
`
	vetExample = `  # Validate the JSON data using the kcl code
  kcl vet data.json code.k

  # Validate the YAML data using the kcl code
  kcl vet data.yaml code.k --format yaml

  # Validate the JSON data using the kcl code with the schema name
  kcl vet data.json code.k -s Schema

  # Validate and output results as JSON for CI/CD integration
  kcl vet data.json code.k --output json`
)

// VetOptions holds the options for the vet command.
type VetOptions struct {
	validate.ValidateOptions
	// Output specifies the output format: "text" (default) or "json"
	Output string
}

// VetResult represents a structured validation result for JSON output.
type VetResult struct {
	Success  bool       `json:"success"`
	ErrCount int        `json:"errCount,omitempty"`
	Errors   []VetError `json:"errors,omitempty"`
	Message  string     `json:"message,omitempty"`
}

// VetError represents a single validation error in structured format.
type VetError struct {
	ErrorType   string       `json:"errorType,omitempty"`
	File        string       `json:"file,omitempty"`
	Line        int          `json:"line,omitempty"`
	Column      int          `json:"column,omitempty"`
	Message     string       `json:"message,omitempty"`
	CodeSnippet string       `json:"codeSnippet,omitempty"`
	Schema      *SchemaError `json:"schema,omitempty"`
}

// SchemaError represents schema-related error details.
type SchemaError struct {
	Filepath string `json:"filepath,omitempty"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Details  string `json:"details,omitempty"`
}

// NewVetCmd returns the vet command.
func NewVetCmd() *cobra.Command {
	o := VetOptions{}
	cmd := &cobra.Command{
		Use:     "vet",
		Short:   "KCL validation tool",
		Long:    vetDesc,
		Example: vetExample,
		RunE: func(_ *cobra.Command, args []string) error {
			dataFile := args[0]
			codeFile := args[1]
			return doValidate(dataFile, codeFile, &o)
		},
		SilenceUsage: true,
	}

	// Two positional arguments <data_file> <kcl_file>
	cmd.Args = cobra.ExactArgs(2)
	cmd.Flags().StringVarP(&o.Schema, "schema", "s", "",
		"Specify the validate schema.")
	cmd.Flags().StringVarP(&o.Schema, "attribute_name", "a", "",
		"Specify the validate config attribute name.")
	cmd.Flags().StringVar(&o.Format, "format", "",
		"Specify the validate data format. e.g., yaml, json. Default is json")
	cmd.Flags().StringVar(&o.Output, "output", "text",
		"Specify the output format. e.g., text, json. Default is text")

	return cmd
}

func doValidate(dataFile, codeFile string, o *VetOptions) error {
	var ok bool
	var errMsg string
	if dataFile == "-" {
		// Read data from stdin
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return outputResult(o.Output, false, "", err)
		}
		code, err := os.ReadFile(codeFile)
		if err != nil {
			return outputResult(o.Output, false, "", err)
		}
		ok, err = validate.ValidateCode(string(input), string(code), &o.ValidateOptions)
		if err != nil {
			return outputResult(o.Output, false, err.Error(), nil)
		}
	} else {
		// Read data from files
		dataFiles, err := fs.ExpandInputFiles([]string{dataFile}, false)
		if err != nil {
			return outputResult(o.Output, false, "", err)
		}
		for _, dataFile := range dataFiles {
			ok, errMsg, err = validateFile(dataFile, codeFile, &o.ValidateOptions)
			if err != nil {
				return outputResult(o.Output, false, "", err)
			}
			if !ok {
				return outputResult(o.Output, false, errMsg, nil)
			}
		}
	}
	return outputResult(o.Output, ok, "", nil)
}

func validateFile(dataFile, codeFile string, opts *validate.ValidateOptions) (ok bool, errMsg string, err error) {
	if opts == nil {
		opts = &validate.ValidateOptions{}
	}
	svc := kcl.Service()
	resp, err := svc.ValidateCode(&gpyrpc.ValidateCodeArgs{
		Datafile:      dataFile,
		File:          codeFile,
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
	})
	if err != nil {
		return false, "", err
	}
	return resp.Success, resp.ErrMessage, nil
}

// outputResult outputs the validation result in the specified format.
func outputResult(outputFormat string, success bool, errMsg string, err error) error {
	if strings.ToLower(outputFormat) == "json" {
		return outputJSON(success, errMsg, err)
	}
	// Default text output
	return outputText(success, errMsg, err)
}

// outputText outputs the validation result in text format (original behavior).
func outputText(success bool, errMsg string, err error) error {
	if err != nil {
		return err
	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	if success {
		fmt.Println("Validate success!")
	}
	return nil
}

// outputJSON outputs the validation result in JSON format.
func outputJSON(success bool, errMsg string, err error) error {
	result := VetResult{
		Success: success,
	}

	if err != nil {
		result.Errors = []VetError{{
			ErrorType: "Error",
			Message:   stripansi.Strip(err.Error()),
		}}
		result.ErrCount = 1
	} else if errMsg != "" {
		// Strip ANSI codes from the error message before parsing
		cleanErrMsg := stripansi.Strip(errMsg)
		result.Errors = parseErrorMessage(cleanErrMsg)
		result.ErrCount = len(result.Errors)
	} else if success {
		result.Message = "Validate success!"
	}

	jsonOutput, jsonErr := json.MarshalIndent(result, "", "  ")
	if jsonErr != nil {
		return jsonErr
	}
	fmt.Println(string(jsonOutput))
	return nil
}

// parseErrorMessage attempts to parse the error message into structured errors.
// KCL error format example:
//
//	EvaluationError
//	 --> path/to/file.yaml:3:9
//	  |
//	3 | app_name: "test"
//	  |         ^ Instance check failed
func parseErrorMessage(errMsg string) []VetError {
	var vetErrors []VetError

	// Pattern to match error location: --> filepath:line:column
	locationPattern := regexp.MustCompile(`-->\s*([^:]+):(\d+):(\d+)`)
	// Pattern to match error type at the start
	// Pattern to match error type at the start
	errorTypePattern := regexp.MustCompile(`^(\w+Error|\w*Exception)`)
	// Pattern to match the error message after ^
	messagePattern := regexp.MustCompile(`\^\s*(.+)$`)
	// Pattern to match code snippet (line number | code)
	snippetPattern := regexp.MustCompile(`^\s*\d+\s*\|\s*(.+)$`)

	lines := strings.Split(errMsg, "\n")

	var currentError *VetError
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "|" {
			continue
		}

		// Check for error type
		if matches := errorTypePattern.FindStringSubmatch(line); matches != nil {
			if currentError != nil {
				vetErrors = append(vetErrors, *currentError)
			}
			currentError = &VetError{
				ErrorType: matches[1],
			}
			continue
		}

		// Check for location
		if matches := locationPattern.FindStringSubmatch(line); matches != nil {
			if currentError == nil {
				currentError = &VetError{}
			}
			currentError.File = matches[1]
			if lineNum, err := strconv.Atoi(matches[2]); err == nil {
				currentError.Line = lineNum
			}
			if colNum, err := strconv.Atoi(matches[3]); err == nil {
				currentError.Column = colNum
			}
			continue
		}

		// Check for error message (contains ^)
		if matches := messagePattern.FindStringSubmatch(line); matches != nil {
			if currentError != nil {
				currentError.Message = strings.TrimSpace(matches[1])
			}
			continue
		}

		// Check for code snippet
		if matches := snippetPattern.FindStringSubmatch(line); matches != nil {
			if currentError != nil && currentError.CodeSnippet == "" {
				currentError.CodeSnippet = strings.TrimSpace(matches[1])
			}
			continue
		}
	}

	// Don't forget the last error
	if currentError != nil {
		vetErrors = append(vetErrors, *currentError)
	}

	// If parsing failed, return the raw message as a single error
	if len(vetErrors) == 0 {
		vetErrors = append(vetErrors, VetError{
			Message: errMsg,
		})
	}

	return vetErrors
}
