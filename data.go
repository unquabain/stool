package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Unmarshaller func([]byte, any) error

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
	var unmarshallers []Unmarshaller
	switch format {
	case FormatJSON:
		unmarshallers = []Unmarshaller{json.Unmarshal, yaml.Unmarshal}
	case FormatYAML:
		unmarshallers = []Unmarshaller{yaml.Unmarshal, json.Unmarshal}
	default:
		return nil, fmt.Errorf(`could not determine format of %q`, filename)
	}
	parsed := make(map[string]any)
	for _, um := range unmarshallers {
		err = um(data, &parsed)
		if err == nil {
			return parsed, nil
		}
		log.Printf(`WARN: unable to parse %q as %s: %s`, filename, format, err.Error())
	}
	return nil, fmt.Errorf(`no unmarshallers were able to parse %q`, filename)
}
