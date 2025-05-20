package iso8583parser

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ElementsData object
type ElementsData struct {
	mu       sync.RWMutex
	elements map[int]string
}

// setElement addd element data based on a specific field
func (e *ElementsData) setElement(field int, data string) {
	e.mu.Lock()
	e.elements[field] = data
	e.mu.Unlock()
}

// getElement retrieving element data based on a specific field
// from the elements map
func (e *ElementsData) getElement(field int) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.elements[field]
}

// getElements retrieve all element data from the elements map
func (e *ElementsData) getElements() map[int]string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.elements
}

// Iso8583Data object
type Iso8583Data struct {
	Spec       SpecData
	Mti        MtiData
	Bitmap     []int
	BitmapSize int
	Elements   ElementsData
}

// Create a new Iso8583Data object from a yaml specification file
// and can set or use option secondary bitmap to true or false
func New(filename string, secondaryBitmap bool) (iso *Iso8583Data, err error) {
	spec, err := SpecFromFile(filename)
	if err != nil {
		return iso, err
	}

	return createIsoObject(spec, secondaryBitmap)
}

// Create a new Iso8583Data object from a predefined data specification
// and can set or use option secondary bitmap to true or false
func NewFromSpec(spec SpecData, secondaryBitmap bool) (iso *Iso8583Data, err error) {
	return createIsoObject(spec, secondaryBitmap)
}

// Private function that create a new Iso8583Data object from a predefined data specification.
func createIsoObject(spec SpecData, secondaryBitmap bool) (iso *Iso8583Data, err error) {
	if len(spec.Fields) == 0 {
		return iso, ErrEmptySpec
	}

	if !spec.hasAtLeastOneDataField() {
		return iso, ErrSpecMinHasOneField
	}

	iso = &Iso8583Data{
		Spec:       spec,
		Mti:        MtiData{},
		Bitmap:     make([]int, 64),
		BitmapSize: 64,
		Elements:   ElementsData{elements: make(map[int]string)},
	}

	if secondaryBitmap {
		iso.Bitmap = make([]int, 128)
		iso.BitmapSize = 128
		iso.Bitmap[0] = 1
	}

	return iso, nil
}

// Set MTI data for the iso8583 message
func (iso *Iso8583Data) AddMTI(mti string) error {
	mtiData := MtiData{mti: mti}
	if err := mtiData.validate(); err != nil {
		return err
	}
	iso.Mti = mtiData
	return nil
}

// Define specific field data by field number.
// An error may occur if the field number entered is less than 2 or exceeds the capacity of the bitmap.
func (iso *Iso8583Data) SetField(field int, data string) error {
	if field < 2 || field > len(iso.Bitmap) {
		return fmt.Errorf("expected field to be between %d and %d found %d instead", 2, len(iso.Bitmap), field)
	}

	iso.Bitmap[field-1] = 1
	iso.Elements.setElement(field, data)
	return nil
}

// Retrieves specific field data by field number.
// An error may occur if the field number entered is less than 2 or exceeds the capacity of the bitmap.
func (iso *Iso8583Data) GetField(field int) (string, error) {
	if field < 2 || field > len(iso.Bitmap) {
		return "", fmt.Errorf("expected field to be between %d and %d found %d instead", 2, len(iso.Bitmap), field)
	}

	return iso.Elements.getElement(field), nil
}

// Retrieves all field data in all elements.
// An error may occur if the elements is zero len
func (iso *Iso8583Data) GetAllFields() (allField map[int]string, err error) {
	allField = iso.Elements.getElements()
	if len(allField) <= 0 {
		return nil, ErrEmptyDataElements
	}

	return allField, nil
}

// Perform ISO8583 data packaging based on fields and data that have been set returning bytes data iso message
// Errors can occur if the data length in a particular field exceeds the field capacity in the configuration,
// and if the field type depends on the length of the variable in the configuration
// and is not part of the type (llvar, lllvar, llllvar),
// and if bitmap is invalid
func (iso *Iso8583Data) Marshal() ([]byte, error) {
	return iso.marshal()
}

