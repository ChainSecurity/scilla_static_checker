package ir

import (
	"errors"
	"fmt"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ast"
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
	typeVarStack     map[string][]Type
	Constructor      *Proc
	Transitions      map[string]*Proc
	Procedures       map[string]*Proc
	fieldTypeMap     map[string]Type

	mapTypeMap map[Type]map[Type]*MapType

	genericTypeConstructors map[string]*AbsTT
	genericDataConstructors map[string]*AbsTD
	nodeCounter             int64
	currentProc             string
	kind                    Kind
}

func (b *CFGBuilder) newIDNode() IDNode {
	id := b.nodeCounter
	b.nodeCounter++
	return IDNode{id}
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
		IDNode: builder.newIDNode(),
		Args:   varTypes,
		To:     typeConstructor,
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

			IDNode: builder.newIDNode(),
			Vars: []DataVar{
				DataVar{
					DataType: v0,
				},
				DataVar{
					DataType: v1,
				},
			},
			Term: &Builtin{
				builder.newIDNode(),
				setDefaultType(builder.primitiveTypeMap, resType, &RawType{builder.newIDNode(), raw0.Size + raw1.Size}),
				"raw_concat",
			},
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
		*bind = Bind{
			IDNode:   builder.newIDNode(),
			BindType: t,
		}
	case *ast.BinderPattern:
		*bind = Bind{
			IDNode:   builder.newIDNode(),
			BindType: t,
		}
		varNames = append(varNames, pat.Variable.Id)
		varBinds = append(varBinds, bind)

	case *ast.ConstructorPattern:
		var typeList []Type
		switch typ := t.(type) {
		case *EnumType:
			typeList = typ.Constructors[pat.ConstrName]
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
			constrTypes := enumType.Constructors[pat.ConstrName]
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
			IDNode:   builder.newIDNode(),
			BindType: t,
			Cond: &Cond{
				IDNode: builder.newIDNode(),
				Case:   pat.ConstrName,
				Data:   whenData,
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
			IDNode:  builder.newIDNode(),
			StrType: strTyp,
			Data:    lit.Val,
		}
		return &str
	case *ast.BNumLiteral:
		typ, ok := builder.primitiveTypeMap["BNum"]

		if !ok {
			panic(errors.New(fmt.Sprintf("Type not found: BNum")))
		}

		bnrTyp, ok := typ.(*BnrType)

		if !ok {
			panic(errors.New(fmt.Sprintf("Type exception: BnrType")))
		}

		bnr := Bnr{
			IDNode:  builder.newIDNode(),
			BnrType: bnrTyp,
			Data:    lit.Val,
		}
		return &bnr
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
			IDNode:  builder.newIDNode(),
			IntType: typ, Data: lit.Val,
		}
	case *ast.UintLiteral:
		typ, ok := builder.natWidthTypeMap[lit.Width]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown width: %d", lit.Width)))
		}
		return &Nat{
			IDNode:  builder.newIDNode(),
			NatType: typ, Data: lit.Val,
		}
	case *ast.MapLiteral:
		ktyp := builder.visitASTType(lit.KeyType)
		vtyp := builder.visitASTType(lit.ValType)
		maptyp := builder.setdefaultMap(ktyp, vtyp)
		return &Map{
			IDNode:  builder.newIDNode(),
			MapType: maptyp,
			Data:    map[string]string{},
		}

	case *ast.ADTValLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	default:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	}
}

func (builder *CFGBuilder) setdefaultMap(keyType Type, valType Type) *MapType {
	_, ok := builder.mapTypeMap[keyType]
	if !ok {
		builder.mapTypeMap[keyType] = map[Type]*MapType{}
	}
	_, ok = builder.mapTypeMap[keyType][valType]
	if !ok {
		mtype := &MapType{
			IDNode:  builder.newIDNode(),
			KeyType: keyType,
			ValType: valType,
		}
		builder.mapTypeMap[keyType][valType] = mtype
	}
	return builder.mapTypeMap[keyType][valType]
}

