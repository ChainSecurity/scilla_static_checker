package ir

import (
	"errors"
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"strconv"
)

type CFGBuilder struct {
	//Names []string
	typeMap map[string]Type
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

func (builder *CFGBuilder) visitLibEntry(le *ast.LibEntry) {
	switch n := (*le).(type) {
	case *ast.LibraryVariable:
		//fmt.Printf("skipping %T\n", n)
	case *ast.LibraryType:
		name := n.Name.Id
		//fmt.Printf("t %T %s\n", n, name)
		typ := EnumType{}
		for _, ctr := range n.CtrDefs {
			name, types := builder.visitCtr(ctr)
			typ[name] = types
		}
		builder.typeMap[name] = &typ
	default:
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func (builder *CFGBuilder) Visit(node ast.AstNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.LibEntry:
		builder.visitLibEntry(n)
	default:
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
		builder.typeMap[intName] = &intTyp
		builder.typeMap[uintName] = &uintTyp
	}
	builder.typeMap["ByStr20"] = &RawType{20}
	builder.typeMap["ByStr32"] = &RawType{32}

	stdLib := StdLib()
	builder.typeMap["Bool"] = stdLib.Boolean

}

func BuildCFG(n ast.AstNode) bool {
	builder := CFGBuilder{map[string]Type{}}
	builder.initPrimitiveTypes()
	ast.Walk(&builder, n)
	fmt.Println("-----------------------")
	for k, v := range builder.typeMap {
		fmt.Println(k, v)
	}
	fmt.Println("-----------------------")
	return true
}