// Perform ISO8583 data packaging based on fields and data that have been set returning text of iso message
// Errors can occur if the data length in a particular field exceeds the field capacity in the configuration,
// and if the field type depends on the length of the variable in the configuration
// and is not part of the type (llvar, lllvar, llllvar),
// and if bitmap is invalid
func (iso *Iso8583Data) MarshalString() (string, error) {
	bytesData, err := iso.marshal()
	if err != nil {
		return "", err
	}

	return string(bytesData), nil
}

func (iso *Iso8583Data) marshal() ([]byte, error) {
	buf := make([]byte, 0, 512)
	buf = append(buf, []byte(iso.Mti.Get())...)

	bitmap := iso.Bitmap
	elementsSpec := iso.Spec

	bitmapString, err := BitsIntArrayToHex(iso.Bitmap)
	if err != nil {
		return nil, err
	}
	buf = append(buf, []byte(bitmapString)...)

	//Loop from Field 2 (bitmap index 1),
	// because Field 0 (MTI) and Field 1 (bitmap index 0) is bitmap auto-generated
	for index := 1; index < iso.BitmapSize; index++ {
		if bitmap[index] == 1 {
			var (
				field     = index + 1
				fieldSpec = elementsSpec.Fields[field]
				data      = iso.Elements.getElement(field)
				maxLen    = fieldSpec.MaxLen
				dataLen   = len(data)
			)

			if dataLen > maxLen {
				return nil, fmt.Errorf("failed to marshal field %d with max length %d but data length %d", field, maxLen, dataLen)
			}

			if strings.ToLower(fieldSpec.LenType) == "fixed" {
				if fieldSpec.ContentType == "n" {
					data = leftPad(data, maxLen, "0")
				} else {
					data = rightPad(data, maxLen, " ")
				}

				buf = append(buf, data...)
			} else {
				lengthType, err := getVariableLengthFromString(fieldSpec.LenType)
				if err != nil {
					return nil, err
				}

				paddedLength := leftPad(strconv.Itoa(dataLen), lengthType, "0")
				buf = append(buf, paddedLength...)
				buf = append(buf, data...)
			}
		}
	}

	return buf, nil
}

// Perform ISO8583 data parsing according to predetermined specifications
// form the data sent is in the form of a byte array
func (iso *Iso8583Data) Unmarshal(bytesIso []byte) error {
	if len(bytesIso) < 20 {
		return ErrIsoMessageTooShort
	}

	specs := iso.Spec
	mti := string(bytesIso[:4])

	mtiData, _ := extractMti(mti)
	if err := mtiData.validate(); err != nil {
		return err
	}

	iso.Mti = mtiData

	bitmapHex := bytesIso[4:20]
	bitmap := make([]byte, 8)

	if _, err := hex.Decode(bitmap, bitmapHex); err != nil {
		return err
	}

	bitmapSize := 64
	if bitmap[0]&0x80 != 0 {
		bitmapSize = 128
		bitmap = append(bitmap, make([]byte, 8)...)
		if len(bytesIso) < 36 {
			return ErrDataToShortSecondaryBitmap
		}

		if _, err := hex.Decode(bitmap[8:], bytesIso[20:36]); err != nil {
			return err
		}

		iso.Bitmap[0] = 1
		bytesIso = bytesIso[36:]
	} else {
		bytesIso = bytesIso[20:]
	}

	iso.BitmapSize = bitmapSize

	pos := 0
	for i := 2; i <= bitmapSize; i++ {
		bytePos := (i - 1) / 8
		bitPos := uint(7 - ((i - 1) % 8))

		if bitmap[bytePos]&(1<<bitPos) != 0 {
			spec, ok := specs.Fields[i]
			if !ok {
				return fmt.Errorf("no field spec for field %d", i)
			}

			var fieldLen int
			switch spec.LenType {
			case "fixed":
				fieldLen = spec.MaxLen
			case "llvar":
				if pos+2 > len(bytesIso) {
					return fmt.Errorf("field %d: LLVAR prefix too short", i)
				}

				n, err := strconv.Atoi(string(bytesIso[pos : pos+2]))
				if err != nil {
					return fmt.Errorf("field %d: LLVAR prefix is not an integer", i)
				}

				fieldLen = n
				pos += 2
			case "lllvar":
				if pos+3 > len(bytesIso) {
					return fmt.Errorf("field %d: LLLVAR prefix too short", i)
				}

				n, err := strconv.Atoi(string(bytesIso[pos : pos+3]))
				if err != nil {
					return fmt.Errorf("field %d: LLLVAR prefix is not an integer", i)
				}

				fieldLen = n
				pos += 3
			case "llllvar":
				if pos+3 > len(bytesIso) {
					return fmt.Errorf("field %d: LLLLVAR prefix too short", i)
				}

				n, err := strconv.Atoi(string(bytesIso[pos : pos+4]))
				if err != nil {
					return fmt.Errorf("field %d: LLLLVAR prefix is not an integer", i)
				}

				fieldLen = n
				pos += 4
			}

			if pos+fieldLen > len(bytesIso) {
				return fmt.Errorf("field %d: value too short", i)
			}

			iso.Bitmap[i-1] = 1
			iso.SetField(i, string(bytesIso[pos:pos+fieldLen]))
			pos += fieldLen
		}
	}

	return nil
}

