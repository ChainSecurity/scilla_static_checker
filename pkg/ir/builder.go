package ir

import (
	"errors"
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"strconv"
)

type CFGBuilder struct {
	builtinOpMap map[string]Data

	constrTypeMap   map[string]string
	intWidthTypeMap map[int]*IntType
	natWidthTypeMap map[int]*NatType
	typeMap         map[string]Type
	varStack        map[string][]Data
	fieldStack      map[string][]Data
	constructor     *Proc
}

func (builder *CFGBuilder) visitCtr(ctr *ast.CtrDef) (string, []Type) {
	name := ctr.CDName.Id
	var types []Type
	for _, typ := range ctr.CArgTypes {
		typ := builder.visitASTType(typ)
		types = append(types, typ)
	}
	return name, types
}

func (builder *CFGBuilder) visitPattern(p ast.Pattern, t Type) *Bind {
	switch pat := p.(type) {
	case *ast.WildcardPattern:
		return &Bind{BindType: t}
	case *ast.BinderPattern:
		return &Bind{BindType: t}
	case *ast.ConstructorPattern:
		var typeList []Type
		switch typ := t.(type) {
		case *EnumType:
			typeList = (*typ)[pat.ConstrName]
		default:
			panic(errors.New(fmt.Sprintf("Not constructor pattern type: %T", t)))
		}
		if len(typeList) != len(pat.Pats) {
			panic(errors.New(fmt.Sprintf("Constructor pattern argument length mistmatch: %s", pat.ConstrName)))
		}
		whenData := []*Bind{}
		for i, subp := range pat.Pats {
			whenData = append(whenData, builder.visitPattern(subp, typeList[i]))
		}
		return &Bind{
			BindType: t,
			Cond: &Cond{
				Case: pat.ConstrName,
				Data: whenData,
			},
		}
	default:
		panic(errors.New("Unknown pattern "))
	}
}

func (builder *CFGBuilder) visitLiteral(l ast.Literal) Data {
	switch lit := l.(type) {
	case *ast.StringLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	case *ast.BNumLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	case *ast.ByStrLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	case *ast.ByStrXLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	case *ast.IntLiteral:
		typ, ok := builder.intWidthTypeMap[lit.Width]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown width: %d", lit.Width)))
		}
		return &Int{
			IntType: typ, Data: lit.Val,
		}
	case *ast.UintLiteral:
		typ, ok := builder.natWidthTypeMap[lit.Width]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown width: %d", lit.Width)))
		}
		return &Nat{
			NatType: typ, Data: lit.Val,
		}
	case *ast.MapLiteral:
		ktyp := builder.visitASTType(lit.KeyType)
		vtyp := builder.visitASTType(lit.ValType)
		maptyp := MapType{ktyp, vtyp}
		return &Map{
			MapType: &maptyp,
			Data:    map[string]string{},
		}

	case *ast.ADTValLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	}
	return nil
}

