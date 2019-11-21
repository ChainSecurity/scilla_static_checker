package ir

import (
	"errors"
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"strconv"
	"strings"
)

type CFGBuilder struct {
	builtinOpMap map[string]Data

	constructorType  map[string]string
	intWidthTypeMap  map[int]*IntType
	natWidthTypeMap  map[int]*NatType
	primitiveTypeMap map[string]Type
	definedADT       map[string]Type
	builtinADT       map[string]Type
	varStack         map[string][]Data
	constructor      *Proc
	Transitions      map[string]*Proc
	procedures       map[string]*Proc
	fieldTypeMap     map[string]Type

	mapTypeMap map[Type]map[Type]Type

	genericTypeConstructors map[string]*AbsTT
	genericDataConstructors map[string]*AbsTD
}

func (builder *CFGBuilder) constructGenericType(typeName string, varTypes []Type) Type {

	typeConstructor, ok := builder.genericTypeConstructors[typeName]
	if !ok {
		panic(errors.New(fmt.Sprintf("Unknown type constructor type: %s", typeName)))
	}

	if len(typeConstructor.Vars) != len(varTypes) {
		panic(errors.New(fmt.Sprintf("Wrong number of arguments for the constructor: %s", typeName)))
	}
	typ := &AppTT{
		Args: varTypes,
		To:   typeConstructor,
	}
	return typ
}

