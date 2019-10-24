package ir

// ProcType :
type ProcType struct{ Vars []*DataVar }

func (*ProcType) isType() {}

// Proc :
type Proc struct {
	Vars    []Type
	ProcTyp *ProcType
	Plan    []Unit
	Jump    Jump
}

func (*Proc) isData()      {}
func (x *Proc) Type() Type { return x.ProcTyp }

type (
	// Jump :
	Jump interface{ isJump() }

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
type Unit interface{ isUnit() }

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
