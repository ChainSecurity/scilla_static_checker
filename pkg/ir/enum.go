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

func (*Enum) isData()      {}
func (x *Enum) Type() Type { return x.EnumType }

type (
	// Bind :
	Bind struct {
		BindType Type
		When     *When
	}

	// When :
	When struct {
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

func (*PickData) isData()       {}
func (pd *PickData) Type() Type { return pd.With[0].Body.Type() }
