package ir

// Typed λ-calculus

type (
	// Data :
	Data interface{ isData() }

	// Type :
	Type interface{ isType() }

	// Kind :
	Kind interface{ isKind() }
)

type (
	// DataVar :
	DataVar struct{ Type Type }

	// TypeVar :
	TypeVar struct{ Kind Kind }

	// SetKind :
	SetKind struct{}
)

func (*DataVar) isData() {}
func (*TypeVar) isType() {}
func (*SetKind) isKind() {}

type (
	// AllDD :
	AllDD struct {
		Vars []DataVar
		Term Type
	}

	// AllTD :
	AllTD struct {
		Vars []TypeVar
		Term Type
	}

	// AllTT :
	AllTT struct {
		Vars []TypeVar
		Term Kind
	}
)

func (*AllDD) isType() {}
func (*AllTD) isType() {}
func (*AllTT) isKind() {}

type (
	// AppDD :
	AppDD struct {
		Args []Data
		To   Data
	}

	// AppTD :
	AppTD struct {
		Args []Type
		To   Data
	}

	// AppTT :
	AppTT struct {
		Args []Type
		To   Type
	}
)

func (*AppDD) isData() {}
func (*AppTD) isData() {}
func (*AppTT) isType() {}

type (
	// AbsDD :
	AbsDD struct {
		Vars []DataVar
		Term Data
	}

	// AbsTD :
	AbsTD struct {
		Vars []TypeVar
		Term Data
	}

	// AbsTT :
	AbsTT struct {
		Vars []TypeVar
		Term Type
	}
)

func (*AbsDD) isData() {}
func (*AbsTD) isData() {}
func (*AbsTT) isType() {}

// Scilla
type (
	// IntType :
	IntType struct{ Size int }

	// NatType :
	NatType struct{ Size int }

	// RawType :
	RawType struct{ Size int }

	// StrType :
	StrType struct{}

	// BnrType :
	BnrType struct{}

	// ExcType :
	ExcType struct{}

	// MsgType :
	MsgType struct{}

	// MapType :
	MapType struct{ Key, Val Type }
)

func (*IntType) isType() {}
func (*NatType) isType() {}
func (*RawType) isType() {}
func (*StrType) isType() {}
func (*BnrType) isType() {}
func (*MsgType) isType() {}
func (*ExcType) isType() {}
func (*MapType) isType() {}

type (
	// Int :
	Int struct {
		Type *IntType
		Data string
	}

	// Nat :
	Nat struct {
		Type *NatType
		Data string
	}

	// Raw :
	Raw struct {
		Type *RawType
		Data string
	}

	// Str :
	Str struct {
		Type *StrType
		Data string
	}

	// Bnr :
	Bnr struct {
		Type *BnrType
		Data string
	}

	// Exc :
	Exc struct {
		Type *ExcType
		Data string
	}

	// Msg :
	Msg struct {
		Type *MsgType
		Data map[string]string
	}

	// Map :
	Map struct {
		Type *MapType
		Data map[string]string
	}
)

func (*Int) isData() {}
func (*Nat) isData() {}
func (*Raw) isData() {}
func (*Str) isData() {}
func (*Bnr) isData() {}
func (*Exc) isData() {}
func (*Msg) isData() {}
func (*Map) isData() {}
