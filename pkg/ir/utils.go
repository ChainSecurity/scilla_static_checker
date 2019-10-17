package ir

import (
	"errors"
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"strconv"
)

type CFGBuilder struct {
	GlobalVarMap map[string]Data

	constrTypeMap   map[string]string
	intWidthTypeMap map[int]*IntType
	natWidthTypeMap map[int]*NatType
	typeMap         map[string]Type
	varStack        map[string][]Data
}

func stackMapPush(s map[string][]Data, k string, v Data) {
	s[k] = append(s[k], v)
}

func stackMapPop(s map[string][]Data, k string) {
	n := len(s[k]) - 1
	s[k] = s[k][:n]
}

func stackMapPeek(s map[string][]Data, k string) (Data, bool) {
	l := len(s[k])
	if l == 0 {
		return nil, false
	}
	return s[k][l-1], true
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
			When: &When{
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
		vars := make([]Data, len(n.Args))
		for i, a := range n.Args {
			v, ok := stackMapPeek(builder.varStack, a.Id)
			if !ok {
				panic(errors.New(fmt.Sprintf("variable not found: %s", a.Id)))
			}
			vars[i] = v
		}
		return nil
	case *ast.LetExpression:
		var_name := n.Var.Id
		expr := builder.visitExpression(n.Expr)
		stackMapPush(builder.varStack, var_name, expr)
		defer stackMapPop(builder.varStack, var_name)
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
		builder.GlobalVarMap[name] = v
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

func visitParams(b *CFGBuilder, p *ast.Parameter) (string, Type) {
	name := p.Name.Id
	t := b.visitASTType(p.Type)
	return name, t
}

func visitField(builder *CFGBuilder, f *ast.Field) (string, Data) {
	name := f.Name.Id
	//t := builder.typeMap[f.Type]
	data := builder.visitExpression(f.Expr)
	return name, data
}

func (builder *CFGBuilder) Visit(node ast.AstNode) ast.Visitor {
	switch n := node.(type) {
	case ast.LibEntry:
		builder.visitLibEntry(n)
	case *ast.Contract:
		//name := n.Name.Id
		vars := make([]string, len(n.Params))
		params := map[string]Type{}
		for i, p := range n.Params {
			pName, pType := visitParams(builder, p)
			params[pName] = pType
			vars[i] = pName
			dataVar := DataVar{pType}
			stackMapPush(builder.varStack, pName, &dataVar)
		}

		for _, f := range n.Fields {
			n, d := visitField(builder, f)
			stackMapPush(builder.varStack, n, d)
		}

		for _, pName := range vars {
			stackMapPop(builder.varStack, pName)
		}

		//t := n.Type
		//expr := builder.visitExpression(n.Expr)
		//builder.Field[name] = Save{}
		//builder.visitLibEntry(n)
	default:
		//fmt.Printf("Unhandled type: %T\n", n)
		// do nothing
	}
	return builder
}

func (builder *CFGBuilder) initPrimitiveTypes() {
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
	}
	builder.typeMap["ByStr20"] = &RawType{20}
	builder.typeMap["ByStr32"] = &RawType{32}
	builder.typeMap["BNum"] = &BnrType{}
	builder.typeMap["Message"] = &MsgType{}

	stdLib := StdLib()
	builder.typeMap["Bool"] = stdLib.Boolean
	//builder.constructorTypeMap["Bool"] = stdLib.Boolean

}

func BuildCFG(n ast.AstNode) *CFGBuilder {
	builder := CFGBuilder{
		GlobalVarMap:    map[string]Data{},
		typeMap:         map[string]Type{},
		constrTypeMap:   map[string]string{},
		intWidthTypeMap: map[int]*IntType{},
		natWidthTypeMap: map[int]*NatType{},
		varStack:        map[string][]Data{},
	}
	builder.initPrimitiveTypes()

	//stdlib := StdLib()
	//builder.typeMap[] =

	ast.Walk(&builder, n)

	//a := builder.typeMap["Piece"]
	//switch s := a.(type) {
	//case *EnumType:
	//d := *s
	//for k, a := range d {
	//fmt.Println(k, a)
	//}
	//default:
	//fmt.Printf("%T\n", s)
	//}
	//fmt.Println("-----------------------")
	//for k, v := range builder.typeMap {
	//fmt.Println(k, v)
	//}
	//fmt.Println("-----------------------")
	return &builder
}
