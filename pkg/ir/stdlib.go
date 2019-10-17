package ir

import ()

type BuiltinADTs struct {
	star Kind

	Boolean Type
	Tt      Data
	Ff      Data

	List    Type
	Empty   Data
	Stack   Data
	isEmpty Data

	Option Type
	None   Data
	Some   Data

	Pair   Type
	Pair_c Data
}

func StdLib() BuiltinADTs {
	var (
		star SetKind

		boolean EnumType
		tt      Enum
		ff      Enum

		list     AbsTT
		listEnum EnumType

		empty AbsTD

		stack      AbsTD
		stackAbsDD AbsDD

		isEmpty      AbsTD
		isEmptyAbsDD AbsDD

		option     AbsTT
		optionEnum EnumType

		none      AbsTD
		some      AbsTD
		someAbsDD AbsDD

		pair     AbsTT
		pairEnum EnumType

		pair_c    AbsTD
		pairAbsDD AbsDD
	)

	boolean = EnumType{
		"tt": {},
		"ff": {},
	}
	tt = Enum{
		EnumType: &boolean,
		Case:     "tt",
	}
	ff = Enum{
		EnumType: &boolean,
		Case:     "ff",
	}

	list = AbsTT{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &listEnum,
	}
	listEnum = EnumType{
		"empty": {},
		"stack": {&list.Vars[0], &listEnum},
	}

	empty = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
	}
	empty.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&empty.Vars[0]},
			To:   &list,
		},
		Case: "empty",
		Data: []Data{},
	}

	stack = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &stackAbsDD,
	}
	stackAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &stack.Vars[0]},
			&DataVar{
				DataType: &AppTT{
					Args: []Type{&stack.Vars[0]},
					To:   &list,
				},
			},
		},
	}
	stackAbsDD.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&stack.Vars[0]},
			To:   &list,
		},
		Case: "stack",
		Data: []Data{stackAbsDD.Vars[0], stackAbsDD.Vars[1]},
	}

	isEmpty = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &isEmptyAbsDD,
	}
	isEmptyAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: &AppTT{
					Args: []Type{&isEmpty.Vars[0]},
					To:   &list,
				},
			},
		},
	}
	isEmptyAbsDD.Term = &PickData{
		From: isEmptyAbsDD.Vars[0],
		With: []DataCase{
			DataCase{
				Bind: &Bind{
					BindType: isEmptyAbsDD.Vars[0].Type(),
					When:     &When{Case: "empty", Data: []*Bind{}},
				},
				Body: &tt,
			},
			DataCase{
				Bind: &Bind{
					BindType: isEmptyAbsDD.Vars[0].Type(),
					When: &When{Case: "stack", Data: []*Bind{
						&Bind{BindType: &isEmpty.Vars[0]},
						&Bind{BindType: isEmptyAbsDD.Vars[0].Type()},
					}},
				},
				Body: &ff,
			},
		},
	}

	option = AbsTT{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &optionEnum,
	}
	optionEnum = EnumType{
		"none": {},
		"some": {&option.Vars[0]},
	}

	none = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
	}
	none.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&none.Vars[0]},
			To:   &option,
		},
		Case: "none",
		Data: []Data{},
	}

	some = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &someAbsDD,
	}
	someAbsDD = AbsDD{
		Vars: []*DataVar{&DataVar{DataType: &some.Vars[0]}},
	}
	someAbsDD.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&some.Vars[0]},
			To:   &option,
		},
		Case: "some",
		Data: []Data{someAbsDD.Vars[0]},
	}

	pair = AbsTT{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
			TypeVar{Kind: &star},
		},
		Term: &pairEnum,
	}

	pairEnum = EnumType{
		"pair": {&pair.Vars[0], &pair.Vars[1]},
	}

	pair_c = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
			TypeVar{Kind: &star},
		},
		Term: &pairAbsDD,
	}

	pairAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &pair.Vars[0]},
			&DataVar{DataType: &pair.Vars[1]},
		},
	}
	pairAbsDD.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&pair.Vars[0], &pair.Vars[1]},
			To:   &pair,
		},
		Case: "pair",
		Data: []Data{pairAbsDD.Vars[0], pairAbsDD.Vars[1]},
	}

	return BuiltinADTs{
		star: &star,

		Boolean: &boolean,
		Tt:      &tt,
		Ff:      &ff,

		List:    &list,
		Empty:   &empty,
		Stack:   &stack,
		isEmpty: &isEmpty,

		Option: &option,
		None:   &none,
		Some:   &some,

		Pair:   &pair,
		Pair_c: &pair_c,
	}
}
