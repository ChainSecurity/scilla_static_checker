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

func walkPayload(v Visitor, p *Payload) {
	switch n := (*p).(type) {
	case *PayloadLitral:
		Walk(v, n.Lit)
	case *PayloadVariable:
		Walk(v, n.Val)
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func walkExpression(v Visitor, p *Expression) {
	switch n := (*p).(type) {
	case *LiteralExpression:
		Walk(v, n.Loc)
		Walk(v, n.Val)
	case *VarExpression:
		Walk(v, n.Loc)
		Walk(v, n.Var)
	case *LetExpression:
		Walk(v, n.Loc)
		Walk(v, n.Var)
		Walk(v, n.Expr)
		Walk(v, n.Body)
	case *MessageExpression:
		Walk(v, n.Loc)
		for _, x := range n.MArgs {
			Walk(v, x)
		}
	case *FunExpression:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
		Walk(v, n.RhsExpr)
	case *AppExpression:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
		for _, x := range n.RhsList {
			Walk(v, x)
		}
	case *ConstrExpression:
		Walk(v, n.Loc)
		for _, x := range n.Args {
			Walk(v, x)
		}
	case *MatchExpression:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
		for _, x := range n.Cases {
			Walk(v, x)
		}
	case *BuiltinExpression:
		Walk(v, n.Loc)
		for _, x := range n.Args {
			Walk(v, x)
		}
		Walk(v, n.Bf)
	case *TFunExpression:
		Walk(v, n.Loc)
	case *TAppExpression:
		Walk(v, n.Loc)
	case *FixpointExpression:
		Walk(v, n.Loc)
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func walkLiteral(v Visitor, p *Literal) {
	switch n := (*p).(type) {
	case *StringLiteral:
		// do nothing
	case *BNumLiteral:
		// do nothing
	case *ByStrLiteral:
		// do nothing
	case *ByStrXLiteral:
		// do nothing
	case *IntLiteral:
		// do nothing
	case *UintLiteral:
		// do nothing
	case *MapLiteral:
		for _, x := range n.MVals {
			Walk(v, x)
		}
	case *ADTValLiteral:
		// do nothing
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func walkPattern(v Visitor, p *Pattern) {
	switch n := (*p).(type) {
	case *WildcardPattern:
		// do nothing
	case *BinderPattern:
		Walk(v, n.Variable)
	case *ConstructorPattern:
		for _, x := range n.Pats {
			Walk(v, x)
		}
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

func walkStatement(v Visitor, p *Statement) {
	switch n := (*p).(type) {
	case *LoadStatement:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
		Walk(v, n.Rhs)
	case *StoreStatement:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
		Walk(v, n.Rhs)
	case *BindStatement:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
		Walk(v, n.RhsExpr)
	case *MapUpdateStatement:
		Walk(v, n.Loc)
		Walk(v, n.Name)
		if n.Rhs != nil {
			Walk(v, n.Rhs)
		}
		for _, x := range n.Keys {
			Walk(v, x)
		}
	case *MapGetStatement:
		Walk(v, n.Loc)
		Walk(v, n.Name)
		Walk(v, n.Lhs)
		for _, x := range n.Keys {
			Walk(v, x)
		}
	case *MatchStatement:
		Walk(v, n.Loc)
		Walk(v, n.Arg)
		for _, x := range n.Cases {
			Walk(v, x)
		}
	case *ReadFromBCStatement:
		Walk(v, n.Loc)
		Walk(v, n.Lhs)
	case *AcceptPaymentStatement:
		Walk(v, n.Loc)
	case *SendMsgsStatement:
		Walk(v, n.Loc)
		Walk(v, n.Arg)
	case *CreateEvntStatement:
		Walk(v, n.Loc)
		Walk(v, n.Arg)
	case *CallProcStatement:
		Walk(v, n.Loc)
		Walk(v, n.Arg)
		for _, x := range n.Messages {
			Walk(v, x)
		}
	case *ThrowStatement:
		Walk(v, n.Loc)
		if n.Arg != nil {
			Walk(v, n.Arg)
		}
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}
func walkLibEntry(v Visitor, p *LibEntry) {
	switch n := (*p).(type) {
	case *LibraryVariable:
		Walk(v, n.Expr)
	case *LibraryType:
		fmt.Printf("Library %s\n", n.Name.Id)
		for _, x := range n.CtrDefs {
			Walk(v, x)
		}
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
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
	case *LibraryModule:
		Walk(v, n.Library)
		for _, x := range n.ELibs {
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
		fmt.Printf("\tCtrDef %s\n", n.CDName.Id)
		for _, x := range n.CArgTypes {
			fmt.Printf("\t\t %s\n", x)
		}
		Walk(v, n.CDName)
	case *Builtin:
		Walk(v, n.Loc)
	case *MatchExpressionCase:
		Walk(v, n.Pat)
		Walk(v, n.Expr)
	case *Pattern:
		walkPattern(v, n)
	case *Literal:
		walkLiteral(v, n)
	case *Statement:
		walkStatement(v, n)
	case *Payload:
		walkPayload(v, n)
	case *Expression:
		walkExpression(v, n)
	case *LibEntry:
		walkLibEntry(v, n)
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
