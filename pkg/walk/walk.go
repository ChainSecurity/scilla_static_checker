package walk

import (
    //"fmt"
    //"log"
    //"errors"
    "gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
)

type Visitor interface {
    Visit(node ast.AstNode) (w Visitor)
}

//func walkGenericLibEntry(v Visitor, n *GenericLibEntry){
    //t := n.NodeType
    //switch t{
    //case "LibraryVariable":
        //Walk(v, n.Expr)
    //case "LibraryType":
        //for _, x := range n.CtrDefs {
            //Walk(v, x)
        //}
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled GenericLibEntry type: %s", t)))
    //}
//}
//func walkGenericPayload(v Visitor, n *GenericPayload){
    //t := n.NodeType
    //switch t{
    //case "PayloadLitral":
        //Walk(v, n.Lit)
    //case "PayloadVariable":
        //Walk(v, n.Val)
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled GenericPattern type: %s", t)))
    //}
//}

//func walkGenericPattern(v Visitor, n *GenericPattern){
    //t := n.NodeType
    //switch t{
    //case "WildcardPattern":
        ////do nothing
    //case "BinderPattern":
        //Walk(v, n.Variable)
    //case "ConstructorPattern":
        //for _, x := range n.Pats {
            //Walk(v, x)
        //}
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled GenericPattern type: %s", t)))
    //}
//}

//func walkGenericLiteral(v Visitor, n *GenericLiteral) {
    //t := n.NodeType
    //switch t{
    //case "StringLiteral":
        ////do nothing
    //case "BNumLiteral":
        ////do nothing
    //case "ByStrLiteral":
        ////do nothing
    //case "ByStrXLiteral":
        ////do nothing
    //case "IntLiteral":
        ////do nothing
    //case "UintLiteral":
        ////do nothing
    //case "MapLiteral":
        //for _, x := range n.MVals {
            //Walk(v, x)
        //}
    //case "ADTValLiteral":
        ////do nothing
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled GenericExpression type: %s", t)))
    //}

//}

//func walkGenericStatement(v Visitor, n *GenericStatement) {
    //Walk(v, n.Loc)
    //t := n.NodeType
    //switch t{
    //case "LoadStatement": case "StoreStatement":
        //Walk(v, n.Lhs)
        //Walk(v, n.Rhs)
    //case "BindStatement":
        //Walk(v, n.Lhs)
        //Walk(v, n.RhsExpr)
    //case "MapUpdateStatement":
        //Walk(v, n.Name)
        //if n.Rhs != nil {
            //Walk(v, n.Rhs)
        //}
        //for _, x := range n.Keys {
            //Walk(v, x)
        //}
    //case "MapGetStatement":
        //Walk(v, n.Name)
        //Walk(v, n.Lhs)
        //for _, x := range n.Keys {
            //Walk(v, x)
        //}
    //case "MatchStatement":
        //Walk(v, n.Arg)
        //for _, x := range n.Cases {
            //Walk(v, x)
        //}
    //case "ReadFromBCStatement":
        //Walk(v, n.Lhs)
    //case "AcceptPaymentStatement":
    //case "SendMsgsStatement":
        //Walk(v, n.Arg)
    //case "CreateEvntStatement":
        //Walk(v, n.Arg)
    //case "CallProcStatement":
        //Walk(v, n.Arg)
        //for _, x := range n.Messages {
            //Walk(v, x)
        //}
    //case "ThrowStatement":
        //if n.Arg != nil {
            //Walk(v, n.Arg)
        //}
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled GenericStatement type: %s", t)))
    //}
