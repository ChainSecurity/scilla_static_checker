package ast

import (
)

type AstNode interface {
}

type Expression interface {
    AstNode
    expressionNode()
}

type Statement interface {
    AstNode
    statementNode()
}

type Literal interface {
    AstNode
    literalNode()
}

type Payload interface {
    AstNode
    payloadNode()
}

type Pattern interface {
    AstNode
    patternNode()
}

type LibraryEntry interface {
    AstNode
    libraryEntryNode()
}

type Location struct{
    SourceFile string `json:"source_file"`
    Line int `json:"line"`
    Column int `json:"column"`
}

type Identifier struct{
    Loc Location `json:"loc"`
    Id string `json:"identifier"`
}

type MapValue struct{
    Key GeneralLiteral
    Value GeneralLiteral
}

type CtrDef struct {
    CrtDefName Identifier `json:"crt_def_name"`
    CArgTypes []string `json:"c_arg_types"`
}

type GeneralLiteral struct{
    Value string `json:"value"`
    String string `json:"string"`
    KeyType string `json:"key_type"`
    ValueType string `json:"value_type"`
    MValues []MapValue `json:"mvalues"`
}

type StringLiteral struct{
    Value string `json:"value"`
}

type BNumLiteral struct{
    Value string `json:"value"`

}
type ByStrLiteral struct{
    Value string `json:"value"`
    String string `json:"string"`

}
type ByStrXLiteral struct{
    Value string `json:"value"`

}
type IntLiteral struct{
    Value string `json:"value"`

}
type UintLiteral struct{
    Value string `json:"value"`
}

type MapLiteral struct{
    KeyType string `json:"key_type"`
    ValueType string `json:"value_type"`
    MValues []MapValue `json:"mvalues"`
}

type ADTValueLiteral struct {
}


type AnnotatedNode struct{
    Loc Location `json:"loc"`
}

type GeneralExpression struct{
    AnnotatedNode
    Value GeneralLiteral `json:"value"`
    Variable Identifier `json:"variable"`
    VariableType string `json:"variable_type"` //Optional 
    Expr *GeneralExpression `json:"expression"`
    Body *GeneralExpression `json:"body"`
    MArgs []MessageArgument `json:"margs"`
    Lhs Identifier `json:"lhs"`
    RhsExpr *GeneralExpression `json:"rhs_expr"`
    FunType string `json:"fun_type"`
    RhsList []Identifier `json:"rhs_list"`
    Types []string `json:"types"`
    ConstructorName string `json:"constructor_name"`
    Args []Identifier `json:"args"`
    Cases []MatchExpressionCase `json:"cases"`
    BuiltintFunction Builtin `json:"builtin_function"`
    NodeType string `json:"node_type"`
}

type LiteralExpression struct{
    AnnotatedNode
    Value Literal
}

type VarExpression struct{
    AnnotatedNode
    Variable Identifier
}

type LetExpression struct{
    AnnotatedNode
    Variable Identifier
    VariableType string //Optional
    Expr Expression
    Body Expression
}

type MessageExpression struct{
    AnnotatedNode
    Arguments []MessageArgument
}

type FunExpression struct{
    AnnotatedNode
    Lhs Identifier
    RhsExpr Expression
    FunType string
}

type AppExpression struct{
    AnnotatedNode
    Lhs Identifier
    Rhs []Identifier
}

type ConstrExpression struct{
    AnnotatedNode
    Types []string
    ConstructorName string
    Arguments []Identifier
}

type MatchExpression struct{
    AnnotatedNode
    Lhs Identifier
    Rhs []MatchExpressionCase
}

