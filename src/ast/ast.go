package ast

import (
    "fmt"
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
    SourceFile string
    Line int
    Column int
}

type Identifier struct{
    Loc Location
    Id string
}

type MapValue struct{
    Key Literal
    Value Literal
}

type CtrDef struct {
    CrtDefName Identifier
    CArgTypes []string
}

type StringLiteral struct{
    Value string
}

type BNumLiteral struct{
    Value string

}
type ByStrLiteral struct{
    Value string
    String string

}
type ByStrXLiteral struct{
    Value string
    value string

}
type IntLiteral struct{
    Value string

}
type UintLiteral struct{
    Value string
    value string
}

type MapLiteral struct{
    KeyType string
    ValueType string
    NodeType string
    Value []MapValue
}

type ADTValueLiteral struct {
}

func (*StringLiteral) literalNode() {}
func (*BNumLiteral) literalNode() {}
func (*ByStrXLiteral) literalNode() {}
func (*IntLiteral) literalNode() {}
func (*UintLiteral) literalNode() {}
func (*MapLiteral) literalNode() {}
func (*ADTValueLiteral) literalNode() {}


type AnnotatedNode struct{
    Loc Location
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
    arguments []MessageArgument
}

type FunExpression struct{
    AnnotatedNode
    Lhs Identifier
    Rhs Expression
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

type PayloadLitral struct{
    Lit Literal
}

type PayloadVariable struct{
    Value Identifier
}


func (*PayloadLitral) payloadNode() {}
func (*PayloadVariable) payloadNode() {}

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

func (*WildcardPattern) patternNode() {}
func (*BinderPattern) patternNode() {}
func (*ConstructorPattern) patternNode() {}

type MatchExpressionCase struct{
    Pat Pattern
    Expr Expression
}

type MatchStatementCase struct{
    Pat Pattern
    PatternBody []Statement
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
    Rhs Expression

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
    Rhs string

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

type LibraryVariable struct{
    VariableType string // Optional
    Expr Expression
}

type LibraryType struct{
    CtrDefs []CtrDef
}

func (*LibraryVariable) libraryEntryNode() {}
func (*LibraryType) libraryEntryNode() {}

type Library struct{
    Name Identifier
    Entries []LibraryEntry
}

type ExternalLibrary struct{
    Name Identifier
    Alias Identifier // Optional

}
type ContractModule struct{
    ScillaMajorVersion int
    Name Identifier
    Library Library // Optional
    ExternalLibraries []ExternalLibrary
    //TODO name
    C Contract
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
    Name Identifier
    Params []Parameter
    Body []Statement

}
type Contract struct{
    Name Identifier
    Params []Parameter
    Fields []Field
    Components []Component
}
