package ir

// Typed λ-calculus

type (
	// Type :
	Type interface{ isType() }

	// Data :
	Data interface {
		isData()
		Type() Type
	}

	// Kind :
	Kind interface{ isKind() }
)

type (
	// DataVar :
	DataVar struct{ DataType Type }

	// TypeVar :
	TypeVar struct{ Kind Kind }

	// SetKind :
	SetKind struct{}

	//Builtin
	Builtin struct{ BuiltinType Type }
)

func (*DataVar) isData()      {}
func (x *DataVar) Type() Type { return x.DataType }

func (*Builtin) isData()      {}
func (x *Builtin) Type() Type { return x.BuiltinType }

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

func (*AppDD) isData()      {}
func (a *AppDD) Type() Type { return a.To.Type() }
func (*AppTD) isData()      {}
func (a *AppTD) Type() Type { return a.To.Type() }
func (*AppTT) isType()      {}

type (
	// AbsDD :
	AbsDD struct {
		Vars []*DataVar
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

func (*AbsDD) isData()      {}
func (a *AbsDD) Type() Type { return a.Term.Type() }
func (*AbsTD) isData()      {}
func (a *AbsTD) Type() Type { return a.Term.Type() }
func (*AbsTT) isType()      {}

// Scilla
type (
	// Int
	IntType struct{ Size int }

	// Uint
	NatType struct{ Size int }

	// ByStr
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
	MapType struct{ KeyType, ValType Type }
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
		IntType *IntType
		Data    string
	}

	// Nat :
	Nat struct {
		NatType *NatType
		Data    string
	}

	// Raw :
	Raw struct {
		RawType *RawType
		Data    string
	}

	// Str :
	Str struct {
		StrType *StrType
		Data    string
	}

	// Bnr :
	Bnr struct {
		BnrType *BnrType
		Data    string
	}

	// Exc :
	Exc struct {
		ExcType *ExcType
		Data    string
	}

	// Msg :
	Msg struct {
		MsgType *MsgType
		Data    map[string]string
	}

	// Map :
	Map struct {
		MapType *MapType
		Data    map[string]string
	}
)

func (*Int) isData()      {}
func (*Nat) isData()      {}
func (*Raw) isData()      {}
func (*Str) isData()      {}
func (*Bnr) isData()      {}
func (*Exc) isData()      {}
func (*Msg) isData()      {}
func (*Map) isData()      {}
func (x *Int) Type() Type { return x.IntType }
func (x *Nat) Type() Type { return x.NatType }
func (x *Raw) Type() Type { return x.RawType }
func (x *Str) Type() Type { return x.StrType }
func (x *Bnr) Type() Type { return x.BnrType }
func (x *Exc) Type() Type { return x.ExcType }
func (x *Msg) Type() Type { return x.MsgType }
func (x *Map) Type() Type { return x.MapType }
