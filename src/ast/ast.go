package ast

import (
)


type AstNode interface {
}

type Literal interface {
    AstNode
    litNode()
}

type Expression interface {
    AstNode
    exprNode()
}

type Statement interface {
    AstNode
    stmtNode()
}

type Payload interface {
    AstNode
    payloadNode()
}

type LibEntry interface {
    AstNode
    libEntryNode()
}

type Pattern interface {
    AstNode
    patternNode()
}

type Location struct{
    SourceFile string `json:"source_file"`
    Line int `json:"line"`
    Column int `json:"column"`
}


type Identifier struct{
    Loc *Location `json:"loc"`
    Id string `json:"identifier"`
}

type MapVal struct{
    Key *Literal `json:"key"`
    Val *Literal `json:"value"`
}

type CtrDef struct {
    CDName *Identifier `json:"ctr_def_name"`
    CArgTypes []string `json:"c_arg_types"`
}

type StringLiteral struct{
    Val string `json:"value"`
}

type BNumLiteral struct{
    Val string `json:"value"`
}

type ByStrLiteral struct{
    Val string `json:"value"`
}

type ByStrXLiteral struct{
    Val string `json:"value"`

}
type IntLiteral struct{
    Val string `json:"value"`

}
type UintLiteral struct{
    Val string `json:"value"`
}

type MapLiteral struct{
    KeyType string `json:"key_type"`
    ValType string `json:"value_type"`
    MVals []MapVal `json:"mvalues"`
}

type ADTValLiteral struct {
}

func (*StringLiteral) litNode() {}
func (*BNumLiteral) litNode() {}
func (*ByStrLiteral) litNode() {}
func (*ByStrXLiteral) litNode() {}
func (*IntLiteral) litNode() {}
func (*UintLiteral) litNode() {}
func (*MapLiteral) litNode() {}
func (*ADTValLiteral) litNode() {}

type AnnotatedNode struct{
    Loc *Location `json:"loc"`
}

type LiteralExpression struct{
    AnnotatedNode
    Val *Literal `json:"value"`
}

type VarExpression struct{
    AnnotatedNode
    Var *Identifier `json:"variable"`
}

type LetExpression struct{
    AnnotatedNode
    Var *Identifier `json:"variable"`
    VarType string `json:"variable_type"` //Optional 
    Expr *Expression `json:"expression"`
    Body *Expression `json:"body"`
}

type MessageExpression struct{
    AnnotatedNode
    MArgs []*MessageArgument `json:"margs"`
}

type FunExpression struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    RhsExpr *Expression `json:"rhs_expr"`
    FunType string `json:"fun_type"`
}

type AppExpression struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    RhsList []*Identifier `json:"rhs_list"`
}

type ConstrExpression struct{
    AnnotatedNode
    Types []string `json:"types"`
    ConstructorName string `json:"constructor_name"`
    Args []*Identifier `json:"args"`
}

type MatchExpression struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    Cases []*MatchExpressionCase `json:"cases"`
}

type BuiltinExpression struct{
    AnnotatedNode
    Args []*Identifier `json:"args"`
    Bf *Builtin `json:"builtin_function"`
}

type TFunExpression struct{
    AnnotatedNode
}

type TAppExpression struct{
    AnnotatedNode
}

type FixpointExpression struct{
    AnnotatedNode
}

func (*LiteralExpression) exprNode() {}
func (*VarExpression) exprNode() {}
func (*LetExpression) exprNode() {}
func (*MessageExpression) exprNode() {}
func (*FunExpression) exprNode() {}
func (*AppExpression) exprNode() {}
func (*ConstrExpression) exprNode() {}
func (*MatchExpression) exprNode() {}
func (*BuiltinExpression) exprNode() {}
func (*TFunExpression) exprNode() {}
func (*TAppExpression) exprNode() {}
func (*FixpointExpression) exprNode() {}


type PayloadLitral struct{
    Lit *Literal `json:"literal"`
}

type PayloadVariable struct{
    Val *Identifier `json:"value"`
}


func (*PayloadLitral) payloadNode() {}
func (*PayloadVariable) payloadNode() {}

type MessageArgument struct{
    Var string `json:"variable"`
    Pl *Payload `json:"payload"`
}

type WildcardPattern struct{
}

type BinderPattern struct{
    Variable *Identifier `json:"variable"`
}

