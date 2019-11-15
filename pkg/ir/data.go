package ir

// Typed λ-calculus

type (
	IRNode interface{}

	// Type :
	Type interface {
		IRNode
		isType()
	}
	// Data :
	Data interface {
		IRNode
		isData()
	}
	// Kind :
	Kind interface {
		IRNode
		isKind()
	}

	// DataVar :
	DataVar struct{ DataType Type }
	// TypeVar :
	TypeVar struct{ Kind Kind }
	// SetKind :
	SetKind struct{}

	//Builtin
	Builtin struct{ BuiltinType Type }
)

func (*DataVar) isData() {}
func (*Builtin) isData() {}
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

func (*AllDD) isType() {}
func (*AllTD) isType() {}
func (*AllTT) isKind() {}
func (*AppDD) isData() {}
func (*AppTD) isData() {}
func (*AppTT) isType() {}
func (*AbsDD) isData() {}
func (*AbsTD) isData() {}
func (*AbsTT) isType() {}

// Scilla Types
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
		Data    map[string]Data
	}

	// Map :
	Map struct {
		MapType *MapType
		Data    map[string]string
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
		Data []Bind
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
		Bind Bind
		Body Data
	}
)

func (*PickData) isData() {}
func (*Bind) isData()     {}

// ProcType :
type ProcType struct {
	Vars []DataVar
}

func (*ProcType) isType() {}

// Proc :
type Proc struct {
	Vars []DataVar
	Plan []Unit
	Jump Jump
}

func (*Proc) isData() {}

type (
	// Jump :
	Jump interface {
		IRNode
		isJump()
	}

	// CallProc :
	CallProc struct {
		Args []Data
		To   Data
	}

	// PickProc :
	PickProc struct {
		From Data
		With []ProcCase
	}

	// ProcCase :
	ProcCase struct {
		Bind Bind
		Body Proc
	}
)

func (*CallProc) isJump() {}
func (*PickProc) isJump() {}

// Unit :
type Unit interface {
	IRNode
	isUnit()
}

type (
	// Load :
	Load struct {
		Slot string
		Path []Data
	}

	// Save :
	Save struct {
		Slot string
		Path []Data
		Data Data
	}

	// Emit :
	Emit struct {
		Data Data
	}

	// Send :
	Send struct {
		Data Data
	}

	// Have :
	Have struct{}
)

func (*Load) isData() {}

func (*Load) isUnit() {}
func (*Save) isUnit() {}
func (*Emit) isUnit() {}
func (*Send) isUnit() {}
func (*Have) isUnit() {}

func (*AbsTD) isUnit() {}
func (*AbsDD) isUnit() {}
func (*AppDD) isUnit() {}
func (*AppTD) isUnit() {}

func (*Int) isUnit() {}
func (*Nat) isUnit() {}
func (*Raw) isUnit() {}
func (*Str) isUnit() {}
func (*Bnr) isUnit() {}
func (*Exc) isUnit() {}
func (*Msg) isUnit() {}
func (*Map) isUnit() {}