type BuiltinExpression struct{
    AnnotatedNode
    Arguments []Identifier
    BuiltintFunction Builtin
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


type PayloadLitral struct{
    Lit Literal
}

type PayloadVariable struct{
    Value Identifier
}


type MessageArgument struct{
    Var string
    // TODO fix name
    P Payload
}


type WildcardPattern struct{
}

type BinderPattern struct{
    variable Identifier
}

type ConstructorPattern struct{
    ConstructorName string
    Pats []Pattern
}

type MatchExpressionCase struct{
    Pat Pattern
    Expr Expression
}

type GenericStatement struct{
    AnnotatedNode
    Lhs Identifier `json:"lhs"`
    Rhs Identifier `json:"rhs"`
    RhsExpr GeneralExpression `json:"rhs_expr"`
    Name Identifier `json:"name"`
    Keys []Identifier `json:"keys"`
    IsValueRetrieve bool `json:"is_value_retrieve"`
    Arg Identifier `json:"arg"`
    Cases []MatchStatementCase `json:"cases"`
    RhsStr string `json:"rhs_str"`
    Messages []Identifier `json:"messages"`
    NodeType string `json:"node_type"`
}

type MatchStatementCase struct{
    Pat Pattern
    PatternBody []GenericStatement
}
type Builtin struct{
    Loc Location
    Type string
}

type LoadStatement struct{
    AnnotatedNode
    Lhs Identifier
    Rhs Identifier

}
type StoreStatement struct{
    AnnotatedNode
    Lhs Identifier
    Rhs Identifier

}
type BindStatement struct{
    AnnotatedNode
    Lhs Identifier
    RhsExpr Expression

}
type MapUpdateStatement struct{
    AnnotatedNode
    Name Identifier
    Rhs Identifier //Optional
    Keys []Identifier

}
type MapGetStatement struct{
    AnnotatedNode
    Name Identifier
    Lhs Identifier
    Keys []Identifier
    IsValueRetrieve bool
}

type MatchStatement struct{
    AnnotatedNode
    Arg Identifier
    Cases []MatchStatementCase

}
type ReadFromBCStatement struct{
    AnnotatedNode
    Lhs Identifier
    RhsStr string

}
type AcceptPaymentStatement struct{
    AnnotatedNode
}

type SendMsgsStatement struct{
    AnnotatedNode
    Arg Identifier

}
type CreateEvntStatement struct{
    AnnotatedNode
    Arg Identifier
}

type CallProcStatement struct{
    AnnotatedNode
    Arg Identifier
    Messages []Identifier

}
type ThrowStatement struct{
    AnnotatedNode
    Arg Identifier // Optional
}


type LibraryVariable struct{
    VariableType string // Optional
    Expr Expression
}

type LibraryType struct{
    CtrDefs []CtrDef
}

type Library struct{
    Name Identifier
    Entries []LibraryEntry
}

type ExternalLibrary struct{
    Name Identifier
    Alias Identifier // Optional

}
type ContractModule struct{
    ScillaMajorVersion int `json:"scilla_major_version"`
    Name Identifier `json:"name"`
    Library Library `json:"library"`
    ExternalLibraries []ExternalLibrary `json:"external_libraries"`
    //TODO name
    C Contract `json:"contract"`
}

type Field struct{
    Name Identifier
    Type string
    Expr Expression
}

type Parameter struct{
    Name Identifier
    Type string
}

type Component struct{
    Name Identifier `json:"name"`
    Params []Parameter `json:"params"`
    Body []GenericStatement `json:"body"`

}
type Contract struct{
    Name Identifier `json:"name"`
    Params []Parameter `json:"params"`
    Fields []Field `json:"fields"`
    Components []Component `json:"components"`
}

func (*GeneralLiteral) literalNode() {}
func (*StringLiteral) literalNode() {}
func (*BNumLiteral) literalNode() {}
func (*ByStrXLiteral) literalNode() {}
func (*IntLiteral) literalNode() {}
func (*UintLiteral) literalNode() {}
func (*MapLiteral) literalNode() {}
func (*ADTValueLiteral) literalNode() {}

func (*GeneralExpression) expressionNode() {}
func (*LiteralExpression) expressionNode() {}
func (*VarExpression) expressionNode() {}
func (*LetExpression) expressionNode() {}
func (*MessageExpression) expressionNode() {}
func (*FunExpression) expressionNode() {}
func (*AppExpression) expressionNode() {}
func (*ConstrExpression) expressionNode() {}
func (*MatchExpression) expressionNode() {}
func (*BuiltinExpression) expressionNode() {}
func (*TFunExpression) expressionNode() {}
func (*TAppExpression) expressionNode() {}
func (*FixpointExpression) expressionNode() {}

func (*PayloadLitral) payloadNode() {}
func (*PayloadVariable) payloadNode() {}

func (*WildcardPattern) patternNode() {}
func (*BinderPattern) patternNode() {}
func (*ConstructorPattern) patternNode() {}

func (*GenericStatement) statementNode() {}
func (*LoadStatement) statementNode() {}
func (*StoreStatement) statementNode() {}
func (*BindStatement) statementNode() {}
func (*MapUpdateStatement) statementNode() {}
func (*MapGetStatement) statementNode() {}
func (*MatchStatement) statementNode() {}
func (*ReadFromBCStatement) statementNode() {}
func (*AcceptPaymentStatement) statementNode() {}
func (*SendMsgsStatement) statementNode() {}
func (*CreateEvntStatement) statementNode() {}
func (*CallProcStatement) statementNode() {}
func (*ThrowStatement) statementNode() {}

func (*LibraryVariable) libraryEntryNode() {}
func (*LibraryType) libraryEntryNode() {}