func (builder *CFGBuilder) getBuiltinOp(opName string, varTypes []Type) Data {

	if opName != "concat" {
		return builder.builtinOpMap[opName]
	}
	v0 := varTypes[0]
	switch raw0 := v0.(type) {
	case *RawType:
		v1 := varTypes[1]
		raw1, ok := v1.(*RawType)
		if !ok {
			panic(errors.New(fmt.Sprintf("Builtin concat wrong type: %T", v1)))
		}
		resType := "ByStr" + strconv.Itoa(raw0.Size+raw1.Size)
		op := &AbsDD{
			Vars: []DataVar{
				DataVar{
					DataType: v0,
				},
				DataVar{
					DataType: v1,
				},
			},
			Term: &Builtin{setDefaultType(builder.primitiveTypeMap, resType, &RawType{raw0.Size + raw1.Size})},
		}
		return op
	case *StrType:
		return builder.builtinOpMap["str_concat"]
	default:
		panic(errors.New(fmt.Sprintf("Builtin concat wrong type: %T", v0)))
	}
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

//Populates the Bind with the data. Return the lists of Name and Data of the varaibles that need to be in the new scope

func (builder *CFGBuilder) visitPattern(p ast.Pattern, t Type, bind *Bind) ([]string, []Data) {
	varNames := []string{}
	varBinds := []Data{}
	switch pat := p.(type) {
	case *ast.WildcardPattern:
		*bind = Bind{BindType: t}
	case *ast.BinderPattern:
		*bind = Bind{BindType: t}
		varNames = append(varNames, pat.Variable.Id)
		varBinds = append(varBinds, bind)

	case *ast.ConstructorPattern:
		fmt.Printf("ConstructorPattern %s\n", pat.ConstrName)
		var typeList []Type
		switch typ := t.(type) {
		case *EnumType:
			typeList = (*typ)[pat.ConstrName]
		case *AppTT:
			to := typ.To
			absTT, ok := to.(*AbsTT)
			if !ok {
				panic(errors.New(fmt.Sprintf("Not supported AppTT ConstructorPattern %T", to)))
			}
			args := typ.Args
			vars := absTT.Vars
			term := absTT.Term
			enumType, ok := term.(*EnumType)
			if !ok {
				panic(errors.New(fmt.Sprintf("Not supported AppTT ConstructorPattern %T", to)))
			}
			constrTypes := (*enumType)[pat.ConstrName]
			varToIndex := make(map[Type]int)
			for i, _ := range vars {
				varToIndex[&vars[i]] = i
			}
			indicies := make([]int, len(constrTypes))
			for i, c := range constrTypes {
				index, ok := varToIndex[c]
				if !ok {
					panic(errors.New(fmt.Sprintf("Not found Type order")))
				}
				indicies[i] = index
			}
			typeList = make([]Type, len(constrTypes))
			for i, _ := range indicies {
				typeList[i] = args[indicies[i]]
			}
		default:
			panic(errors.New(fmt.Sprintf("In %s; Not constructor pattern type: %T %d", pat.ConstrName, t, len(pat.Pats))))
		}
		if len(typeList) != len(pat.Pats) {
			panic(errors.New(fmt.Sprintf("constructor pattern argument length mistmatch: %s", pat.ConstrName)))
		}
		whenData := make([]Bind, len(pat.Pats))
		for i, subp := range pat.Pats {
			var bNames []string
			var bBinds []Data
			bNames, bBinds = builder.visitPattern(subp, typeList[i], &whenData[i])

			varNames = append(varNames, bNames...)
			varBinds = append(varBinds, bBinds...)
			//whenData[i] = b
		}
		*bind = Bind{
			BindType: t,
			Cond: &Cond{
				Case: pat.ConstrName,
				Data: whenData,
			},
		}
	default:
		panic(errors.New("Unknown pattern "))
	}
	return varNames, varBinds
}

func (builder *CFGBuilder) visitLiteral(l ast.Literal) Data {
	switch lit := l.(type) {
	case *ast.StringLiteral:
		typ := builder.primitiveTypeMap["String"]
		strTyp, ok := typ.(*StrType)

		if !ok {
			panic(errors.New(fmt.Sprintf("Type exception: String")))
		}

		str := Str{
			StrType: strTyp,
			Data:    lit.Val,
		}
		return &str
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
	//case *ast.MapLiteral:
	//ktyp := builder.visitASTType(lit.KeyType)
	//vtyp := builder.visitASTType(lit.ValType)
	//maptyp := MapType{ktyp, vtyp}
	//return &Map{
	//MapType: &maptyp,
	//Data:    map[string]string{},
	//}

	case *ast.ADTValLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	}
	return nil
}

func (builder *CFGBuilder) visitASTType(e ast.ASTType) Type {

	switch n := e.(type) {
	case *ast.PrimType:
		t, ok := builder.primitiveTypeMap[n.Name]
		if ok {
			return t
		}
		if strings.HasPrefix(n.Name, "ByStr") {
			width, err := strconv.Atoi(n.Name[5:])
			if err != nil {
				panic(err)
			}
			t = &RawType{width}
			builder.primitiveTypeMap[n.Name] = t
			return t
		}

		panic(errors.New(fmt.Sprintf("PrimType not found : %s", n.Name)))
	case *ast.ADT:
		typ, ok := builder.definedADT[n.Name]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown type: %s", n.Name)))
		}
		return typ

	case *ast.MapType:
		keyType := builder.visitASTType(n.KeyType)
		valType := builder.visitASTType(n.ValType)
		_, ok := builder.mapTypeMap[keyType]
		if !ok {
			mtype := &MapType{keyType, valType}
			builder.mapTypeMap[keyType] = map[Type]Type{}
			builder.mapTypeMap[keyType][valType] = mtype
			return mtype
		}

		_, ok = builder.mapTypeMap[keyType][valType]
		if !ok {
			mtype := &MapType{keyType, valType}
			builder.mapTypeMap[keyType][valType] = mtype
			return mtype
		}
		return builder.mapTypeMap[keyType][valType]

		//fmt.Printf("MapType %T %T\n", keyType, valType)
		//panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
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

func (builder *CFGBuilder) visitStatement(p *Proc, s ast.Statement) *Proc {
	var u Unit
	switch n := s.(type) {
	case *ast.LoadStatement:
		lhs := n.Lhs.Id
		rhs := n.Rhs.Id
		load := Load{
			Slot: rhs,
			Path: []Data{},
		}
		stackMapPush(builder.varStack, lhs, &load)
		u = &load
	case *ast.BindStatement:
		lhs := n.Lhs.Id
		rhs := builder.visitExpression(n.RhsExpr)
		stackMapPush(builder.varStack, lhs, rhs)
		switch r := rhs.(type) {
		case *AbsTD:
			u = r
		case *AbsDD:
			u = r
		case *AppTD:
			u = r
		case *AppDD:
			u = r
		case *Int:
			u = r
		case *Nat:
			u = r
		case *Raw:
			u = r
		case *Str:
			u = r
		case *Bnr:
			u = r
		case *Exc:
			u = r
		case *Msg:
			u = r
		case *Map:
			u = r
		default:
			u = nil
		}
	case *ast.StoreStatement:
		lhs := n.Lhs.Id
		rhs := n.Rhs.Id
		fmt.Println("StoreStatement", lhs, rhs)
		data, ok := stackMapPeek(builder.varStack, rhs)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", rhs)))
		}
		save := Save{
			Slot: lhs,
			Path: []Data{},
			Data: data,
		}
		u = &save
	case *ast.AcceptPaymentStatement:
		u = &Accept{}
	case *ast.SendMsgsStatement:
		d, ok := stackMapPeek(builder.varStack, n.Arg.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Arg.Id)))
		}
		u = &Send{Data: d}
	case *ast.MatchStatement:

		initialVarStack := stackMapCopy(builder.varStack)
		d, ok := stackMapPeek(builder.varStack, n.Arg.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Arg.Id)))
		}
		procCases := make([]ProcCase, len(n.Cases))
		contProc := Proc{
			Plan: []Unit{},
		}
		for i, mc := range n.Cases {
			//TODO create the DataVar and pass it as the arg for the call
			procCases[i] = ProcCase{
				Body: Proc{
					Plan: []Unit{},
				},
			}
			copyVarStack := stackMapCopy(initialVarStack)
			builder.varStack = copyVarStack
			varNames, varBinds := builder.visitPattern(mc.Pat, builder.TypeOf(d), &procCases[i].Bind)
			for j, name := range varNames {
				stackMapPush(builder.varStack, name, varBinds[j])
			}

			curProc := &procCases[i].Body
			for _, s := range mc.Body {
				nextProc := builder.visitStatement(curProc, s)
				if nextProc != nil {
					curProc = nextProc
				}
			}
			curProc.Jump = &CallProc{
				Args: []Data{},
				To:   &contProc,
			}

			for _, name := range varNames {
				stackMapPop(builder.varStack, name)
			}

		}

		pp := PickProc{
			From: d,
			With: procCases,
		}
		p.Jump = &pp
		builder.restoreVarStack(initialVarStack)
		return &contProc
	case *ast.MapGetStatement:
		if !n.IsValRetrieve {
			panic(errors.New(fmt.Sprintf("Only value retrive is implemeted for MapGetStatement")))
		}

		slot := n.Name.Id
		path := make([]Data, len(n.Keys))
		var ok bool
		for i, _ := range n.Keys {
			path[i], ok = stackMapPeek(builder.varStack, n.Keys[i].Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", n.Keys[i].Id)))
			}
		}
		load := Load{slot, path}
		u = &load
		stackMapPush(builder.varStack, n.Lhs.Id, &load)
	case *ast.MapUpdateStatement:
		if n.Rhs == nil {
			panic(errors.New(fmt.Sprintf("Unhandled Map delete: %T", n)))
		}
		slot := n.Name.Id
		data, ok := stackMapPeek(builder.varStack, n.Rhs.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Rhs.Id)))
		}
		path := make([]Data, len(n.Keys))
		for i, _ := range n.Keys {
			path[i], ok = stackMapPeek(builder.varStack, n.Keys[i].Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", n.Keys[i].Id)))
			}
		}
		u = &Save{slot, path, data}

	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
	if u != nil {
		p.Plan = append(p.Plan, u)
	}
	return nil
}