func (builder *CFGBuilder) visitASTType(e ast.ASTType) Type {

	switch n := e.(type) {
	case *ast.PrimType:
		t, ok := builder.typeMap[n.Name]
		if !ok {
			panic(errors.New(fmt.Sprintf("PrimType not found : %s", n.Name)))
		}
		return t
	//case *ast.MapType:
	//fmt.Printf("%T", n)
	case *ast.ADT:
		typ, ok := builder.typeMap[n.Name]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown type: %s", n.Name)))
		}
		return typ
	//case *ast.FunType:
	//fmt.Printf("%T", n)
	//case *ast.TypeVar:
	//fmt.Printf("%T", n)
	//case *ast.PolyFun:
	//fmt.Printf("%T", n)
	//case *ast.Unit:
	//fmt.Printf("%T", n)
	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func (builder *CFGBuilder) visitExpression(e ast.Expression) Data {
	switch n := e.(type) {
	case *ast.ConstrExpression:
		constrName := n.ConstructorName
		typeName, ok := builder.constrTypeMap[constrName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown constructor: %s", constrName)))
		}

		for k := range builder.constrTypeMap {
			fmt.Println("->", k, builder.constrTypeMap[k])
		}
		fmt.Println("____")
		typ, ok := builder.typeMap[typeName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown type: %s", typeName)))
		}
		d := []Data{}
		//TODO handle arguments
		//TODO do a check on types
		//for _, arg := range n.Args {
		//fmt.Println(arg.Id)
		//d = d.append(builderu)
		//}
		res := Enum{
			EnumType: typ,
			Case:     constrName,
			Data:     d,
		}
		return &res
	case *ast.FunExpression:
		lhs := n.Lhs.Id
		lhsTyp := builder.visitASTType(n.LhsType)
		dataVar := DataVar{lhsTyp}
		stackMapPush(builder.varStack, lhs, &dataVar)
		defer stackMapPop(builder.varStack, lhs)

		rhs := builder.visitExpression(n.RhsExpr)

		return &AbsDD{Vars: []*DataVar{&dataVar}, Term: rhs}
	case *ast.MatchExpression:
		lhs := n.Lhs.Id

		data, ok := stackMapPeek(builder.varStack, lhs)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", lhs)))
		}
		cases := []DataCase{}
		for _, c := range n.Cases {
			b := builder.visitPattern(c.Pat, data.Type())
			e := builder.visitExpression(c.Expr)
			mec := DataCase{Bind: b, Body: e}
			cases = append(cases, mec)
		}
		return &PickData{
			From: data, With: cases,
		}
	case *ast.LiteralExpression:
		return builder.visitLiteral(n.Val)
	case *ast.VarExpression:
		data, ok := stackMapPeek(builder.varStack, n.Var.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Var.Id)))
		}
		return data
	case *ast.BuiltinExpression:
		opName := n.BuintinFunc.BuiltinOp
		vars := make([]Data, len(n.Args))
		for i, a := range n.Args {
			v, ok := stackMapPeek(builder.varStack, a.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", a.Id)))
			}
			vars[i] = v
		}
		op := builder.builtinOpMap[opName]
		if opName == "concat" {
			v0 := vars[0].Type()
			switch raw0 := v0.(type) {
			case *RawType:
				v1 := vars[1].Type()
				raw1, ok := v1.(*RawType)
				if !ok {
					panic(errors.New(fmt.Sprintf("Builtin concat wrong type: %T", v1)))
				}
				resType := "ByStr" + strconv.Itoa(raw0.Size+raw1.Size)
				op = &AbsDD{
					Vars: []*DataVar{
						&DataVar{
							DataType: v0,
						},
						&DataVar{
							DataType: v1,
						},
					},
					Term: &Builtin{setDefaultType(builder.typeMap, resType, &RawType{raw0.Size + raw1.Size})},
				}
			case *StrType:
				op = builder.builtinOpMap["str_concat"]
			default:
				panic(errors.New(fmt.Sprintf("Builtin concat wrong type: %T", v0)))
			}
		}
		res := AppDD{
			Args: vars,
			To: &AppTD{
				Args: []Type{vars[0].Type()},
				To:   op,
			},
		}
		return &res
	case *ast.LetExpression:
		varName := n.Var.Id
		expr := builder.visitExpression(n.Expr)
		fmt.Println("LetExpression", varName, expr.Type())
		stackMapPush(builder.varStack, varName, expr)
		defer stackMapPop(builder.varStack, varName)
		body := builder.visitExpression(n.Body)
		return body
	default:
		fmt.Printf("Unhandled Expression type: %T\n", n)
		//panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
		return nil
	}
}

