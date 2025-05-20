package iso8583parser

import (
	"os"

	"gopkg.in/yaml.v2"
)

// FieldSpec contains fields that describes an iso8583 Field
type FieldSpec struct {
	ContentType string `yaml:"ContentType"`
	MaxLen      int    `yaml:"MaxLen"`
	MinLen      int    `yaml:"MinLen"`
	LenType     string `yaml:"LenType"`
	Label       string `yaml:"Label"`
}

// Spec contains the fields that describes an iso8583 specification
type SpecData struct {
	Fields map[int]FieldSpec
}

// Read specification from the spesific yaml configuration file
func (s *SpecData) readFromFile(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return err
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, &s.Fields)
}

// Check field excluding Field 0 (MTI) and Field 1 (bitmap auto-generated)
func (s *SpecData) hasAtLeastOneDataField() bool {
	for field := range s.Fields {
		if field != 0 && field != 1 {
			return true
		}
	}
	return false
}

// Create new SpecData object from the file specification
// Errors can occur if have an error from file like file not found, failed to read file, etc
// and when the file does not match the specified specifications
func SpecFromFile(filename string) (spec SpecData, err error) {
	if err := spec.readFromFile(filename); err != nil {
		return spec, err
	}
	return spec, nil
}
