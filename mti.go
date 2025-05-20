package iso8583parser

import (
	"strconv"
)

type MtiData struct {
	mti string
}

// Retrieving MTI code
func (m *MtiData) Get() string {
	return m.mti
}

// Private function to validate MTI code
// Errors can occur if invalid length of MTI or MTI is not an integer value
func (m *MtiData) validate() error {
	if len(m.mti) != 4 {
		return ErrInvalidMtiLength
	}

	_, err := strconv.ParseInt(m.mti, 10, 64)
	if err != nil {
		return ErrInvalidMtiInteger
	}

	return nil
}
