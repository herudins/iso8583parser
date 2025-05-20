package iso8583parser

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	bitArray         = []int{1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bitArrayTertiary = []int{1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	msgiso           = "2200bf38404109e30000000000001300000010070000000000150000000000000500000000000607111702150000000823edfr20230707110014260232dhfte4736fge40 41      42             43                                      0031470051234506123456101234567890006654321"
	msgisoTertiary   = "2200b8000000000000008000000013000000c000000000000000100700000000001500000000000005 061234561012345678900066543210000000500000005"
)

func setDataIso(isoParser *Iso8583Data) {
	isoParser.AddMTI("2200")
	isoParser.SetField(3, "100700")
	isoParser.SetField(4, "1500")
	isoParser.SetField(5, "5")
	isoParser.SetField(6, "6")
	isoParser.SetField(7, "0711170215")
	isoParser.SetField(8, "8")
	isoParser.SetField(11, "23edfr")
	isoParser.SetField(12, "202307")
	isoParser.SetField(13, "0711")
	isoParser.SetField(18, "0014")
	isoParser.SetField(26, "26")
	isoParser.SetField(32, "32")
	isoParser.SetField(37, "dhfte4736fge")
	isoParser.SetField(40, "40")
	isoParser.SetField(41, "41")
	isoParser.SetField(42, "42")
	isoParser.SetField(43, "43")
	isoParser.SetField(47, "147")
	isoParser.SetField(48, "12345")
	isoParser.SetField(100, "123456")
	isoParser.SetField(103, "1234567890")
	isoParser.SetField(104, "654321")
}

func TestMarshal(t *testing.T) {
	isoParser, err := New("spec1987.yml")
	assert.Nil(t, err, "Error should be nil")

	setDataIso(isoParser)
	isoMsg, err := isoParser.Marshal()
	assert.Nil(t, err, "Error should be nil")

	require.Equal(t, bitArray, isoParser.Bitmap, "Expected bit string to be equal")
	require.Equal(t, msgiso, string(isoMsg), "Expected iso message to be equal")
}

func TestMarshalString(t *testing.T) {
	isoParser, err := New("spec1987.yml")
	assert.Nil(t, err, "Error should be nil")

	setDataIso(isoParser)
	isoMsg, err := isoParser.Marshal()
	assert.Nil(t, err, "Error should be nil")

	require.Equal(t, bitArray, isoParser.Bitmap, "Expected bit string to be equal")
	require.Equal(t, msgiso, string(isoMsg), "Expected iso message to be equal")
}

func TestMarshalRace(t *testing.T) {
	isoParser, err := New("spec1987.yml")
	assert.Nil(t, err, "Error should be nil")

	source := map[int]string{
		3:   "100700",
		4:   "1500",
		5:   "5",
		6:   "6",
		7:   "0711170215",
		8:   "8",
		11:  "23edfr",
		12:  "202307",
		13:  "0711",
		47:  "147",
		48:  "12345",
		100: "123456",
		103: "1234567890",
		104: "654321",
	}

	var wg sync.WaitGroup
	for field, data := range source {
		wg.Add(1)
		go func(g *sync.WaitGroup) {
			defer wg.Done()

			isoParser.SetField(field, data)
		}(&wg)
	}
	wg.Wait()

	_, err = isoParser.Marshal()
	assert.Nil(t, err, "Error should be nil")
}

func setDataIsoTertiary(isoParser *Iso8583Data) {
	isoParser.AddMTI("2200")
	isoParser.SetField(3, "100700")
	isoParser.SetField(4, "1500")
	isoParser.SetField(5, "5")
	isoParser.SetField(100, "123456")
	isoParser.SetField(103, "1234567890")
	isoParser.SetField(104, "654321")
	isoParser.SetField(129, "5")
	isoParser.SetField(130, "5")
}

func TestMarshalTertiary(t *testing.T) {
	isoParser, err := New("./spec1987_tertiary.yml")
	assert.Nil(t, err, "Error should be nil")

	setDataIsoTertiary(isoParser)
	isoMsg, err := isoParser.Marshal()
	assert.Nil(t, err, "Error should be nil")

	require.Equal(t, bitArrayTertiary, isoParser.Bitmap, "Expected bit string to be equal")
	require.Equal(t, msgisoTertiary, string(isoMsg), "Expected iso message to be equal")
}

func BenchmarkMarshal(b *testing.B) {
	isoParser, err := New("spec1987.yml")
	if err != nil {
		b.Fatalf("Failed to initialize parser: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setDataIso(isoParser)
		_, err := isoParser.Marshal()
		if err != nil {
			b.Fatalf("Failed to Marshal: %v", err)
		}
	}
}

func BenchmarkMarshalString(b *testing.B) {
	isoParser, err := New("spec1987.yml")
	if err != nil {
		b.Fatalf("Failed to initialize parser: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setDataIso(isoParser)
		_, err := isoParser.MarshalString()
		if err != nil {
			b.Fatalf("Failed to Marshal: %v", err)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	isoParser, err := New("spec1987.yml")
	assert.Nil(t, err, "Error should be nil")

	err = isoParser.Unmarshal([]byte(msgiso))
	assert.Nil(t, err, "Error should be nil")

	bit3, err := isoParser.GetField(3)
	assert.Nil(t, err, "Error should be nil")

	bit4, err := isoParser.GetField(4)
	assert.Nil(t, err, "Error should be nil")

	bit5, err := isoParser.GetField(5)
	assert.Nil(t, err, "Error should be nil")

	require.Equal(t, bitArray, isoParser.Bitmap, "Expected bit string to be equal")
	require.Equal(t, "2200", isoParser.Mti.Get(), "Expected MTI to be equal")
	require.Equal(t, "100700", bit3, "Expected Bit3 to be equal")
	require.Equal(t, "000000001500", bit4, "Expected Bit4 to be equal")
	require.Equal(t, "000000000005", bit5, "Expected Bit5 to be equal")
}

func TestUnmarshalTertiary(t *testing.T) {
	isoParser, err := New("./spec1987_tertiary.yml")
	assert.Nil(t, err, "Error should be nil")

	err = isoParser.Unmarshal([]byte(msgisoTertiary))
	assert.Nil(t, err, "Error should be nil")

	bit3, err := isoParser.GetField(3)
	assert.Nil(t, err, "Error should be nil")

	bit4, err := isoParser.GetField(4)
	assert.Nil(t, err, "Error should be nil")

	bit100, err := isoParser.GetField(100)
	assert.Nil(t, err, "Error should be nil")

	bit129, err := isoParser.GetField(129)
	assert.Nil(t, err, "Error should be nil")

	require.Equal(t, bitArrayTertiary, isoParser.Bitmap, "Expected bit string to be equal")
	require.Equal(t, "2200", isoParser.Mti.Get(), "Expected MTI to be equal")
	require.Equal(t, "100700", bit3, "Expected Bit3 to be equal")
	require.Equal(t, "000000001500", bit4, "Expected Bit4 to be equal")
	require.Equal(t, "123456", bit100, "Expected Bit100 to be equal")
	require.Equal(t, "00000005", bit129, "Expected Bit129 to be equal")
}

func TestUnmarshalString(t *testing.T) {
	isoParser, err := New("spec1987.yml")
	assert.Nil(t, err, "Error should be nil")

	err = isoParser.UnmarshalString(msgiso)
	assert.Nil(t, err, "Error should be nil")

	bit3, err := isoParser.GetField(3)
	assert.Nil(t, err, "Error should be nil")

	bit4, err := isoParser.GetField(4)
	assert.Nil(t, err, "Error should be nil")

	bit5, err := isoParser.GetField(5)
	assert.Nil(t, err, "Error should be nil")

	require.Equal(t, bitArray, isoParser.Bitmap, "Expected bit string to be equal")
	require.Equal(t, "2200", isoParser.Mti.Get(), "Expected MTI to be equal")
	require.Equal(t, "100700", bit3, "Expected Bit3 to be equal")
	require.Equal(t, "000000001500", bit4, "Expected Bit4 to be equal")
	require.Equal(t, "000000000005", bit5, "Expected Bit5 to be equal")
}

func BenchmarkUnmarshal(b *testing.B) {
	isoParser, err := New("spec1987.yml")
	if err != nil {
		b.Fatalf("Failed to initialize parser: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := isoParser.Unmarshal([]byte(msgiso)); err != nil {
			b.Fatalf("Failed to Unmarshal: %v", err)
		}
	}
}

func BenchmarkUnmarshalString(b *testing.B) {
	isoParser, err := New("spec1987.yml")
	if err != nil {
		b.Fatalf("Failed to initialize parser: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := isoParser.UnmarshalString(msgiso); err != nil {
			b.Fatalf("Failed to Unmarshal: %v", err)
		}
	}
}