func (builder *CFGBuilder) visitASTType(e ast.ASTType) Type {
	fmt.Printf("visitASTType %T\n", e)

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
			t = &RawType{
				builder.newIDNode(),
				width,
			}
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
		return builder.setdefaultMap(keyType, valType)

		//fmt.Printf("MapType %T %T\n", keyType, valType)
		//panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	case *ast.FunType:
		argType := builder.visitASTType(n.ArgType)
		valType := builder.visitASTType(n.ValType)
		argDataVar := DataVar{
			IDNode:   builder.newIDNode(),
			DataType: argType,
		}
		fmt.Printf("%T %T\n", argDataVar.DataType, valType)
		return &AllDD{
			IDNode: builder.newIDNode(),
			Vars:   []DataVar{argDataVar},
			Term:   valType,
		}
		//panic(errors.New(fmt.Sprintf("CFGBuilder visitASTType unhandled type: %T", n)))
	case *ast.TypeVar:
		fmt.Printf("TypeVar %s\n", n.Name)
		t, ok := typeStackMapPeek(builder.typeVarStack, n.Name)
		if !ok {
			panic(errors.New(fmt.Sprintf("TypeVar not found: %s", n.Name)))
		}
		return t
	case *ast.PolyFun:
		fmt.Printf("PolyFun %s\n", n.TypeVal)
		tv := &TypeVar{
			IDNode: builder.newIDNode(),
			Kind:   builder.kind,
		}
		typeStackMapPush(builder.typeVarStack, n.TypeVal, tv)
		t := builder.visitASTType(n.Body)
		return t
	//case *ast.Unit:
	//fmt.Printf("%T", n)
	default:
		panic(errors.New(fmt.Sprintf("CFGBuilder visitASTType unhandled type: %T", n)))
	}
}

func (builder *CFGBuilder) visitStatement(p *Proc, s ast.Statement) *Proc {
	var u Unit
	switch n := s.(type) {
	case *ast.LoadStatement:
		lhs := n.Lhs.Id
		rhs := n.Rhs.Id
		load := Load{
			IDNode: builder.newIDNode(),
			Slot:   rhs,
			Path:   []Data{},
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
			IDNode: builder.newIDNode(),
			Slot:   lhs,
			Path:   []Data{},
			Data:   data,
		}
		u = &save
	case *ast.AcceptPaymentStatement:
		u = &Accept{
			IDNode: builder.newIDNode(),
		}
	case *ast.SendMsgsStatement:
		d, ok := stackMapPeek(builder.varStack, n.Arg.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Arg.Id)))
		}
		u = &Send{
			IDNode: builder.newIDNode(),
			Data:   d,
		}
	case *ast.MatchStatement:

		initialVarStack := stackMapCopy(builder.varStack)
		d, ok := stackMapPeek(builder.varStack, n.Arg.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Arg.Id)))
		}
		procCases := make([]ProcCase, len(n.Cases))
		contProc := Proc{
			IDNode:   builder.newIDNode(),
			ProcName: builder.currentProc,
			Plan:     []Unit{},
		}
		for i, mc := range n.Cases {
			//TODO create the DataVar and pass it as the arg for the call
			procCases[i] = ProcCase{
				IDNode: builder.newIDNode(),
				Body: Proc{
					IDNode:   builder.newIDNode(),
					ProcName: builder.currentProc,
					Plan:     []Unit{},
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
				IDNode: builder.newIDNode(),
				Args:   []Data{},
				To:     &contProc,
			}

			for _, name := range varNames {
				stackMapPop(builder.varStack, name)
			}

		}

		pp := PickProc{
			IDNode: builder.newIDNode(),
			From:   d,
			With:   procCases,
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
		load := Load{
			IDNode: builder.newIDNode(),
			Slot:   slot,
			Path:   path,
		}
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
		u = &Save{
			IDNode: builder.newIDNode(),
			Slot:   slot,
			Path:   path,
			Data:   data,
		}
	case *ast.ReadFromBCStatement:
		switch n.RhsStr {
		case "BLOCKNUMBER":
			bn := BuiltinVar{
				IDNode:         builder.newIDNode(),
				BuiltinVarType: builder.primitiveTypeMap["BNum"],
				Label:          "BLOCKNUMBER",
			}
			stackMapPush(builder.varStack, n.Lhs.Id, &bn)
		default:
			panic(errors.New(fmt.Sprintf("Unhandled ReadFromBCStatement: %s", n.RhsStr)))
		}

	case *ast.CreateEvntStatement:
		d, ok := stackMapPeek(builder.varStack, n.Arg.Id)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Arg.Id)))
		}
		u = &Event{
			IDNode: builder.newIDNode(),
			Data:   d,
		}
	case *ast.CallProcStatement:
		contProc := Proc{
			IDNode:   builder.newIDNode(),
			ProcName: builder.currentProc,
			Plan:     []Unit{},
		}

		procedureProc, ok := builder.Procedures[n.Arg.Id]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown procedure: %s", n.Arg.Id)))
		}

		cp := CallProc{
			IDNode: builder.newIDNode(),
			Args:   make([]Data, len(n.Messages)+3),
			To:     procedureProc,
		}
		cp.Args[0] = &p.Vars[0]
		cp.Args[1] = &p.Vars[1]
		cp.Args[len(cp.Args)-1] = &contProc
		for i, m := range n.Messages {
			d, ok := stackMapPeek(builder.varStack, m.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", m.Id)))
			}
			cp.Args[2+i] = d
		}
		p.Jump = &cp

		return &contProc
	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
	if u != nil {
		p.Plan = append(p.Plan, u)
	}
	return nil
}

