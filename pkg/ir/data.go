package ir

// Typed λ-calculus

type (
	Node interface {
		ID() uint64
	}

	IDNode struct {
		Id uint64
	}

	// Type :
	Type interface {
		Node
		isType()
	}
	// Data :
	Data interface {
		Node
		isData()
	}
	// Kind :
	Kind interface {
		Node
		isKind()
	}

	// DataVar :
	DataVar struct {
		IDNode
		DataType Type
	}
	// TypeVar :
	TypeVar struct {
		IDNode
		Kind Kind
	}
	// SetKind :
	SetKind struct {
		IDNode
	}

	//Builtin
	Builtin struct {
		IDNode
		BuiltinType Type
	}
)

func (i *IDNode) ID() uint64 {
	return i.Id
}

func (*DataVar) isData() {}
func (*Builtin) isData() {}
func (*TypeVar) isType() {}
func (*SetKind) isKind() {}

type (
	// AllDD :
	AllDD struct {
		IDNode
		Vars []DataVar
		Term Type
	}

	// AllTD :
	AllTD struct {
		IDNode
		Vars []TypeVar
		Term Type
	}

	// AllTT :
	AllTT struct {
		IDNode
		Vars []TypeVar
		Term Kind
	}

	// AppDD :
	AppDD struct {
		IDNode
		Args []Data
		To   Data
	}

	// AppTD :
	AppTD struct {
		IDNode
		Args []Type
		To   Data
	}

	// AppTT :
	AppTT struct {
		IDNode
		Args []Type
		To   Type
	}

	// AbsDD :
	AbsDD struct {
		IDNode
		Vars []DataVar
		Term Data
	}

	// AbsTD :
	AbsTD struct {
		IDNode
		Vars []TypeVar
		Term Data
	}

	// AbsTT :
	AbsTT struct {
		IDNode
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
	IntType struct {
		IDNode
		Size int
	}

	// Uint
	NatType struct {
		IDNode
		Size int
	}

	// ByStr
	RawType struct {
		IDNode
		Size int
	}

	// StrType :
	StrType struct {
		IDNode
	}

	// BnrType :
	BnrType struct {
		IDNode
	}

	// ExcType :
	ExcType struct {
		IDNode
	}

	// MsgType :
	MsgType struct {
		IDNode
	}

	// MapType :
	MapType struct {
		IDNode
		KeyType Type
		ValType Type
	}
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
		IDNode
		IntType *IntType
		Data    string
	}

	// Nat :
	Nat struct {
		IDNode
		NatType *NatType
		Data    string
	}

	// Raw :
	Raw struct {
		IDNode
		RawType *RawType
		Data    string
	}

	// Str :
	Str struct {
		IDNode
		StrType *StrType
		Data    string
	}

	// Bnr :
	Bnr struct {
		IDNode
		BnrType *BnrType
		Data    string
	}

	// Exc :
	Exc struct {
		IDNode
		ExcType *ExcType
		Data    string
	}

	// Msg :
	Msg struct {
		IDNode
		MsgType *MsgType
		Data    map[string]Data
	}

	// Map :
	Map struct {
		IDNode
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
type EnumType struct {
	IDNode
	Constructors map[string][]Type
}

func (*EnumType) isType() {}

// Enum :
type Enum struct {
	IDNode
	EnumType Type
	Case     string
	Data     []Data
}

func (*Enum) isData() {}

type (
	// Bind :
	Bind struct {
		IDNode
		BindType Type
		Cond     *Cond
	}

	// Cond :
	Cond struct {
		IDNode
		Case string
		Data []Bind
	}
)

type (
	// PickData :
	PickData struct {
		IDNode
		From Data
		With []DataCase
	}

	// DataCase :
	DataCase struct {
		IDNode
		Bind Bind
		Body Data
	}
)

func (*PickData) isData() {}
func (*Bind) isData()     {}

// ProcType :
type ProcType struct {
	IDNode
	Vars []DataVar
}

func (*ProcType) isType() {}

// Proc :
type Proc struct {
	IDNode
	ProcName string
	Vars     []DataVar
	Plan     []Unit
	Jump     Jump
}

func (*Proc) isData() {}

type (
	// Jump :
	Jump interface {
		Node
		isJump()
	}

	// CallProc :
	CallProc struct {
		IDNode
		Args []Data
		To   Data
	}

	// PickProc :
	PickProc struct {
		IDNode
		From Data
		With []ProcCase
	}

	// ProcCase :
	ProcCase struct {
		IDNode
		Bind Bind
		Body Proc
	}
)

func (*CallProc) isJump() {}
func (*PickProc) isJump() {}

// Unit :
type Unit interface {
	Node
	isUnit()
}

type (
	// Load :
	Load struct {
		IDNode
		Slot string
		Path []Data
	}

	// Save :
	Save struct {
		IDNode
		Slot string
		Path []Data
		Data Data
	}

	// Event :
	Event struct {
		IDNode
		Data Data
	}

	// Send :
	Send struct {
		IDNode
		Data Data
	}

	// Accept :
	Accept struct {
		IDNode
	}
)

func (*Load) isData() {}

func (*Load) isUnit()   {}
func (*Save) isUnit()   {}
func (*Event) isUnit()  {}
func (*Send) isUnit()   {}
func (*Accept) isUnit() {}

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
