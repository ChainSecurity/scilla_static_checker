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
	None   *AbsTD
	Some   *AbsTD

	Product Type
	Pair    Data
}

func StdLib(builder *CFGBuilder) BuiltinADTs {
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

	star = SetKind{
		IDNode: builder.newIDNode(),
	}
	boolean = EnumType{
		IDNode: builder.newIDNode(),
		Constructors: map[string][]Type{
			"True":  {},
			"False": {},
		},
	}
	tt = Enum{
		IDNode:   builder.newIDNode(),
		EnumType: &boolean,
		Case:     "True",
	}
	ff = Enum{
		IDNode:   builder.newIDNode(),
		EnumType: &boolean,
		Case:     "False",
	}

	list = AbsTT{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
		Term: &listEnum,
	}
	listEnum = EnumType{
		IDNode: builder.newIDNode(),
		Constructors: map[string][]Type{
			"Nil":  {},
			"Cons": {&list.Vars[0], &listEnum},
		},
	}

	listNil = AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
	}
	listNil.Term = &AbsDD{
		IDNode: builder.newIDNode(),
		Vars:   []DataVar{},
		Term: &Enum{
			IDNode: builder.newIDNode(),
			EnumType: &AppTT{
				IDNode: builder.newIDNode(),
				Args:   []Type{&listNil.Vars[0]},
				To:     &list,
			},
			Case: "Nil",
			Data: []Data{},
		},
	}

	listCons = AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
		Term: &listConsAbsDD,
	}
	listType := &AppTT{
		IDNode: builder.newIDNode(),
		Args:   []Type{&listCons.Vars[0]},
		To:     &list,
	}
	listConsAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &listCons.Vars[0],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: listType,
			},
		},
	}
	listConsAbsDD.Term = &Enum{
		IDNode:   builder.newIDNode(),
		EnumType: listType,
		//&AppTT{
		//Args: []Type{&listCons.Vars[0]},
		//To:   &list,
		//},
		Case: "Cons",
		Data: []Data{&listConsAbsDD.Vars[0], &listConsAbsDD.Vars[1]},
	}

	option = AbsTT{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
		Term: &optionEnum,
	}
	optionEnum = EnumType{
		IDNode: builder.newIDNode(),
		Constructors: map[string][]Type{
			"None": {},
			"Some": {&option.Vars[0]},
		},
	}

	none = AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
	}
	none.Term = &Enum{
		IDNode: builder.newIDNode(),
		EnumType: &AppTT{
			IDNode: builder.newIDNode(),
			Args:   []Type{&none.Vars[0]},
			To:     &option,
		},
		Case: "None",
		Data: []Data{},
	}

	some = AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
		Term: &someAbsDD,
	}
	someAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars:   []DataVar{DataVar{DataType: &some.Vars[0]}},
	}
	someAbsDD.Term = &Enum{
		IDNode: builder.newIDNode(),
		EnumType: &AppTT{
			IDNode: builder.newIDNode(),
			Args:   []Type{&some.Vars[0]},
			To:     &option,
		},
		Case: "Some",
		Data: []Data{&someAbsDD.Vars[0]},
	}

	product = AbsTT{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
		Term: &pairEnum,
	}

	pairEnum = EnumType{
		IDNode: builder.newIDNode(),
		Constructors: map[string][]Type{
			"pair": {&product.Vars[0], &product.Vars[1]},
		},
	}

	pair = AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   &star,
			},
		},
		Term: &pairAbsDD,
	}

	pairAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &product.Vars[0],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &product.Vars[1],
			},
		},
	}
	pairAbsDD.Term = &Enum{
		IDNode: builder.newIDNode(),
		EnumType: &AppTT{
			IDNode: builder.newIDNode(),
			Args:   []Type{&product.Vars[0], &product.Vars[1]},
			To:     &product,
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