func (builder *CFGBuilder) visitExpression(e ast.Expression) Data {
	fmt.Printf("visitExpression %T\n", e)
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
				IDNode:   builder.newIDNode(),
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
			IDNode: builder.newIDNode(),
			Args:   ts,
			To:     constr,
		}
		appDD := AppDD{
			IDNode: builder.newIDNode(),
			Args:   ds,
			To:     &appTD,
		}
		fmt.Printf("Constructor %s\n\t%s\n\t%s\n\t%T\n\t%T\n", constrName, ts, ds, typ, constr)
		return &appDD
	case *ast.FunExpression:
		fmt.Println("FunExpression", n.AnnotatedNode.Loc, n.Lhs.Id)
		lhs := n.Lhs.Id
		lhsTyp := builder.visitASTType(n.LhsType)
		absDD := &AbsDD{
			IDNode: builder.newIDNode(),
			Vars:   make([]DataVar, 1),
			Term:   nil,
		}
		absDD.Vars[0] = DataVar{
			IDNode:   builder.newIDNode(),
			DataType: lhsTyp,
		}
		stackMapPush(builder.varStack, lhs, &absDD.Vars[0])

		defer stackMapPop(builder.varStack, lhs)
		rhs := builder.visitExpression(n.RhsExpr)
		absDD.Term = rhs
		fmt.Printf("FunExpression return %T\n", absDD)
		return absDD
	case *ast.MatchExpression:
		lhs := n.Lhs.Id

		data, ok := stackMapPeek(builder.varStack, lhs)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", lhs)))
		}
		pd := PickData{
			IDNode: builder.newIDNode(),
			From:   data,
			With:   make([]DataCase, len(n.Cases)),
		}
		for i, c := range n.Cases {
			pd.With[i].IDNode = builder.newIDNode()
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
			fmt.Println(op.Vars, vars, opName)
			fmt.Println(len(op.Vars), len(vars))

			types := make([]Type, len(op.Vars))
			for i, _ := range op.Vars {
				types[i] = builder.TypeOf(vars[i])
			}
			appTD := &AppTD{
				IDNode: builder.newIDNode(),
				Args:   types,
				To:     op,
			}
			appDD := AppDD{
				IDNode: builder.newIDNode(),
				Args:   vars,
				To:     appTD,
			}
			return &appDD
		case *AbsDD:
			if len(op.Vars) != len(vars) {
				panic(errors.New(fmt.Sprintf("Wrong number of Builtin AbsDD args")))
			}
			appDD := AppDD{
				IDNode: builder.newIDNode(),
				Args:   vars,
				To:     op,
			}
			return &appDD
		default:
			panic(errors.New(fmt.Sprintf("Unhandled Builtin op type: %T\n", n)))
		}
	case *ast.LetExpression:
		fmt.Println("LetExpression", n.AnnotatedNode.Loc)
		varName := n.Var.Id
		expr := builder.visitExpression(n.Expr)
		stackMapPush(builder.varStack, varName, expr)
		defer stackMapPop(builder.varStack, varName)
		body := builder.visitExpression(n.Body)
		return body
	case *ast.AppExpression:
		fmt.Println("AppExpression", n.AnnotatedNode.Loc, n.Lhs.Id)
		rhsData := make([]Data, len(n.RhsList))
		for i, rhs := range n.RhsList {
			fmt.Println("\t", rhs.Id)
			data, ok := stackMapPeek(builder.varStack, rhs.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", rhs.Id)))
			}
			rhsData[i] = data
		}
		lhs, ok := stackMapPeek(builder.varStack, n.Lhs.Id)
		fmt.Printf("Lhs type %T\n", lhs)
		if !ok {
			panic(errors.New(fmt.Sprintf("variable not found: %s", n.Lhs.Id)))
		}

		i := 0
		curr := lhs
		accum := lhs
		var nextCurr Data
		for i < len(rhsData) {
			fmt.Printf("%d %T %T\n", i, curr, rhsData[i])
			argsCount := 0
			switch c := curr.(type) {
			case *AbsDD:
				argsCount = len(c.Vars)
				nextCurr = c.Term
			case *DataVar:
				t := c.DataType
				allDD, ok := t.(*AllDD)
				if !ok {
					panic(errors.New(fmt.Sprintf("AppExpression DataVar wrong type: %T", t)))
				}
				argsCount = len(allDD.Vars)
				nextCurr = nil
			default:
				panic(errors.New(fmt.Sprintf("AppExpression absDD wrong type: %T", curr)))
			}
			currData := rhsData[i : i+argsCount]
			i = i + argsCount
			accum = &AppDD{
				IDNode: builder.newIDNode(),
				Args:   currData,
				To:     accum,
			}
			curr = nextCurr
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
			IDNode:  builder.newIDNode(),
			MsgType: msgType,
			Data:    data,
		}
	case *ast.TFunExpression:
		fmt.Println("TFunExpression", n.AnnotatedNode.Loc)
		lhs := n.Lhs.Id
		dataVar := DataVar{
			IDNode:   builder.newIDNode(),
			DataType: nil,
		}
		stackMapPush(builder.varStack, lhs, &dataVar)
		defer stackMapPop(builder.varStack, lhs)

		rhs := builder.visitExpression(n.RhsExpr)

		return &AbsDD{
			IDNode: builder.newIDNode(),
			Vars:   []DataVar{dataVar},
			Term:   rhs,
		}
	default:
		//fmt.Printf("Unhandled Expression type: %T\n", n)
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
		//return nil
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
		fmt.Println("LibraryVariable", n.Name.Loc, n.Name.Id)
		name := n.Name.Id
		//typ := builder.visitASTType(n.VarType)
		//fmt.Printf("%T\n", typ)
		v := builder.visitExpression(n.Expr)
		fmt.Printf("LibraryVariable %s %T\n", name, v)
		stackMapPush(builder.varStack, name, v)
	case *ast.LibraryType:
		typeName := n.Name.Id
		enumType := EnumType{
			IDNode:       builder.newIDNode(),
			Constructors: make(map[string][]Type),
		}
		for _, ctr := range n.CtrDefs {
			constrName, types := builder.visitCtr(ctr)
			enumType.Constructors[constrName] = types
			builder.constructorType[constrName] = typeName
		}
		builder.definedADT[typeName] = &enumType
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
		paramVars[i] = DataVar{
			IDNode:   builder.newIDNode(),
			DataType: pType,
		}
		stackMapPush(builder.varStack, pName, &paramVars[i])
		defer stackMapPop(builder.varStack, pName)
	}

	implicitVars := []DataVar{
		DataVar{
			IDNode:   builder.newIDNode(),
			DataType: builder.primitiveTypeMap["Uint128"]},
		DataVar{
			IDNode:   builder.newIDNode(),
			DataType: builder.primitiveTypeMap["ByStr20"],
		},
	}
	stackMapPush(builder.varStack, "_amount", &implicitVars[0])
	stackMapPush(builder.varStack, "_sender", &implicitVars[1])
	defer stackMapPop(builder.varStack, "_amount")
	defer stackMapPop(builder.varStack, "_sender")
	dataVars := append(implicitVars, paramVars...)

	firstProc := Proc{
		IDNode:   builder.newIDNode(),
		ProcName: comp.Name.Id,
		Vars:     dataVars,
		Plan:     []Unit{},
	}
	builder.currentProc = comp.Name.Id
	proc := &firstProc
	for _, s := range comp.Body {
		contProc := builder.visitStatement(proc, s)
		if contProc != nil {
			proc = contProc
		}

	}

	//fmt.Printf("Component %s type: %s\n\tvars: %s\n\tplan: %s\n", comp.Name.Id, comp.ComponentType, dataVars, proc.Plan)
	if comp.ComponentType == "procedure" {
		builder.Procedures[comp.Name.Id] = &firstProc
		continuationVar := DataVar{
			IDNode:   builder.newIDNode(),
			DataType: &ProcType{builder.newIDNode(), []DataVar{}},
		}
		firstProc.Vars = append(firstProc.Vars, continuationVar)
		proc.Jump = &CallProc{
			IDNode: builder.newIDNode(),
			Args:   []Data{},
			To:     &firstProc.Vars[len(firstProc.Vars)-1],
		}
	} else if comp.ComponentType == "transition" {
		builder.Transitions[comp.Name.Id] = &firstProc
	} else {
		panic(errors.New(fmt.Sprintf("Wrong Component type: %s", comp.ComponentType)))
	}
	builder.currentProc = ""
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
			dataVars[i] = DataVar{
				IDNode:   builder.newIDNode(),
				DataType: pType,
			}
			stackMapPush(builder.varStack, pName, &dataVars[i])
		}

		builder.Constructor = &Proc{
			IDNode:   builder.newIDNode(),
			ProcName: "constructor",
			Vars:     dataVars,
			Plan:     make([]Unit, len(n.Fields)),
		}
		for i, f := range n.Fields {
			n, d := builder.visitField(f)
			fmt.Println("Field", n, d)
			stackMapPush(builder.varStack, n, d)
			builder.Constructor.Plan[i] = &Save{
				IDNode: builder.newIDNode(),
				Slot:   n,
				Path:   []Data{},
				Data:   d,
			}
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
	case *ast.ExternalLibrary:
		fmt.Println("ExternalLibrary", n.Name.Id)
		//do nothing
	default:
		//fmt.Printf("CFGBuilder Visit unhandled type: %T\n", n)
		// do nothing
	}
	return builder
}

