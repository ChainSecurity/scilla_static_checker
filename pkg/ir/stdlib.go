package ir

import ()

type BuiltinADTs struct {
	star Kind

	Boolean Type
	TT      Data
	FF      Data

	List  Type
	Empty Data
	Stack Data

	Option Type
	None   Data
	Some   Data

	Product Type
	Pair    Data
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

		option     AbsTT
		optionEnum EnumType

		none      AbsTD
		some      AbsTD
		someAbsDD AbsDD

		product  AbsTT
		pairEnum EnumType

		pair      AbsTD
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

	product = AbsTT{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
			TypeVar{Kind: &star},
		},
		Term: &pairEnum,
	}

	pairEnum = EnumType{
		"pair": {&product.Vars[0], &product.Vars[1]},
	}

	pair = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
			TypeVar{Kind: &star},
		},
		Term: &pairAbsDD,
	}

	pairAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &product.Vars[0]},
			&DataVar{DataType: &product.Vars[1]},
		},
	}
	pairAbsDD.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&product.Vars[0], &product.Vars[1]},
			To:   &product,
		},
		Case: "pair",
		Data: []Data{pairAbsDD.Vars[0], pairAbsDD.Vars[1]},
	}

	return BuiltinADTs{
		star: &star,

		Boolean: &boolean,
		TT:      &tt,
		FF:      &ff,

		List:  &list,
		Empty: &empty,
		Stack: &stack,

		Option: &option,
		None:   &none,
		Some:   &some,

		Product: &product,
		Pair:    &pair,
	}
}