// Perform ISO8583 data parsing according to predetermined specifications
// form the data sent is in the form of string
func (iso *Iso8583Data) UnmarshalString(isoMessage string) error {
	if len(isoMessage) < 20 {
		return ErrIsoMessageTooShort
	}

	specs := iso.Spec
	mti := string(isoMessage[:4])

	mtiData, _ := extractMti(mti)
	if err := mtiData.validate(); err != nil {
		return err
	}

	iso.Mti = mtiData

	bitmapHex := isoMessage[4:20]
	bitmap, err := HexToBitsString(bitmapHex)
	if err != nil {
		return err
	}

	if bitmap[0] == '1' {
		if len(isoMessage) < 36 {
			return ErrDataToShortSecondaryBitmap
		}

		second := isoMessage[20:36]
		secondBitmap, err := HexToBitsString(second)
		if err != nil {
			return err
		}
		bitmap += secondBitmap

		iso.Bitmap[0] = 1
		isoMessage = isoMessage[36:]
	} else {
		isoMessage = isoMessage[20:]
	}

	pos := 0
	for i, c := range bitmap {
		bit := i + 1
		if c != '1' || i == 0 || i == 1 {
			continue
		}

		spec, ok := specs.Fields[bit]
		if !ok {
			return fmt.Errorf("no field spec for field %d", i)
		}

		var fieldLen int
		switch spec.LenType {
		case "fixed":
			fieldLen = spec.MaxLen
		case "llvar":
			if pos+2 > len(isoMessage) {
				return fmt.Errorf("field %d: LLVAR prefix too short", i)
			}

			n, err := strconv.Atoi(string(isoMessage[pos : pos+2]))
			if err != nil {
				return fmt.Errorf("field %d: LLVAR prefix is not an integer", i)
			}

			fieldLen = n
			pos += 2
		case "lllvar":
			if pos+3 > len(isoMessage) {
				return fmt.Errorf("field %d: LLLVAR prefix too short", i)
			}

			n, err := strconv.Atoi(string(isoMessage[pos : pos+3]))
			if err != nil {
				return fmt.Errorf("field %d: LLLVAR prefix is not an integer", i)
			}

			fieldLen = n
			pos += 3
		case "llllvar":
			if pos+3 > len(isoMessage) {
				return fmt.Errorf("field %d: LLLLVAR prefix too short", i)
			}

			n, err := strconv.Atoi(string(isoMessage[pos : pos+4]))
			if err != nil {
				return fmt.Errorf("field %d: LLLLVAR prefix is not an integer", i)
			}

			fieldLen = n
			pos += 4
		}

		if pos+fieldLen > len(isoMessage) {
			return fmt.Errorf("field %d: value too short", i)
		}

		iso.Bitmap[bit-1] = 1
		iso.SetField(bit, string(isoMessage[pos:pos+fieldLen]))
		pos += fieldLen
	}
	iso.BitmapSize = len(bitmap)

	return nil
}
