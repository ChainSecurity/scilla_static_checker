package ir

// EnumType :
type EnumType map[string][]Type

func (*EnumType) isType() {}

// Enum :
type Enum struct {
	EnumType Type
	Case     string
	Data     []Data
}

func (*Enum) isData() {}

type (
	// Bind :
	Bind struct {
		BindType Type
		Cond     *Cond
	}

	// Cond :
	Cond struct {
		Case string
		Data []*Bind
	}
)

type (
	// PickData :
	PickData struct {
		From Data
		With []DataCase
	}

	// DataCase :
	DataCase struct {
		Bind *Bind
		Body Data
	}
)

func (*PickData) isData() {}