type ConstructorPattern struct{
    ConstrName string `json:"constructor_name"`
    Pats []*Pattern `json:"patterns"`
}

func (*WildcardPattern) patternNode() {}
func (*BinderPattern) patternNode() {}
func (*ConstructorPattern) patternNode() {}

type MatchExpressionCase struct{
    Pat *Pattern `json:"pattern"`
    Expr *Expression `json:"expression"`
}

type MatchStatementCase struct{
    Pat *Pattern `json:"pattern"`
    Body []*Statement `json:"pattern_body"`
}
type Builtin struct{
    Loc *Location
    Type string
}

type LoadStatement struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    Rhs *Identifier `json:"rhs"`

}
type StoreStatement struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    Rhs *Identifier `json:"rhs"`

}
type BindStatement struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    RhsExpr *Expression `json:"rhs_expr"`

}
type MapUpdateStatement struct{
    AnnotatedNode
    Name *Identifier `json:"map_name"`
    Rhs *Identifier `json:"rhs"` //Optional
    Keys []*Identifier `json:"keys"`

}
type MapGetStatement struct{
    AnnotatedNode
    Name *Identifier `json:"map_name"`
    Lhs *Identifier `json:"lhs"`
    Keys []*Identifier `json:"keys"`
    IsValRetrieve bool `json:"is_value_retrieve"`
}

type MatchStatement struct{
    AnnotatedNode
    Arg *Identifier `json:"arg"`
    Cases []*MatchStatementCase `json:"cases"`

}
type ReadFromBCStatement struct{
    AnnotatedNode
    Lhs *Identifier `json:"lhs"`
    RhsStr string `json:"rhs_str"`

}
type AcceptPaymentStatement struct{
    AnnotatedNode
}

type SendMsgsStatement struct{
    AnnotatedNode
    Arg *Identifier `json:"arg"`

}
type CreateEvntStatement struct{
    AnnotatedNode
    Arg *Identifier `json:"arg"`
}

type CallProcStatement struct{
    AnnotatedNode
    Arg *Identifier `json:"arg"`
    Messages []*Identifier `json:"messages"`

}
type ThrowStatement struct{
    AnnotatedNode
    Arg *Identifier `json:"arg"` // Optional
}

func (*LoadStatement) stmtNode() {}
func (*StoreStatement) stmtNode() {}
func (*BindStatement) stmtNode() {}
func (*MapUpdateStatement) stmtNode() {}
func (*MapGetStatement) stmtNode() {}
func (*MatchStatement) stmtNode() {}
func (*ReadFromBCStatement) stmtNode() {}
func (*AcceptPaymentStatement) stmtNode() {}
func (*SendMsgsStatement) stmtNode() {}
func (*CreateEvntStatement) stmtNode() {}
func (*CallProcStatement) stmtNode() {}
func (*ThrowStatement) stmtNode() {}

type LibraryVariable struct{
    VarType string `json:"variable_type"` // Optional
    Expr *Expression `json:"expression"`
}

type LibraryType struct{
    CtrDefs []*CtrDef `json:"ctr_defs"`
}

func (*LibraryVariable) libEntryNode() {}
func (*LibraryType) libEntryNode() {}

type Library struct{
    Name *Identifier  `json:"library_name"`
    Entries []*LibEntry `json:"library_entries"`
}

type ExternalLibrary struct{
    Name *Identifier `json:"name"`
    Alias *Identifier `json:"alias"` // Optional

}
type ContractModule struct{
    ScillaMajorVersion int `json:"scilla_major_version"`
    Name *Identifier `json:"name"`
    Library *Library `json:"library"` // Optional
    ELibs []*ExternalLibrary `json:"external_libraries"`
    C *Contract `json:"contract"`
}

type Field struct{
    Name *Identifier `json:"field_name"`
    Type string `json:"field_type"`
    Expr *Expression `json:"expression"`
}

type Parameter struct{
    Name *Identifier `json:"parameter_name"`
    Type string `json:"parameter_type"`
}

type Component struct{
    Name *Identifier `json:"name"`
    Params []*Parameter `json:"params"`
    Body []*Statement `json:"body"`

}
type Contract struct{
    Name *Identifier `json:"name"`
    Params []*Parameter `json:"params"`
    Fields []*Field `json:"fields"`
    Components []*Component `json:"components"`
}
