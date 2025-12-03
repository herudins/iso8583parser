
# ISO8583 Parser

Simple golang library for pack and unpack iso8583 message.



## Installation

Install my-project with go get

```go
  go get https://github.com/herudins/iso8583parser
```
    
## Quick Start
The following example demonstrates how to:

- Pack and Unpack a message using structure from yaml file
- Pack and Unpack a message using structure using Go types

### Pack and Unpack using structure from yaml file
```go
// Load iso parser
parser, err := iso8583parser.New("spec1987.yml")
if err != nil {
    panic(err)
}

//Set data for ISO
parser.AddMTI("2200")
parser.SetField(3, "100700")
parser.SetField(4, "1500")
parser.SetField(5, "5")
parser.SetField(7, "0711170215")

// Pack the message
isoMsg, err := parser.MarshalString()
if err != nil {
    panic(err)
}

// Send packed message to the server
// ...

//Unpack message
if err := parser.Unmarshal([]byte(isoMsg)); err != nil {
    panic(err)
}
```

### Pack and Unpack using structure from GO Types
```go
//Define spec from Go Types
specData := iso8583parser.SpecData{
	Fields: map[int]iso8583parser.FieldSpec{
		3: iso8583parser.FieldSpec{
			ContentType: "n",
			MaxLen:      6,
			LenType:     "fixed",
			Label:       "Processing code",
		},
		4: iso8583parser.FieldSpec{
			ContentType: "n",
			MaxLen:      12,
			LenType:     "fixed",
			Label:       "Amount, transaction",
		},
		5: iso8583parser.FieldSpec{
			ContentType: "n",
			MaxLen:      12,
			LenType:     "fixed",
			Label:       "Amount, settlement",
		},
		7: iso8583parser.FieldSpec{
			ContentType: "n",
			MaxLen:      10,
			LenType:     "fixed",
			Label:       "Transmission date & time",
		},
	},
}

// Load iso parser
parser, err := iso8583parser.NewFromSpec(specData)
if err != nil {
    panic(err)
}

//Set data for ISO
parser.AddMTI("2200")
parser.SetField(3, "100700")
parser.SetField(4, "1500")
parser.SetField(5, "5")
parser.SetField(7, "0711170215")

// Pack the message
isoMsg, err := parser.MarshalString()
if err != nil {
    panic(err)
}

// Send packed message to the server
// ...

//Unpack message
if err := parser.Unmarshal([]byte(isoMsg)); err != nil {
    panic(err)
}
```

