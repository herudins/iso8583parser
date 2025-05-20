package iso8583parser

import (
	"fmt"
	"strings"
)

// Extract the MTI code from ISO8583 message
func extractMti(str string) (mti MtiData, err error) {
	if len(str) < 4 {
		return MtiData{}, ErrInvalidMtiLength
	}

	mtiData := str[0:4]
	return MtiData{mti: mtiData}, nil
}

// Get variable length form field type in field spesification
// the field type is: llvar, lllvar and llllvar
func getVariableLengthFromString(str string) (num int, err error) {
	str = strings.ToLower(str)
	if str == "llvar" {
		return 2, nil
	}
	if str == "lllvar" {
		return 3, nil
	}
	if str == "llllvar" {
		return 4, nil
	}

	return num, fmt.Errorf("%s is an invalid LenType", str)
}

// Create text with prefix padding if the text length is less than the maximum length
// If the data exceeds the maximum length, the original data will be returned.
func leftPad(s string, l int, pad string) string {
	if len(s) >= l {
		return s
	}
	padding := strings.Repeat(pad, l-len(s))

	return padding + s
}

// Create text with suffix padding if the text length is less than the maximum length
// If the data exceeds the maximum length, the original data will be returned.
func rightPad(s string, l int, pad string) string {
	if len(s) >= l {
		return s
	}
	padding := strings.Repeat(pad, l-len(s))

	return s + padding
}

// HexToBits converts a hexadecimal string to an array of bits (0 and 1) in the form []byte.
func HexToBits(hexStr string) ([]byte, error) {
	n := len(hexStr)
	bits := make([]byte, 0, n*4)

	for i := range hexStr {
		var value byte
		ch := hexStr[i]
		switch {
		case '0' <= ch && ch <= '9':
			value = ch - '0'
		case 'a' <= ch && ch <= 'f':
			value = ch - 'a' + 10
		case 'A' <= ch && ch <= 'F':
			value = ch - 'A' + 10
		default:
			return nil, fmt.Errorf("invalid hex character: %c", ch)
		}

		// Get 4 bit (MSB to LSB)
		bits = append(bits,
			(value>>3)&1,
			(value>>2)&1,
			(value>>1)&1,
			value&1,
		)
	}

	return bits, nil
}

// HexToBitsString converts a hexadecimal string to a binary bit string (e.g. "1fa3" -> "0001111110100011")
func HexToBitsString(hexStr string) (string, error) {
	var builder strings.Builder
	builder.Grow(len(hexStr) * 4) // Pre-allocate memory for efficiency

	for i := range hexStr {
		var value byte
		ch := hexStr[i]
		switch {
		case '0' <= ch && ch <= '9':
			value = ch - '0'
		case 'a' <= ch && ch <= 'f':
			value = ch - 'a' + 10
		case 'A' <= ch && ch <= 'F':
			value = ch - 'A' + 10
		default:
			return "", fmt.Errorf("invalid hex character: %c", ch)
		}

		// Append 4 bits as characters ('0' or '1')
		for j := 3; j >= 0; j-- {
			if (value>>j)&1 == 1 {
				builder.WriteByte('1')
			} else {
				builder.WriteByte('0')
			}
		}
	}

	return builder.String(), nil
}

// BitsToHex converts an array of bits (0 and 1) to a hexadecimal string.
func BitsToHex(bits []byte) (string, error) {
	if len(bits)%4 != 0 {
		return "", ErrInvalidBitLength
	}

	hexBytes := make([]byte, len(bits)/4)

	for i := 0; i < len(bits); i += 4 {
		nibble := (bits[i] << 3) | (bits[i+1] << 2) | (bits[i+2] << 1) | bits[i+3]
		if nibble < 10 {
			hexBytes[i/4] = '0' + nibble
		} else {
			hexBytes[i/4] = 'a' + (nibble - 10)
		}
	}

	return string(hexBytes), nil
}

// BitsStringToHex converts a binary bit string (e.g. "11001100") to a hexadecimal string.
func BitsStringToHex(bitStr string) (string, error) {
	n := len(bitStr)
	if n%4 != 0 {
		return "", ErrInvalidBitLength
	}

	hex := make([]byte, n/4)

	for i := 0; i < n; i += 4 {
		var nibble byte
		for j := 0; j < 4; j++ {
			ch := bitStr[i+j]
			if ch != '0' && ch != '1' {
				return "", fmt.Errorf("invalid bit character: %c", ch)
			}

			nibble = (nibble << 1) | (ch - '0')
		}

		if nibble < 10 {
			hex[i/4] = '0' + nibble
		} else {
			hex[i/4] = 'a' + (nibble - 10)
		}
	}

	return string(hex), nil
}

// BitsIntArrayToHex converts an array of int (0 and 1) to a hexadecimal string.
func BitsIntArrayToHex(bits []int) (string, error) {
	length := len(bits)
	if length%4 != 0 || ((length/4)%2) != 0 {
		return "", ErrInvalidBitLength
	}

	hex := make([]byte, length/4)

	for i := 0; i < length; i += 4 {
		var nibble byte
		for j := 0; j < 4; j++ {
			b := bits[i+j]
			if b != 0 && b != 1 {
				return "", fmt.Errorf("invalid bit value at index %d: %d", i+j, b)
			}
			nibble = (nibble << 1) | byte(b)
		}

		if nibble < 10 {
			hex[i/4] = '0' + nibble
		} else {
			hex[i/4] = 'a' + (nibble - 10)
		}
	}

	return string(hex), nil
}
