package ir

import (
	"errors"
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"strconv"
)

type CFGBuilder struct {
	//Names []string
	typeMap            map[string]Type
	GlobalVarMap       map[string]Data
	constructorTypeMap map[string]string
	intWidthTypeMap    map[int]*IntType
	natWidthTypeMap    map[int]*NatType

	varMap map[string]Data
}

func (builder *CFGBuilder) visitCtr(ctr *ast.CtrDef) (string, []Type) {
	name := ctr.CDName.Id
	var types []Type
	for _, typName := range ctr.CArgTypes {
		typ, ok := builder.typeMap[typName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown type: %s", typName)))
		}
		types = append(types, typ)
	}
	return name, types
}

func (builder *CFGBuilder) visitPattern(p *ast.Pattern, t Type) *Bind {
	switch pat := (*p).(type) {
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

func (builder *CFGBuilder) visitLiteral(l *ast.Literal) Data {
	switch lit := (*l).(type) {
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
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	case *ast.ADTValLiteral:
		panic(errors.New(fmt.Sprintf("Not implemented: %T", lit)))
	}
	return nil
}

func (builder *CFGBuilder) visitExpression(e *ast.Expression) Data {
	switch n := (*e).(type) {
	case *ast.ConstrExpression:
		constructorName := n.ConstructorName
		typeName, ok := builder.constructorTypeMap[constructorName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown constructor: %s", constructorName)))
		}
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
			Case:     constructorName,
			Data:     d,
		}
		return &res
	case *ast.FunExpression:
		lhs := n.Lhs.Id
		typeName := n.LhsType
		lhsTyp, ok := builder.typeMap[typeName]
		if !ok {
			panic(errors.New(fmt.Sprintf("Unknown type: %s", typeName)))
		}
		dataVar := DataVar{lhsTyp}
		builder.varMap[lhs] = &dataVar
		//fmt.Printf("FunExp:{ var:%s\n\ttypeName:%s}\n", lhs, typeName)
		rhs := builder.visitExpression(n.RhsExpr)
		//fmt.Printf("%T\n", rhs)
		return &AbsDD{Vars: []*DataVar{&dataVar}, Term: rhs}
	case *ast.MatchExpression:
		lhs := n.Lhs.Id
		data, ok := builder.varMap[lhs]
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
	default:
		fmt.Printf("Unhandled Expression type: %T\n", n)
		//panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
		return nil
	}
}

func (builder *CFGBuilder) visitLibEntry(le *ast.LibEntry) {
	switch n := (*le).(type) {
	case *ast.LibraryVariable:
		name := n.Name.Id
		v := builder.visitExpression(n.Expr)
		builder.GlobalVarMap[name] = v
	case *ast.LibraryType:
		typeName := n.Name.Id
		typ := EnumType{}
		for _, ctr := range n.CtrDefs {
			constructorName, types := builder.visitCtr(ctr)
			typ[constructorName] = types
			builder.constructorTypeMap[constructorName] = typeName
		}
		builder.typeMap[typeName] = &typ
	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func (builder *CFGBuilder) Visit(node ast.AstNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.LibEntry:
		builder.visitLibEntry(n)
	default:
		//fmt.Printf("%T\n", n)
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

	stdLib := StdLib()
	builder.typeMap["Bool"] = stdLib.Boolean

}

func BuildCFG(n ast.AstNode) *CFGBuilder {
	builder := CFGBuilder{
		typeMap:            map[string]Type{},
		GlobalVarMap:       map[string]Data{},
		constructorTypeMap: map[string]string{},
		intWidthTypeMap:    map[int]*IntType{},
		natWidthTypeMap:    map[int]*NatType{},
		varMap:             map[string]Data{},
	}
	builder.initPrimitiveTypes()
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