func (builder *CFGBuilder) visitLibEntry(le ast.LibEntry) {
	switch n := le.(type) {
	case *ast.LibraryVariable:
		name := n.Name.Id
		v := builder.visitExpression(n.Expr)
		stackMapPush(builder.varStack, name, v)
	case *ast.LibraryType:
		typeName := n.Name.Id
		typ := EnumType{}
		for _, ctr := range n.CtrDefs {
			constrName, types := builder.visitCtr(ctr)
			typ[constrName] = types
			builder.constrTypeMap[constrName] = typeName
		}
		builder.typeMap[typeName] = &typ
	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func (builder *CFGBuilder) visitParams(p *ast.Parameter) (string, Type) {
	name := p.Name.Id
	t := builder.visitASTType(p.Type)
	return name, t
}

func (builder *CFGBuilder) visitField(f *ast.Field) (string, Data) {
	name := f.Name.Id
	//t := builder.typeMap[f.Type]
	data := builder.visitExpression(f.Expr)
	stackMapPush(builder.varStack, name, data)
	return name, data
}

func (builder *CFGBuilder) visitComponent(comp *ast.Component) {
	fmt.Println(comp.Name.Id)
}

func (builder *CFGBuilder) Visit(node ast.AstNode) ast.Visitor {
	switch n := node.(type) {
	case ast.LibEntry:
		builder.visitLibEntry(n)
	case *ast.Contract:
		name := n.Name.Id
		fmt.Println("Contract", name)
		paramNames := make([]string, len(n.Params))
		params := map[string]Type{}
		dataVars := make([]*DataVar, len(n.Params))
		for i, p := range n.Params {
			pName, pType := builder.visitParams(p)
			params[pName] = pType
			paramNames[i] = pName
			dataVar := DataVar{pType}
			dataVars[i] = &dataVar
			stackMapPush(builder.varStack, pName, &dataVar)
		}

		builder.constructor = &Proc{}
		builder.constructor.Vars = dataVars
		builder.constructor.Plan = make([]Unit, len(n.Params))
		for i, f := range n.Fields {
			n, d := builder.visitField(f)
			stackMapPush(builder.fieldStack, n, d)
			builder.constructor.Plan[i] = &Save{n, []Data{}, d}
		}

		for _, pName := range paramNames {
			stackMapPop(builder.varStack, pName)
		}

		//for

		//t := n.Type
		//expr := builder.visitExpression(n.Expr)
		//builder.Field[name] = Save{}
		//builder.visitLibEntry(n)
	case *ast.Component:
		builder.visitComponent(n)
	default:
		//fmt.Printf("Unhandled type: %T\n", n)
		// do nothing
	}
	return builder
}

func (builder *CFGBuilder) initPrimitiveTypes() {
	stdLib := StdLib()

	sizes := []int{32, 64, 128, 256}
	for _, s := range sizes {
		intName := "Int" + strconv.Itoa(s)
		intTyp := IntType{s}
		uintName := "Uint" + strconv.Itoa(s)
		uintTyp := NatType{s}
		builder.intWidthTypeMap[s] = &intTyp
		builder.natWidthTypeMap[s] = &uintTyp
		builder.typeMap[intName] = &intTyp
		builder.typeMap[uintName] = &uintTyp

		// Conversion functions
		var intAbsDD AbsDD
		intAbsTD := AbsTD{
			Vars: []TypeVar{
				TypeVar{Kind: stdLib.star},
			},
			Term: &intAbsDD,
		}
		intAbsDD = AbsDD{
			Vars: []*DataVar{
				&DataVar{DataType: &intAbsTD.Vars[0]},
			},
			Term: &Builtin{
				&AppTT{
					Args: []Type{&intTyp},
					To:   stdLib.Option,
				},
			},
		}
		intOpName := fmt.Sprintf("to_int%d", s)
		builder.builtinOpMap[intOpName] = &intAbsTD

		var uintAbsDD AbsDD
		uintAbsTD := AbsTD{
			Vars: []TypeVar{
				TypeVar{Kind: stdLib.star},
			},
			Term: &uintAbsDD,
		}
		uintAbsDD = AbsDD{
			Vars: []*DataVar{
				&DataVar{DataType: &uintAbsTD.Vars[0]},
			},
			Term: &Builtin{
				&AppTT{
					Args: []Type{&uintTyp},
					To:   stdLib.Option,
				},
			},
		}
		uintOpName := fmt.Sprintf("to_uint%d", s)
		builder.builtinOpMap[uintOpName] = &uintAbsTD

	}

	builder.typeMap["String"] = &StrType{}
	builder.typeMap["ByStr"] = &RawType{-1}
	builder.typeMap["ByStr32"] = &RawType{32}
	builder.typeMap["ByStr33"] = &RawType{33}
	builder.typeMap["ByStr64"] = &RawType{64}
	builder.typeMap["ByStr20"] = &RawType{20}
	builder.typeMap["BNum"] = &BnrType{}

	builder.typeMap["Bool"] = stdLib.Boolean

	intIntOps := []string{"add", "sub", "mul", "div", "rem"}
	intBoolOps := []string{"eq", "lt"}

	for _, bOp := range intIntOps {
		var bAbsDD AbsDD
		bAbsTD := AbsTD{
			Vars: []TypeVar{
				TypeVar{Kind: stdLib.star},
			},
			Term: &bAbsDD,
		}
		bAbsDD = AbsDD{
			Vars: []*DataVar{
				&DataVar{DataType: &bAbsTD.Vars[0]},
				&DataVar{DataType: &bAbsTD.Vars[0]},
			},
			Term: &Builtin{&bAbsTD.Vars[0]},
		}
		builder.builtinOpMap[bOp] = &bAbsTD
	}

	for _, bOp := range intBoolOps {
		var bAbsDD AbsDD
		bAbsTD := AbsTD{
			Vars: []TypeVar{
				TypeVar{Kind: stdLib.star},
			},
			Term: &bAbsDD,
		}
		bAbsDD = AbsDD{
			Vars: []*DataVar{
				&DataVar{DataType: &bAbsTD.Vars[0]},
				&DataVar{DataType: &bAbsTD.Vars[0]},
			},
			Term: &Builtin{stdLib.Boolean},
		}
		builder.builtinOpMap[bOp] = &bAbsTD
	}

	// builtin pow
	var powAbsDD AbsDD
	powAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
			TypeVar{Kind: stdLib.star},
		},
		Term: &powAbsDD,
	}
	powAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &powAbsTD.Vars[0]},
			&DataVar{DataType: &powAbsTD.Vars[1]},
		},
		Term: &Builtin{&powAbsTD.Vars[0]},
	}
	builder.builtinOpMap["pow"] = &powAbsTD

	// builtin concat
	strConcatOp := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["String"],
			},
			&DataVar{
				DataType: builder.typeMap["String"],
			},
		},
		Term: &Builtin{builder.typeMap["String"]},
	}
	builder.builtinOpMap["str_concat"] = &strConcatOp

	//builtin substr
	strSubstrOp := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["String"],
			},
			&DataVar{
				DataType: builder.natWidthTypeMap[32],
			},
			&DataVar{
				DataType: builder.natWidthTypeMap[32],
			},
		},
		Term: &Builtin{builder.typeMap["String"]},
	}
	builder.builtinOpMap["substr"] = &strSubstrOp

	// builtin to_string
	var toStrAbsDD AbsDD
	toStrAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &toStrAbsDD,
	}
	toStrAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &toStrAbsTD.Vars[0]},
		},
		Term: &Builtin{builder.typeMap["String"]},
	}
	builder.builtinOpMap["to_string"] = &toStrAbsTD

	//builtin strlen

	strlenOp := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["String"],
			},
		},
		Term: &Builtin{builder.natWidthTypeMap[32]},
	}
	builder.builtinOpMap["strlen"] = &strlenOp

	// builtin sha256hash
	var shaAbsDD AbsDD
	shaAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &shaAbsDD,
	}
	shaAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &shaAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.typeMap["ByStr32"],
		},
	}
	builder.builtinOpMap["sha256hash"] = &shaAbsTD

	// builtin keccak256hash
	var keccakAbsDD AbsDD
	keccakAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &keccakAbsDD,
	}
	keccakAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &keccakAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.typeMap["ByStr32"],
		},
	}
	builder.builtinOpMap["keccak256hash"] = &keccakAbsTD

	// builtin ripemd160hash
	var ripemdAbsDD AbsDD
	ripemdAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &ripemdAbsDD,
	}
	ripemdAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &ripemdAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.typeMap["ByStr20"],
		},
	}
	builder.builtinOpMap["ripemd160hash"] = &ripemdAbsTD

	// builtin to_bystr
	var toBystrAbsDD AbsDD
	toBystrAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &toBystrAbsDD,
	}
	toBystrAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &toBystrAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.typeMap["ByStr"],
		},
	}
	builder.builtinOpMap["to_bystr"] = &toBystrAbsTD

	// builtin to_uint256
	var touint256AbsDD AbsDD
	touint256AbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &touint256AbsDD,
	}
	touint256AbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{DataType: &touint256AbsTD.Vars[0]},
		},
		Term: &Builtin{builder.natWidthTypeMap[256]},
	}
	builder.builtinOpMap["to_uint256"] = &touint256AbsTD

	// schnorr_verify
	schnorrVerifyAbsDD := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["ByStr33"],
			},
			&DataVar{
				DataType: builder.typeMap["ByStr"],
			},
			&DataVar{
				DataType: builder.typeMap["ByStr64"],
			},
		},
		Term: &Builtin{builder.typeMap["Bool"]},
	}
	builder.builtinOpMap["schnorr_verify"] = &schnorrVerifyAbsDD

	// ecdsa_verify
	ecdsaVerifyAbsDD := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["ByStr33"],
			},
			&DataVar{
				DataType: builder.typeMap["ByStr"],
			},
			&DataVar{
				DataType: builder.typeMap["ByStr64"],
			},
		},
		Term: &Builtin{builder.typeMap["Bool"]},
	}
	builder.builtinOpMap["ecdsa_verify"] = &ecdsaVerifyAbsDD

	// bech32_to_bystr20
	bech32ToBystr20AbsDD := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["String"],
			},
			&DataVar{
				DataType: builder.typeMap["String"],
			},
		},
		Term: &Builtin{
			&AppTT{
				Args: []Type{builder.typeMap["ByStr20"]},
				To:   stdLib.Option,
			},
		},
	}
	builder.builtinOpMap["bech32_to_bystr20"] = &bech32ToBystr20AbsDD

	// bystr20_to_bech32
	bystr20ToBech32AbsDD := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["String"],
			},
			&DataVar{
				DataType: builder.typeMap["ByStr20"],
			},
		},
		Term: &Builtin{
			&AppTT{
				Args: []Type{builder.typeMap["String"]},
				To:   stdLib.Option,
			},
		},
	}

	builder.builtinOpMap["bystr20_to_bech32"] = &bystr20ToBech32AbsDD

	//Maps

	// put
	var putAbsDD AbsDD
	putAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
			TypeVar{Kind: stdLib.star},
		},
		Term: &putAbsDD,
	}
	putAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				&MapType{
					&putAbsTD.Vars[0],
					&putAbsTD.Vars[1],
				},
			},
			&DataVar{
				&putAbsTD.Vars[0],
			},
			&DataVar{
				&putAbsTD.Vars[1],
			},
		},
	}
	putAbsDD.Term = &Builtin{putAbsDD.Vars[0].DataType}
	builder.builtinOpMap["put"] = &putAbsTD

	// get
	var getAbsDD AbsDD
	getAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
			TypeVar{Kind: stdLib.star},
		},
		Term: &getAbsDD,
	}
	getAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				&MapType{
					&getAbsTD.Vars[0],
					&getAbsTD.Vars[1],
				},
			},
			&DataVar{
				&getAbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			&AppTT{
				Args: []Type{&getAbsTD.Vars[1]},
				To:   stdLib.Option,
			},
		},
	}
	builder.builtinOpMap["get"] = &getAbsTD

	// contains
	var containsAbsDD AbsDD
	containsAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
			TypeVar{Kind: stdLib.star},
		},
		Term: &containsAbsDD,
	}
	containsAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				&MapType{
					&containsAbsTD.Vars[0],
					&containsAbsTD.Vars[1],
				},
			},
			&DataVar{
				&containsAbsTD.Vars[0],
			},
		},
		Term: &Builtin{builder.typeMap["Bool"]},
	}
	builder.builtinOpMap["contains"] = &containsAbsTD

	// remove
	var removeAbsDD AbsDD
	removeAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
			TypeVar{Kind: stdLib.star},
		},
		Term: &removeAbsDD,
	}
	removeAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				&MapType{
					&removeAbsTD.Vars[0],
					&removeAbsTD.Vars[1],
				},
			},
			&DataVar{
				&removeAbsTD.Vars[0],
			},
		},
	}
	removeAbsDD.Term = &Builtin{removeAbsDD.Vars[0].DataType}
	builder.builtinOpMap["remove"] = &removeAbsTD

	// size
	var sizeAbsDD AbsDD
	sizeAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
			TypeVar{Kind: stdLib.star},
		},
		Term: &sizeAbsDD,
	}
	sizeAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				&MapType{
					&sizeAbsTD.Vars[0],
					&sizeAbsTD.Vars[1],
				},
			},
		},
	}
	sizeAbsDD.Term = &Builtin{builder.typeMap["Bool"]}
	builder.builtinOpMap["size"] = &sizeAbsTD

	// builtin blt
	bltAbsDD := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["BNum"],
			},
			&DataVar{
				DataType: builder.typeMap["BNum"],
			},
		},
		Term: &Builtin{builder.typeMap["Bool"]},
	}
	builder.builtinOpMap["blt"] = &bltAbsDD

	// builtin badd
	var baddAbsDD AbsDD
	baddAbsTD := AbsTD{
		Vars: []TypeVar{
			TypeVar{Kind: stdLib.star},
		},
		Term: &baddAbsDD,
	}
	baddAbsDD = AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: &baddAbsTD.Vars[0],
			},
			&DataVar{
				DataType: builder.typeMap["BNum"],
			},
		},
		Term: &Builtin{builder.typeMap["BNum"]},
	}
	builder.builtinOpMap["badd"] = &baddAbsTD

	// builtin bsub
	bsubAbsDD := AbsDD{
		Vars: []*DataVar{
			&DataVar{
				DataType: builder.typeMap["BNum"],
			},
			&DataVar{
				DataType: builder.typeMap["BNum"],
			},
		},
		Term: &Builtin{builder.typeMap["Int256"]},
	}
	builder.builtinOpMap["bsub"] = &bsubAbsDD

}

func BuildCFG(n ast.AstNode) *CFGBuilder {
	builder := CFGBuilder{
		builtinOpMap:    map[string]Data{},
		typeMap:         map[string]Type{},
		constrTypeMap:   map[string]string{},
		intWidthTypeMap: map[int]*IntType{},
		natWidthTypeMap: map[int]*NatType{},
		varStack:        map[string][]Data{},
		fieldStack:      map[string][]Data{},
		constructor:     nil,
	}
	builder.initPrimitiveTypes()
	ast.Walk(&builder, n)

	return &builder
}
