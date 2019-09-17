package ast

import (
	"errors"
	"fmt"
	"log"
	//"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
)

type Visitor interface {
	Visit(node AstNode) (w Visitor)
}

func Walk(v Visitor, node AstNode) {
	if v = v.Visit(node); v == nil {
		return
	}
	switch n := node.(type) {
	case *Location:
	// do nothing
	case *ExternalLibrary:
		Walk(v, n.Name)
		if n.Alias != nil {
			Walk(v, n.Alias)
		}
	case *Identifier:
		Walk(v, n.Loc)
	case *Parameter:
		Walk(v, n.Name)
	case *Field:
		Walk(v, n.Name)
		Walk(v, n.Expr)
	case *Contract:
		Walk(v, n.Name)
		for _, x := range n.Params {
			Walk(v, x)
		}
		for _, x := range n.Fields {
			Walk(v, x)
		}
		for _, x := range n.Components {
			Walk(v, x)
		}
	case *Library:
		Walk(v, n.Name)
		for _, x := range n.Entries {
			Walk(v, x)
		}
	case *ContractModule:
		Walk(v, n.Name)
		if n.Library != nil {
			Walk(v, n.Library)
		}
		for _, x := range n.ELibs {
			Walk(v, x)
		}
		Walk(v, n.C)
	case *Component:
		Walk(v, n.Name)
		for _, x := range n.Params {
			Walk(v, x)
		}
		for _, x := range n.Body {
			Walk(v, x)
		}
	case *MessageArgument:
		Walk(v, n.Pl)
	case *CtrDef:
		Walk(v, n.CDName)
	case *Builtin:
		Walk(v, n.Loc)
	case *MatchExpressionCase:
		Walk(v, n.Pat)
		Walk(v, n.Expr)
	case *MatchStatementCase:
		Walk(v, n.Pat)
		for _, x := range n.Body {
			Walk(v, x)
		}
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

type inspector func(AstNode) bool

func (f inspector) Visit(node AstNode) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
//
func Inspect(node AstNode, f func(AstNode) bool) {
	Walk(inspector(f), node)
}