func (builder *CFGBuilder) visitExpression(e ast.Expression) Data {
	switch n := e.(type) {
	case *ast.ConstrExpression:
		constrName := n.ConstructorName

		typeName, ok := builder.constructorType[constrName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown constructor: %s", constrName)))
		}

		ds := make([]Data, len(n.Args))
		for i, arg := range n.Args {
			d, ok := stackMapPeek(builder.varStack, arg.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", arg.Id)))
			}
			ds[i] = d
		}

		// Handling user defined ADTs
		typ, ok := builder.definedADT[typeName]
		if ok {
			res := Enum{
				EnumType: typ,
				Case:     constrName,
				Data:     ds,
			}
			return &res
		}

		// Handling builtin ADTs, can be generic type
		ts := make([]Type, len(n.Types))

		for i, arg := range n.Types {
			ts[i] = builder.visitASTType(arg)
		}
		typ = builder.constructGenericType(typeName, ts)

		constr, ok := builder.genericDataConstructors[constrName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown constructor: %s", constrName)))
		}
		appTD := AppTD{
			Args: ts,
			To:   constr,
		}
		appDD := AppDD{
			Args: ds,
			To:   &appTD,
		}
		fmt.Printf("Constructor %s\n\t%s\n\t%s\n\t%T\n\t%T\n", constrName, ts, ds, typ, constr)
		return &appDD
	case *ast.FunExpression:
		lhs := n.Lhs.Id
		lhsTyp := builder.visitASTType(n.LhsType)
		dataVar := DataVar{lhsTyp}
		stackMapPush(builder.varStack, lhs, &dataVar)
		defer stackMapPop(builder.varStack, lhs)

		rhs := builder.visitExpression(n.RhsExpr)

		return &AbsDD{Vars: []DataVar{dataVar}, Term: rhs}
	case *ast.MatchExpression:
		lhs := n.Lhs.Id

		data, ok := stackMapPeek(builder.varStack, lhs)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", lhs)))
		}
		pd := PickData{
			From: data,
			With: make([]DataCase, len(n.Cases)),
		}
		for i, c := range n.Cases {
			mec := &pd.With[i].Bind
			varNames, varBinds := builder.visitPattern(c.Pat, builder.TypeOf(data), mec)
			for j, name := range varNames {
				stackMapPush(builder.varStack, name, varBinds[j])
			}

			pd.With[i].Body = builder.visitExpression(c.Expr)

			for _, name := range varNames {
				stackMapPop(builder.varStack, name)
			}
		}
		return &pd
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
		varTypes := make([]Type, len(n.Args))
		for i, a := range n.Args {
			v, ok := stackMapPeek(builder.varStack, a.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", a.Id)))
			}
			vars[i] = v
			varTypes[i] = builder.TypeOf(v)
		}

		op := builder.getBuiltinOp(opName, varTypes)
		switch op := op.(type) {
		case *AbsTD:
			if len(op.Vars) != len(vars) {
				panic(errors.New(fmt.Sprintf("Wrong number of Builtin AbsTD args")))
			}
			types := make([]Type, len(vars))
			for i, _ := range vars {
				types[i] = builder.TypeOf(vars[i])
			}
			appTD := &AppTD{
				Args: types,
				To:   op,
			}
			appDD := AppDD{
				Args: vars,
				To:   appTD,
			}
			return &appDD
		case *AbsDD:
			if len(op.Vars) != len(vars) {
				panic(errors.New(fmt.Sprintf("Wrong number of Builtin AbsDD args")))
			}
			appDD := AppDD{
				Args: vars,
				To:   op,
			}
			return &appDD
		default:
			panic(errors.New(fmt.Sprintf("Unhandled Builtin op type: %T\n", n)))
		}
	case *ast.LetExpression:
		varName := n.Var.Id
		expr := builder.visitExpression(n.Expr)
		stackMapPush(builder.varStack, varName, expr)
		defer stackMapPop(builder.varStack, varName)
		body := builder.visitExpression(n.Body)
		return body
	case *ast.AppExpression:
		rhsData := make([]Data, len(n.RhsList))
		for i, rhs := range n.RhsList {
			data, ok := stackMapPeek(builder.varStack, rhs.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", rhs.Id)))
			}
			rhsData[i] = data
		}
		lhs, ok := stackMapPeek(builder.varStack, n.Lhs.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Lhs.Id)))
		}

		i := 0
		curr := lhs
		accum := lhs
		for i < len(rhsData) {
			absDD, ok := curr.(*AbsDD)
			if !ok {
				panic(errors.New(fmt.Sprintf("AppExpression absDD wrong type: %T", curr)))
			}
			currData := rhsData[i : i+len(absDD.Vars)]
			i = i + len(absDD.Vars)
			accum = &AppDD{
				Args: currData,
				To:   accum,
			}
			curr = absDD.Term
		}

		return accum
	case *ast.MessageExpression:
		//nameList := make([]string, len(n.MArgs))
		//dataList := make([]Data, len(n.MArgs))
		data := make(map[string]Data)

		for _, marg := range n.MArgs {
			name := marg.Var
			d := builder.visitPayload(marg.Pl)
			data[name] = d
		}

		typ := builder.primitiveTypeMap["Message"]
		msgType, ok := typ.(*MsgType)

		if !ok {
			panic(errors.New(fmt.Sprintf("Type exception: Message")))
		}

		return &Msg{
			MsgType: msgType,
			Data:    data,
		}
	default:
		fmt.Printf("Unhandled Expression type: %T\n", n)
		//panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
		return nil
	}
}

