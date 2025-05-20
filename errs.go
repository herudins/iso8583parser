package iso8583parser

import (
	"errors"
)

var (
	ErrInvalidMtiLength           = errors.New("MTI must be length (4)")
	ErrInvalidMtiInteger          = errors.New("MTI can only contain integers")
	ErrEmptySpec                  = errors.New("specification is empty")
	ErrSpecMinHasOneField         = errors.New("specification minimum has one field or more without field 0 and 1")
	ErrInvalidBitLength           = errors.New("bit length must be multiple of 4")
	ErrDataToShortSecondaryBitmap = errors.New("data too short for secondary bitmap")
	ErrDataToShortTertiaryBitmap  = errors.New("data too short for tertiary bitmap")
	ErrIsoMessageTooShort         = errors.New("data iso message too short")
	ErrEmptyDataElements          = errors.New("elements data empty")
)
