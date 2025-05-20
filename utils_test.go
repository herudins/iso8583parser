package iso8583parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	hexStr   = "bf38404109e300000000000013000000"
	bitStr   = "10111111001110000100000001000001000010011110001100000000000000000000000000000000000000000000000000010011000000000000000000000000"
	bitBytes = []byte{1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

func TestExtractMti(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		mtiData, err := extractMti("2200")
		assert.Nil(t, err)
		assert.Equal(t, MtiData{mti: "2200"}, mtiData, "Expected MTI data to be equal")
	})

	t.Run("Invalid Length", func(t *testing.T) {
		_, err := extractMti("10")
		assert.NotNil(t, err, "Expected error mti length")
	})
}

func TestGetVariableLengthFromString(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		num, err := getVariableLengthFromString("llvar")
		assert.Nil(t, err, "Expected nil value error")
		assert.Equal(t, 2, num, "Expected num to be equal")
	})

	t.Run("Invalid code length", func(t *testing.T) {
		_, err := getVariableLengthFromString("dsdsfe")
		assert.NotNil(t, err, "Expected error length ")
	})
}

func TestHexToBits(t *testing.T) {
	bytes, err := HexToBits(hexStr)
	if err != nil {
		t.Errorf("HexToBits failed: %s", err.Error())
		t.FailNow()
	}

	assert.Equal(t, bitBytes, bytes, "Expected byte array to be equal")
}

func TestHexToBitsString(t *testing.T) {
	bitsData, err := HexToBitsString(hexStr)
	if err != nil {
		t.Errorf("HexToBitsString failed: %s", err.Error())
		t.FailNow()
	}

	assert.Equal(t, bitStr, bitsData, "Expected bit string to be equal")
}

func BenchmarkHexToBits(b *testing.B) {
	b.Run("HexToBits", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := HexToBits(hexStr)
			if err != nil {
				b.Errorf("HexToBits failed: %s", err.Error())
				b.FailNow()
			}
		}
	})

	b.Run("HexToBitsString", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := HexToBitsString(hexStr)
			if err != nil {
				b.Errorf("HexToBitsString failed: %s", err.Error())
				b.FailNow()
			}
		}
	})
}

func TestBitsToHex(t *testing.T) {
	hexData, err := BitsToHex(bitBytes)
	if err != nil {
		t.Errorf("BitsToHex failed: %s", err.Error())
		t.FailNow()
	}

	assert.Equal(t, hexStr, hexData, "Expected hex string to be equal")
}

func TestBitsStringToHex(t *testing.T) {
	hexData, err := BitsStringToHex(bitStr)
	if err != nil {
		t.Errorf("BitsStringToHex failed: %s", err.Error())
		t.FailNow()
	}

	assert.Equal(t, hexStr, hexData, "Expected hex string to be equal")
}

func BenchmarkBitsToHex(b *testing.B) {
	b.Run("BitsToHex", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := BitsToHex(bitBytes)
			if err != nil {
				b.Errorf("BitsToHex failed: %s", err.Error())
				b.FailNow()
			}
		}
	})

	b.Run("BitsStringToHex", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := BitsStringToHex(bitStr)
			if err != nil {
				b.Errorf("BitsStringToHex failed: %s", err.Error())
				b.FailNow()
			}
		}
	})
}

func TestGetShortedKeyFields(t *testing.T) {
	source := map[int]string{
		3:   "100700",
		4:   "1500",
		5:   "5",
		7:   "0711170215",
		8:   "8",
		47:  "147",
		48:  "12345",
		100: "123456",
		103: "1234567890",
		104: "654321",
		6:   "6",
		11:  "23edfr",
		12:  "202307",
		13:  "0711",
	}

	keyShorted := GetShortedKeyFields(source)
	fmt.Println(keyShorted)
}