func (builder *CFGBuilder) visitPayload(pl ast.Payload) Data {
	switch n := pl.(type) {
	case *ast.PayloadLiteral:
		return builder.visitLiteral(n.Lit)
	case *ast.PayloadVariable:
		data, ok := stackMapPeek(builder.varStack, n.Val.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Val.Id)))
		}
		return data
	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
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
			builder.constructorType[constrName] = typeName
		}
		builder.definedADT[typeName] = &typ
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
	t := builder.visitASTType(f.Type)
	builder.fieldTypeMap[name] = t
	data := builder.visitExpression(f.Expr)
	stackMapPush(builder.varStack, name, data)
	return name, data
}

func (builder *CFGBuilder) restoreVarStack(varStack map[string][]Data) {
	builder.varStack = varStack
}

func (builder *CFGBuilder) visitComponent(comp *ast.Component) {
	varStackCopy := stackMapCopy(builder.varStack)
	defer builder.restoreVarStack(varStackCopy)

	paramNames := make([]string, len(comp.Params))
	params := map[string]Type{}
	paramVars := make([]DataVar, len(comp.Params))

	for i, p := range comp.Params {
		pName, pType := builder.visitParams(p)
		params[pName] = pType
		paramNames[i] = pName
		paramVars[i] = DataVar{pType}
		stackMapPush(builder.varStack, pName, &paramVars[i])
		defer stackMapPop(builder.varStack, pName)
	}

	implicitVars := []DataVar{
		DataVar{builder.primitiveTypeMap["Uint128"]},
		DataVar{builder.primitiveTypeMap["ByStr20"]},
	}
	stackMapPush(builder.varStack, "_amount", &implicitVars[0])
	stackMapPush(builder.varStack, "_sender", &implicitVars[1])
	defer stackMapPop(builder.varStack, "_amount")
	defer stackMapPop(builder.varStack, "_sender")
	dataVars := append(implicitVars, paramVars...)

	firstProc := Proc{
		Vars: dataVars,
		Plan: []Unit{},
	}
	proc := &firstProc
	for _, s := range comp.Body {
		contProc := builder.visitStatement(proc, s)
		if contProc != nil {
			proc = contProc
		}

	}

	fmt.Printf("Component %s type: %s\n\tvars: %s\n\tplan: %s\n", comp.Name.Id, comp.ComponentType, dataVars, proc.Plan)

	if comp.ComponentType == "procedure" {
		builder.Transitions[comp.Name.Id] = &firstProc
	} else if comp.ComponentType == "transition" {
		builder.Transitions[comp.Name.Id] = &firstProc
	} else {
		panic(errors.New(fmt.Sprintf("Wrong Component type: %s", comp.ComponentType)))
	}
}

