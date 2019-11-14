package ir

import ()

type BuiltinADTs struct {
	star Kind

	Boolean Type
	TT      Data
	FF      Data

	List *AbsTT
	Nil  *AbsTD
	Cons *AbsTD

	Option *AbsTT
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

		listNil AbsTD

		listCons      AbsTD
		listConsAbsDD AbsDD

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
		"True":  {},
		"False": {},
	}
	tt = Enum{
		EnumType: &boolean,
		Case:     "True",
	}
	ff = Enum{
		EnumType: &boolean,
		Case:     "False",
	}

	list = AbsTT{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &listEnum,
	}
	listEnum = EnumType{
		"Nil":  {},
		"Cons": {&list.Vars[0], &listEnum},
	}

	listNil = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
	}
	listNil.Term = &AbsDD{
		Vars: []DataVar{},
		Term: &Enum{
			EnumType: &AppTT{
				Args: []Type{&listNil.Vars[0]},
				To:   &list,
			},
			Case: "Nil",
			Data: []Data{},
		},
	}

	listCons = AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: &star},
		},
		Term: &listConsAbsDD,
	}
	listType := &AppTT{
		Args: []Type{&listCons.Vars[0]},
		To:   &list,
	}
	listConsAbsDD = AbsDD{
		Vars: []DataVar{
			DataVar{DataType: &listCons.Vars[0]},
			DataVar{
				DataType: listType,
			},
		},
	}
	listConsAbsDD.Term = &Enum{
		EnumType: listType,
		//&AppTT{
		//Args: []Type{&listCons.Vars[0]},
		//To:   &list,
		//},
		Case: "Cons",
		Data: []Data{&listConsAbsDD.Vars[0], &listConsAbsDD.Vars[1]},
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
		Vars: []DataVar{DataVar{DataType: &some.Vars[0]}},
	}
	someAbsDD.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&some.Vars[0]},
			To:   &option,
		},
		Case: "some",
		Data: []Data{&someAbsDD.Vars[0]},
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
		Vars: []DataVar{
			DataVar{DataType: &product.Vars[0]},
			DataVar{DataType: &product.Vars[1]},
		},
	}
	pairAbsDD.Term = &Enum{
		EnumType: &AppTT{
			Args: []Type{&product.Vars[0], &product.Vars[1]},
			To:   &product,
		},
		Case: "pair",
		Data: []Data{&pairAbsDD.Vars[0], &pairAbsDD.Vars[1]},
	}

	return BuiltinADTs{
		star: &star,

		Boolean: &boolean,
		TT:      &tt,
		FF:      &ff,

		List: &list,
		Nil:  &listNil,
		Cons: &listCons,

		Option: &option,
		None:   &none,
		Some:   &some,

		Product: &product,
		Pair:    &pair,
	}
}
