package ir

import (
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
)

type CFGBuilder struct {
	Names []string
}

func (builder *CFGBuilder) Visit(node ast.AstNode) ast.Visitor {
	if e, ok := node.(*ast.LibEntry); ok {
		d := *e
		if lt, ok := d.(*ast.LibraryType); ok {
			typeName := lt.Name.Id
			fmt.Println(typeName)
			builder.Names = append(builder.Names, typeName)
			for _, c := range lt.CtrDefs {
				fmt.Printf("\t %s %s\n", c.CDName.Id, c.CArgTypes)
			}
			fmt.Println(builder.Names)
		}
	}
	return builder
}

func BuildCFG(n ast.AstNode) bool {
	builder := CFGBuilder{}
	ast.Walk(&builder, n)
	return true
}
