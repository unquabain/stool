package main

import (
	"fmt"
	"regexp"
)

// Detector determines if a byte array is a certain structured file type.
type Detector interface {
	Detect([]byte) (bool, error)
}

// DetectorFunc is an adapter for a stateless Detector.
type DetectorFunc func([]byte) (bool, error)

// Detect implements the Detector interface on DetectorFunc adapters.
func (df DetectorFunc) Detect(data []byte) (bool, error) {
	return df(data)
}

// Detectors maps different detectors to format types.
type Detectors map[Format]Detector

// Detect tries to determine the structured data type of a byte slice
func (ds Detectors) Detect(data []byte) (Format, error) {
	for format, dtor := range ds {
		match, err := dtor.Detect(data)
		if err != nil {
			return FormatUnknown, fmt.Errorf(`could not determine structured file type: %w`, err)
		}
		if match {
			return format, nil
		}
	}
	return FormatUnknown, nil
}

type RegexDetector regexp.Regexp

func NewRegexDetector(pattern string) *RegexDetector {
	return (*RegexDetector)(regexp.MustCompile(pattern))
}
func (rd *RegexDetector) Regexp() *regexp.Regexp {
	return (*regexp.Regexp)(rd)
}
func (rd *RegexDetector) Detect(data []byte) (bool, error) {
	return rd.Regexp().Match(data), nil
}

var _detectors Detectors

func init() {
	_detectors = make(Detectors)
	_detectors[FormatJSON] = NewRegexDetector(`"\w+":`)
	_detectors[FormatYAML] = NewRegexDetector(`\s*\w+:`)
}

func Detect(data []byte) (Format, error) {
	return _detectors.Detect(data)
}