func (builder *CFGBuilder) initPrimitiveTypes() {
	stdLib := StdLib(builder)

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
	builder.kind = stdLib.star

	sizes := []int{32, 64, 128, 256}
	for _, s := range sizes {
		intName := "Int" + strconv.Itoa(s)
		intTyp := IntType{
			IDNode: builder.newIDNode(),
			Size:   s,
		}
		uintName := "Uint" + strconv.Itoa(s)
		uintTyp := NatType{
			IDNode: builder.newIDNode(),
			Size:   s,
		}
		builder.intWidthTypeMap[s] = &intTyp
		builder.natWidthTypeMap[s] = &uintTyp
		builder.primitiveTypeMap[intName] = &intTyp
		builder.primitiveTypeMap[uintName] = &uintTyp

		// Conversion functions
		var intAbsDD AbsDD
		intAbsTD := AbsTD{
			IDNode: builder.newIDNode(),
			Vars: []TypeVar{
				TypeVar{
					IDNode: builder.newIDNode(),
					Kind:   stdLib.star,
				},
			},
			Term: &intAbsDD,
		}
		intOpName := fmt.Sprintf("to_int%d", s)
		intAbsDD = AbsDD{
			IDNode: builder.newIDNode(),
			Vars: []DataVar{
				DataVar{
					IDNode:   builder.newIDNode(),
					DataType: &intAbsTD.Vars[0],
				},
			},
			Term: &Builtin{
				builder.newIDNode(),
				&AppTT{
					IDNode: builder.newIDNode(),
					Args:   []Type{&intTyp},
					To:     stdLib.Option,
				},
				intOpName,
			},
		}
		builder.builtinOpMap[intOpName] = &intAbsTD

		var uintAbsDD AbsDD
		uintAbsTD := AbsTD{
			IDNode: builder.newIDNode(),
			Vars: []TypeVar{
				TypeVar{
					IDNode: builder.newIDNode(),
					Kind:   stdLib.star,
				},
			},
			Term: &uintAbsDD,
		}
		uintOpName := fmt.Sprintf("to_uint%d", s)
		uintAbsDD = AbsDD{
			IDNode: builder.newIDNode(),
			Vars: []DataVar{
				DataVar{IDNode: builder.newIDNode(), DataType: &uintAbsTD.Vars[0]},
			},
			Term: &Builtin{
				builder.newIDNode(),
				&AppTT{
					IDNode: builder.newIDNode(),
					Args:   []Type{&uintTyp},
					To:     stdLib.Option,
				},
				uintOpName,
			},
		}
		builder.builtinOpMap[uintOpName] = &uintAbsTD

	}

	builder.primitiveTypeMap["Message"] = &MsgType{builder.newIDNode()}
	builder.primitiveTypeMap["String"] = &StrType{builder.newIDNode()}
	builder.primitiveTypeMap["ByStr"] = &RawType{builder.newIDNode(), -1}
	builder.primitiveTypeMap["ByStr32"] = &RawType{builder.newIDNode(), 32}
	builder.primitiveTypeMap["ByStr33"] = &RawType{builder.newIDNode(), 33}
	builder.primitiveTypeMap["ByStr64"] = &RawType{builder.newIDNode(), 64}
	builder.primitiveTypeMap["ByStr20"] = &RawType{builder.newIDNode(), 20}
	builder.primitiveTypeMap["BNum"] = &BnrType{builder.newIDNode()}

	intIntOps := []string{"add", "sub", "mul", "div", "rem"}
	intBoolOps := []string{"eq", "lt"}

	for _, bOp := range intIntOps {
		var bAbsDD AbsDD
		bAbsTD := AbsTD{
			IDNode: builder.newIDNode(),
			Vars: []TypeVar{
				TypeVar{
					IDNode: builder.newIDNode(),
					Kind:   stdLib.star,
				},
			},
			Term: &bAbsDD,
		}
		bAbsDD = AbsDD{
			IDNode: builder.newIDNode(),
			Vars: []DataVar{
				DataVar{
					IDNode:   builder.newIDNode(),
					DataType: &bAbsTD.Vars[0],
				},
				DataVar{
					IDNode:   builder.newIDNode(),
					DataType: &bAbsTD.Vars[0],
				},
			},
			Term: &Builtin{
				IDNode:      builder.newIDNode(),
				BuiltinType: &bAbsTD.Vars[0],
				Label:       bOp,
			},
		}
		builder.builtinOpMap[bOp] = &bAbsTD
	}

	for _, bOp := range intBoolOps {
		var bAbsDD AbsDD
		bAbsTD := AbsTD{
			IDNode: builder.newIDNode(),
			Vars: []TypeVar{
				TypeVar{
					IDNode: builder.newIDNode(),
					Kind:   stdLib.star,
				},
			},
			Term: &bAbsDD,
		}
		bAbsDD = AbsDD{
			IDNode: builder.newIDNode(),
			Vars: []DataVar{
				DataVar{IDNode: builder.newIDNode(), DataType: &bAbsTD.Vars[0]},
				DataVar{IDNode: builder.newIDNode(), DataType: &bAbsTD.Vars[0]},
			},
			Term: &Builtin{builder.newIDNode(), builder.definedADT["Bool"], bOp},
		}
		builder.builtinOpMap[bOp] = &bAbsTD
	}

	// builtin pow
	var powAbsDD AbsDD
	powAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{IDNode: builder.newIDNode(), Kind: stdLib.star},
			TypeVar{IDNode: builder.newIDNode(), Kind: stdLib.star},
		},
		Term: &powAbsDD,
	}
	powAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{IDNode: builder.newIDNode(), DataType: &powAbsTD.Vars[0]},
			DataVar{IDNode: builder.newIDNode(), DataType: &powAbsTD.Vars[1]},
		},
		Term: &Builtin{builder.newIDNode(), &powAbsTD.Vars[0], "pow"},
	}
	builder.builtinOpMap["pow"] = &powAbsTD

	// builtin concat
	strConcatOp := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
		},
		Term: &Builtin{builder.newIDNode(), builder.primitiveTypeMap["String"], "str_concat"},
	}
	builder.builtinOpMap["str_concat"] = &strConcatOp

	//builtin substr
	strSubstrOp := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.natWidthTypeMap[32],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.natWidthTypeMap[32],
			},
		},
		Term: &Builtin{builder.newIDNode(), builder.primitiveTypeMap["String"], "substr"},
	}
	builder.builtinOpMap["substr"] = &strSubstrOp

	// builtin to_string
	var toStrAbsDD AbsDD
	toStrAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star},
		},
		Term: &toStrAbsDD,
	}
	toStrAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &toStrAbsTD.Vars[0],
			},
		},
		Term: &Builtin{builder.newIDNode(), builder.primitiveTypeMap["String"], "to_string"},
	}
	builder.builtinOpMap["to_string"] = &toStrAbsTD

	//builtin strlen

	strlenOp := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
		},
		Term: &Builtin{builder.newIDNode(), builder.natWidthTypeMap[32], "strlen"},
	}
	builder.builtinOpMap["strlen"] = &strlenOp

	// builtin sha256hash
	var shaAbsDD AbsDD
	shaAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &shaAbsDD,
	}
	shaAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &shaAbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.primitiveTypeMap["ByStr32"],
			Label:       "sha256hash",
		},
	}
	builder.builtinOpMap["sha256hash"] = &shaAbsTD

	// builtin keccak256hash
	var keccakAbsDD AbsDD
	keccakAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &keccakAbsDD,
	}
	keccakAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &keccakAbsTD.Vars[0]},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.primitiveTypeMap["ByStr32"],
			Label:       "keccak256hash",
		},
	}
	builder.builtinOpMap["keccak256hash"] = &keccakAbsTD

	// builtin ripemd160hash
	var ripemdAbsDD AbsDD
	ripemdAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star},
		},
		Term: &ripemdAbsDD,
	}
	ripemdAbsDD = AbsDD{
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &ripemdAbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.primitiveTypeMap["ByStr20"],
			Label:       "ripemd160hash",
		},
	}
	builder.builtinOpMap["ripemd160hash"] = &ripemdAbsTD

	// builtin to_bystr
	var toBystrAbsDD AbsDD
	toBystrAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &toBystrAbsDD,
	}
	toBystrAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &toBystrAbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.primitiveTypeMap["ByStr"],
			Label:       "to_bystr",
		},
	}
	builder.builtinOpMap["to_bystr"] = &toBystrAbsTD

	// builtin to_uint256
	var touint256AbsDD AbsDD
	touint256AbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &touint256AbsDD,
	}
	touint256AbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &touint256AbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.natWidthTypeMap[256],
			Label:       "to_uint256",
		},
	}
	builder.builtinOpMap["to_uint256"] = &touint256AbsTD

	// schnorr_verify
	schnorrVerifyAbsDD := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr33"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr64"],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.definedADT["Bool"],
			Label:       "schnorr_verify",
		},
	}
	builder.builtinOpMap["schnorr_verify"] = &schnorrVerifyAbsDD

	// ecdsa_verify
	ecdsaVerifyAbsDD := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr33"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr64"],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.definedADT["Bool"],
			Label:       "ecdsa_verify",
		},
	}
	builder.builtinOpMap["ecdsa_verify"] = &ecdsaVerifyAbsDD

	// bech32_to_bystr20
	bech32ToBystr20AbsDD := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
		},
		Term: &Builtin{
			IDNode: builder.newIDNode(),
			BuiltinType: &AppTT{
				IDNode: builder.newIDNode(),
				Args:   []Type{builder.primitiveTypeMap["ByStr20"]},
				To:     stdLib.Option,
			},
			Label: "bech32_to_bystr20",
		},
	}
	builder.builtinOpMap["bech32_to_bystr20"] = &bech32ToBystr20AbsDD

	// bystr20_to_bech32
	bystr20ToBech32AbsDD := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["String"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["ByStr20"],
			},
		},
		Term: &Builtin{
			IDNode: builder.newIDNode(),
			BuiltinType: &AppTT{
				IDNode: builder.newIDNode(),
				Args:   []Type{builder.primitiveTypeMap["String"]},
				To:     stdLib.Option,
			},
			Label: "bystr20_to_bech32",
		},
	}

	builder.builtinOpMap["bystr20_to_bech32"] = &bystr20ToBech32AbsDD

	//Maps

	// put
	var putAbsDD AbsDD
	putAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star},
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star},
		},
		Term: &putAbsDD,
	}
	putAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.setdefaultMap(&putAbsTD.Vars[0], &putAbsTD.Vars[1]),
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &putAbsTD.Vars[0],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &putAbsTD.Vars[1],
			},
		},
	}
	putAbsDD.Term = &Builtin{
		IDNode:      builder.newIDNode(),
		BuiltinType: putAbsDD.Vars[0].DataType,
		Label:       "put",
	}
	builder.builtinOpMap["put"] = &putAbsTD

	// get
	var getAbsDD AbsDD
	getAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &getAbsDD,
	}
	getAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.setdefaultMap(&getAbsTD.Vars[0], &getAbsTD.Vars[1]),
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &getAbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			IDNode: builder.newIDNode(),
			BuiltinType: &AppTT{
				IDNode: builder.newIDNode(),
				Args:   []Type{&getAbsTD.Vars[1]},
				To:     stdLib.Option,
			},
			Label: "get",
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
				IDNode:   builder.newIDNode(),
				DataType: builder.setdefaultMap(&containsAbsTD.Vars[0], &containsAbsTD.Vars[1]),
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &containsAbsTD.Vars[0],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.definedADT["Bool"],
			Label:       "contains",
		},
	}
	builder.builtinOpMap["contains"] = &containsAbsTD

	// remove
	var removeAbsDD AbsDD
	removeAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &removeAbsDD,
	}
	removeAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.setdefaultMap(&removeAbsTD.Vars[0], &removeAbsTD.Vars[1]),
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &removeAbsTD.Vars[0],
			},
		},
	}
	removeAbsDD.Term = &Builtin{
		IDNode:      builder.newIDNode(),
		BuiltinType: removeAbsDD.Vars[0].DataType,
		Label:       "remove",
	}
	builder.builtinOpMap["remove"] = &removeAbsTD

	// size
	var sizeAbsDD AbsDD
	sizeAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &sizeAbsDD,
	}
	sizeAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.setdefaultMap(&sizeAbsTD.Vars[0], &sizeAbsTD.Vars[1]),
			},
		},
	}
	sizeAbsDD.Term = &Builtin{
		IDNode:      builder.newIDNode(),
		BuiltinType: builder.definedADT["Bool"],
		Label:       "size",
	}
	builder.builtinOpMap["size"] = &sizeAbsTD

	// builtin blt
	bltAbsDD := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["BNum"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["BNum"],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.definedADT["Bool"],
			Label:       "blt",
		},
	}
	builder.builtinOpMap["blt"] = &bltAbsDD

	// builtin badd
	var baddAbsDD AbsDD
	baddAbsTD := AbsTD{
		IDNode: builder.newIDNode(),
		Vars: []TypeVar{
			TypeVar{
				IDNode: builder.newIDNode(),
				Kind:   stdLib.star,
			},
		},
		Term: &baddAbsDD,
	}
	baddAbsDD = AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: &baddAbsTD.Vars[0],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["BNum"],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.primitiveTypeMap["BNum"],
			Label:       "badd",
		},
	}
	builder.builtinOpMap["badd"] = &baddAbsTD

	// builtin bsub
	bsubAbsDD := AbsDD{
		IDNode: builder.newIDNode(),
		Vars: []DataVar{
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["BNum"],
			},
			DataVar{
				IDNode:   builder.newIDNode(),
				DataType: builder.primitiveTypeMap["BNum"],
			},
		},
		Term: &Builtin{
			IDNode:      builder.newIDNode(),
			BuiltinType: builder.primitiveTypeMap["Int256"],
			Label:       "bsub",
		},
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
		typeVarStack:     map[string][]Type{},
		Constructor:      nil,
		Transitions:      map[string]*Proc{},
		Procedures:       map[string]*Proc{},
		fieldTypeMap:     map[string]Type{},

		mapTypeMap:              map[Type]map[Type]*MapType{},
		genericTypeConstructors: map[string]*AbsTT{},
		genericDataConstructors: map[string]*AbsTD{},
		nodeCounter:             1,
	}

	builder.initPrimitiveTypes()
	ast.Walk(&builder, n)

	return &builder
}
