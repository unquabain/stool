package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func Parse(filename string, format Format) (map[string]any, error) {
	var (
		err    error
		reader io.Reader
	)
	if filename == `-` {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf(`could not open %q: %w`, filename, err)
		}
		reader = file
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf(`could not read %q: %w`, filename, err)
	}
	if format == FormatUnknown {
		format, err = Detect(data)
		if err != nil {
			return nil, fmt.Errorf(`could not determine format of %q: %w`, filename, err)
		}
	}
	if format == FormatUnknown {
		return nil, fmt.Errorf(`could not determine format of %q`, filename)
	}
	parsed := make(map[string]any)
	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf(`could not unmarshal JSON data from %q: %w`, filename, err)
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf(`could not unmarshal YAML data from %q: %w`, filename, err)
		}
	}
	return parsed, nil
}