func (builder *CFGBuilder) Visit(node ast.AstNode) ast.Visitor {
	switch n := node.(type) {
	case ast.LibEntry:
		builder.visitLibEntry(n)
	case *ast.Contract:
		//name := n.Name.Id
		paramNames := make([]string, len(n.Params))
		params := map[string]Type{}
		dataVars := make([]DataVar, len(n.Params))
		for i, p := range n.Params {
			pName, pType := builder.visitParams(p)
			params[pName] = pType
			paramNames[i] = pName
			dataVars[i] = DataVar{pType}
			stackMapPush(builder.varStack, pName, &dataVars[i])
		}

		builder.constructor = &Proc{
			Vars: dataVars,
			Plan: make([]Unit, len(n.Fields)),
		}
		for i, f := range n.Fields {
			n, d := builder.visitField(f)
			stackMapPush(builder.varStack, n, d)
			builder.constructor.Plan[i] = &Save{n, []Data{}, d}
		}

		//builder.Transitions = make[[]*Proc, len(n.Components)]
		//for _, c := range n.Components {
		//builder.Visit(c)
		//builder.Transitions = append(builder.Transitions, &e)
		//}

		//for _, pName := range paramNames {
		//stackMapPop(builder.varStack, pName)
		//}

		//for

		//t := n.Type
		//expr := builder.visitExpression(n.Expr)
		//builder.Field[name] = Save{}
		//builder.visitLibEntry(n)
	case *ast.Component:
		builder.visitComponent(n)
	case *ast.Identifier:
		//do nothing
	case *ast.Location:
		//do nothing
	default:
		//fmt.Printf("Unhandled type: %T\n", n)
		// do nothing
	}
	return builder
}