//}
//func walkGenericExpression(v Visitor, n *GenericExpression) {
    //Walk(v, n.Loc)
    //t := n.NodeType 
    //switch t{
    //case "LiteralExpression":
        //Walk(v, n.Val)
    //case "VarExpression":
        //Walk(v, n.Variable)
    //case "LetExpression":
        //Walk(v, n.Variable)
        //Walk(v, n.Expr)
        //Walk(v, n.Body)
    //case "MessageExpression":
        //for _, x := range n.MArgs {
            //Walk(v, x)
        //}
    //case "FunExpression":
        //Walk(v, n.Lhs)
        //Walk(v, n.RhsExpr)
    //case "AppExpression":
        //Walk(v, n.Lhs)
        //for _, x := range n.RhsList {
            //Walk(v, x)
        //}
    //case "ConstrExpression":
        //for _, x := range n.Args {
            //Walk(v, x)
        //}
    //case "MatchExpression":
        //Walk(v, n.Lhs)
        //for _, x := range n.Cases {
            //Walk(v, x)
        //}
    //case "BuiltinExpression":
        //for _, x := range n.Args {
            //Walk(v, x)
        //}
        //Walk(v, n.Bf)
    //case "TFunExpression", "TAppExpression", "FixpointExpression":
        ////do nothing
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled GenericExpression type: %s", t)))
    //}
//}

//func Walk(v Visitor, node AstNode) {
    //if v = v.Visit(node); v == nil {
        //return
    //}
    //switch n := node.(type) {
    //case *Location:
        //// do nothing
    //case *ExternalLibrary:
        //Walk(v, n.Name)
        //if n.Alias != nil {
            //Walk(v, n.Alias)
        //}
    //case *Identifier:
        //Walk(v, n.Loc)
    //case *Parameter:
        //Walk(v, n.Name)
    //case *Field:
        //Walk(v, n.Name)
        //Walk(v, n.Expr)
    //case *Contract:
        //Walk(v, n.Name)
        //for _, x := range n.Params {
            //Walk(v, x)
        //}
        //for _, x := range n.Fields {
            //Walk(v, x)
        //}
        //for _, x := range n.Components {
            //Walk(v, x)
        //}
    //case *Library:
        //Walk(v, n.Name)
        //for _, x := range n.Entries {
            //Walk(v, x)
        //}
    //case *ContractModule:
        //Walk(v, n.Name)
        //if n.Library != nil {
            //Walk(v, n.Library)
        //}
        //for _, x := range n.ELibs {
            //Walk(v, x)
        //}
        //Walk(v, n.C)
    //case *Component:
        //Walk(v, n.Name)
        //for _, x := range n.Params {
            //Walk(v, x)
        //}
        //for _, x := range n.Body {
            //Walk(v, x)
        //}
    //case *GenericLibEntry:
        //walkGenericLibEntry(v, n)
    //case *MessageArgument:
        //Walk(v, n.P)
    //case *CtrDef:
        //Walk(v, n.CDName)
    //case *Builtin:
        //Walk(v, n.Loc)
    //case *GenericLiteral:
        //walkGenericLiteral(v, n)
    //case *GenericStatement:
        //walkGenericStatement(v, n)
    //case *GenericExpression:
        //walkGenericExpression(v, n)
    //case *MatchExpressionCase:
        //Walk(v, n.Pat)
        //Walk(v, n.Expr)
    //case *MatchStatementCase:
        //Walk(v, n.Pat)
        //for _, x := range n.Body {
            //Walk(v, x)
        //}
    //case *GenericPattern:
        //walkGenericPattern(v, n)
    //case *GenericPayload:
        //walkGenericPayload(v, n)
    //default:
        //log.Fatal(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
    //}
//}

//type inspector func(AstNode) bool

//func (f inspector) Visit(node AstNode) Visitor {
    //if f(node) {
        //return f
    //}
    //return nil
//}

//// Inspect traverses an AST in depth-first order: It starts by calling
//// f(node); node must not be nil. If f returns true, Inspect invokes f
//// recursively for each of the non-nil children of node, followed by a
//// call of f(nil).
////
//func Inspect(node AstNode, f func(AstNode) bool) {
    //Walk(inspector(f), node)
//}