func (builder *CFGBuilder) initPrimitiveTypes() {
	stdLib := StdLib()

	builder.constructorType["Nil"] = "List"
	builder.constructorType["Cons"] = "List"

	builder.genericTypeConstructors["List"] = stdLib.List
	builder.genericDataConstructors["Nil"] = stdLib.Nil
	builder.genericDataConstructors["Cons"] = stdLib.Cons

	builder.genericTypeConstructors["Option"] = stdLib.Option
	builder.genericDataConstructors["Some"] = stdLib.Some
	builder.genericDataConstructors["None"] = stdLib.None
	//builder.constructorType["True"] = "Boolean"
	//builder.constructorType["False"] = "Boolean"
	//builder.primitiveTypeMap["Boolean"] = stdLib.Boolean

	builder.constructorType["True"] = "Bool"
	builder.constructorType["False"] = "Bool"
	builder.definedADT["Bool"] = stdLib.Boolean

	sizes := []int{32, 64, 128, 256}
	for _, s := range sizes {
		intName := "Int" + strconv.Itoa(s)
		intTyp := IntType{s}
		uintName := "Uint" + strconv.Itoa(s)
		uintTyp := NatType{s}
		builder.intWidthTypeMap[s] = &intTyp
		builder.natWidthTypeMap[s] = &uintTyp
		builder.primitiveTypeMap[intName] = &intTyp
		builder.primitiveTypeMap[uintName] = &uintTyp

		// Conversion functions
		var intAbsDD AbsDD
		intAbsTD := AbsTD{
			Vars: []TypeVar{
				TypeVar{Kind: stdLib.star},
			},
			Term: &intAbsDD,
		}
		intAbsDD = AbsDD{
			Vars: []DataVar{
				DataVar{DataType: &intAbsTD.Vars[0]},
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
			Vars: []DataVar{
				DataVar{DataType: &uintAbsTD.Vars[0]},
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

	builder.primitiveTypeMap["Message"] = &MsgType{}
	builder.primitiveTypeMap["String"] = &StrType{}
	builder.primitiveTypeMap["ByStr"] = &RawType{-1}
	builder.primitiveTypeMap["ByStr32"] = &RawType{32}
	builder.primitiveTypeMap["ByStr33"] = &RawType{33}
	builder.primitiveTypeMap["ByStr64"] = &RawType{64}
	builder.primitiveTypeMap["ByStr20"] = &RawType{20}
	builder.primitiveTypeMap["BNum"] = &BnrType{}

	builder.primitiveTypeMap["Bool"] = stdLib.Boolean

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
			Vars: []DataVar{
				DataVar{DataType: &bAbsTD.Vars[0]},
				DataVar{DataType: &bAbsTD.Vars[0]},
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
			Vars: []DataVar{
				DataVar{DataType: &bAbsTD.Vars[0]},
				DataVar{DataType: &bAbsTD.Vars[0]},
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
		Vars: []DataVar{
			DataVar{DataType: &powAbsTD.Vars[0]},
			DataVar{DataType: &powAbsTD.Vars[1]},
		},
		Term: &Builtin{&powAbsTD.Vars[0]},
	}
	builder.builtinOpMap["pow"] = &powAbsTD

	// builtin concat
	strConcatOp := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["String"]},
	}
	builder.builtinOpMap["str_concat"] = &strConcatOp

	//builtin substr
	strSubstrOp := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				DataType: builder.natWidthTypeMap[32],
			},
			DataVar{
				DataType: builder.natWidthTypeMap[32],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["String"]},
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
		Vars: []DataVar{
			DataVar{DataType: &toStrAbsTD.Vars[0]},
		},
		Term: &Builtin{builder.primitiveTypeMap["String"]},
	}
	builder.builtinOpMap["to_string"] = &toStrAbsTD

	//builtin strlen

	strlenOp := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
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
		Vars: []DataVar{
			DataVar{DataType: &shaAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.primitiveTypeMap["ByStr32"],
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
		Vars: []DataVar{
			DataVar{DataType: &keccakAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.primitiveTypeMap["ByStr32"],
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
		Vars: []DataVar{
			DataVar{DataType: &ripemdAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.primitiveTypeMap["ByStr20"],
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
		Vars: []DataVar{
			DataVar{DataType: &toBystrAbsTD.Vars[0]},
		},
		Term: &Builtin{
			BuiltinType: builder.primitiveTypeMap["ByStr"],
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
		Vars: []DataVar{
			DataVar{DataType: &touint256AbsTD.Vars[0]},
		},
		Term: &Builtin{builder.natWidthTypeMap[256]},
	}
	builder.builtinOpMap["to_uint256"] = &touint256AbsTD

	// schnorr_verify
	schnorrVerifyAbsDD := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr33"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr64"],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["Bool"]},
	}
	builder.builtinOpMap["schnorr_verify"] = &schnorrVerifyAbsDD

	// ecdsa_verify
	ecdsaVerifyAbsDD := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr33"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr64"],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["Bool"]},
	}
	builder.builtinOpMap["ecdsa_verify"] = &ecdsaVerifyAbsDD

	// bech32_to_bystr20
	bech32ToBystr20AbsDD := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
			},
		},
		Term: &Builtin{
			&AppTT{
				Args: []Type{builder.primitiveTypeMap["ByStr20"]},
				To:   stdLib.Option,
			},
		},
	}
	builder.builtinOpMap["bech32_to_bystr20"] = &bech32ToBystr20AbsDD

	// bystr20_to_bech32
	bystr20ToBech32AbsDD := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["ByStr20"],
			},
		},
		Term: &Builtin{
			&AppTT{
				Args: []Type{builder.primitiveTypeMap["String"]},
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
		Vars: []DataVar{
			DataVar{
				&MapType{
					&putAbsTD.Vars[0],
					&putAbsTD.Vars[1],
				},
			},
			DataVar{
				&putAbsTD.Vars[0],
			},
			DataVar{
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
		Vars: []DataVar{
			DataVar{
				&MapType{
					&getAbsTD.Vars[0],
					&getAbsTD.Vars[1],
				},
			},
			DataVar{
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
		Vars: []DataVar{
			DataVar{
				&MapType{
					&containsAbsTD.Vars[0],
					&containsAbsTD.Vars[1],
				},
			},
			DataVar{
				&containsAbsTD.Vars[0],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["Bool"]},
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
		Vars: []DataVar{
			DataVar{
				&MapType{
					&removeAbsTD.Vars[0],
					&removeAbsTD.Vars[1],
				},
			},
			DataVar{
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
		Vars: []DataVar{
			DataVar{
				&MapType{
					&sizeAbsTD.Vars[0],
					&sizeAbsTD.Vars[1],
				},
			},
		},
	}
	sizeAbsDD.Term = &Builtin{builder.primitiveTypeMap["Bool"]}
	builder.builtinOpMap["size"] = &sizeAbsTD

	// builtin blt
	bltAbsDD := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["BNum"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["BNum"],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["Bool"]},
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
		Vars: []DataVar{
			DataVar{
				DataType: &baddAbsTD.Vars[0],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["BNum"],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["BNum"]},
	}
	builder.builtinOpMap["badd"] = &baddAbsTD

	// builtin bsub
	bsubAbsDD := AbsDD{
		Vars: []DataVar{
			DataVar{
				DataType: builder.primitiveTypeMap["BNum"],
			},
			DataVar{
				DataType: builder.primitiveTypeMap["BNum"],
			},
		},
		Term: &Builtin{builder.primitiveTypeMap["Int256"]},
	}
	builder.builtinOpMap["bsub"] = &bsubAbsDD

}

func BuildCFG(n ast.AstNode) *CFGBuilder {
	builder := CFGBuilder{
		builtinOpMap:     map[string]Data{},
		primitiveTypeMap: map[string]Type{},
		definedADT:       map[string]Type{},
		builtinADT:       map[string]Type{},
		constructorType:  map[string]string{},
		intWidthTypeMap:  map[int]*IntType{},
		natWidthTypeMap:  map[int]*NatType{},
		varStack:         map[string][]Data{},
		constructor:      nil,
		Transitions:      map[string]*Proc{},
		procedures:       map[string]*Proc{},
		fieldTypeMap:     map[string]Type{},

		mapTypeMap:              map[Type]map[Type]Type{},
		genericTypeConstructors: map[string]*AbsTT{},
		genericDataConstructors: map[string]*AbsTD{},
	}
	builder.initPrimitiveTypes()
	ast.Walk(&builder, n)

	return &builder
}
