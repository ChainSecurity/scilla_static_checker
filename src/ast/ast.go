package ast

import (
)
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
    Key GenericLiteral `json:"key"`
    Value GenericLiteral `json:"value"`
}

type CtrDef struct {
    CrtDefName Identifier `json:"crt_def_name"`
    CArgTypes []string `json:"c_arg_types"`
}

type GenericLiteral struct{
    Value string `json:"value"`
    String string `json:"string"`
    KeyType string `json:"key_type"`
    ValueType string `json:"value_type"`
    MValues []MapValue `json:"mvalues"`
    NodeType string `json:"node_type"`
}

//type StringLiteral struct{
    //Value string `json:"value"`
//}

//type BNumLiteral struct{
    //Value string `json:"value"`

//}
//type ByStrLiteral struct{
    //Value string `json:"value"`
    //String string `json:"string"`

//}
//type ByStrXLiteral struct{
    //Value string `json:"value"`

//}
//type IntLiteral struct{
    //Value string `json:"value"`

//}
//type UintLiteral struct{
    //Value string `json:"value"`
//}

//type MapLiteral struct{
    //KeyType string `json:"key_type"`
    //ValueType string `json:"value_type"`
    //MValues []MapValue `json:"mvalues"`
//}

//type ADTValueLiteral struct {
//}


type AnnotatedNode struct{
    Loc Location `json:"loc"`
}

type GenericExpression struct{
    AnnotatedNode
    Value GenericLiteral `json:"value"`
    Variable Identifier `json:"variable"`
    VariableType string `json:"variable_type"` //Optional 
    Expr *GenericExpression `json:"expression"`
    Body *GenericExpression `json:"body"`
    MArgs []MessageArgument `json:"margs"`
    Lhs Identifier `json:"lhs"`
    RhsExpr *GenericExpression `json:"rhs_expr"`
    FunType string `json:"fun_type"`
    RhsList []Identifier `json:"rhs_list"`
    Types []string `json:"types"`
    ConstructorName string `json:"constructor_name"`
    Args []Identifier `json:"args"`
    Cases []MatchExpressionCase `json:"cases"`
    BuiltintFunction Builtin `json:"builtin_function"`
    NodeType string `json:"node_type"`
}

//type LiteralExpression struct{
    //AnnotatedNode
    //Value GenericLiteral
//}

//type VarExpression struct{
    //AnnotatedNode
    //Variable Identifier
//}

//type LetExpression struct{
    //AnnotatedNode
    //Variable Identifier
    //VariableType string //Optional
    //Expr GenericExpression
    //Body GenericExpression
//}

//type MessageExpression struct{
    //AnnotatedNode
    //Arguments []MessageArgument
//}

//type FunExpression struct{
    //AnnotatedNode
    //Lhs Identifier
    //RhsExpr GenericExpression
    //FunType string
//}

//type AppExpression struct{
    //AnnotatedNode
    //Lhs Identifier
    //Rhs []Identifier
//}

//type ConstrExpression struct{
    //AnnotatedNode
    //Types []string
    //ConstructorName string
    //Arguments []Identifier
//}

//type MatchExpression struct{
    //AnnotatedNode
    //Lhs Identifier
    //Rhs []MatchExpressionCase
//}

//type BuiltinExpression struct{
    //AnnotatedNode
    //Arguments []Identifier
    //BuiltintFunction Builtin
//}

//type TFunExpression struct{
    //AnnotatedNode
//}

//type TAppExpression struct{
    //AnnotatedNode
//}

//type FixpointExpression struct{
    //AnnotatedNode
//}

type GenericPayload struct{
    Lit GenericLiteral `json:"literal"`
    Value Identifier `json:"value"`
    NodeType string `json:"node_type"`
}

//type PayloadLitral struct{
    //Lit GenericLiteral
//}

//type PayloadVariable struct{
    //Value Identifier
//}


type MessageArgument struct{
    Var string `json:"variable"`
    // TODO fix name
    P GenericPayload `json:"payload"`
}

type GenericPattern struct{
    Variable Identifier `json:"variable"`
    ConstructorName string `json:"constructor_name"`
    Pats []GenericPattern `json:"patterns"`
    NodeType string `json:"node_type"`
}

//type WildcardPattern struct{
//}

//type BinderPattern struct{
    //Variable Identifier `json:"variable"`
//}

//type ConstructorPattern struct{
    //ConstructorName string `json:"constructor_name"`
    //Pats []GenericPattern `json:"patterns"`
//}

type MatchExpressionCase struct{
    Pat GenericPattern `json:"pattern"`
    Expr GenericExpression `json:"expression"`
}

type GenericStatement struct{
    AnnotatedNode
    Lhs Identifier `json:"lhs"`
    Rhs Identifier `json:"rhs"`
    RhsExpr GenericExpression `json:"rhs_expr"`
    RhsStr string `json:"rhs_str"`
    Name Identifier `json:"name"`
    Keys []Identifier `json:"keys"`
    IsValueRetrieve bool `json:"is_value_retrieve"`
    Arg Identifier `json:"arg"`
    Cases []MatchStatementCase `json:"cases"`
    Messages []Identifier `json:"messages"`
    NodeType string `json:"node_type"`
}

type MatchStatementCase struct{
    Pat GenericPattern `json:"pattern"`
    Body []GenericStatement `json:"pattern_body"`
}
type Builtin struct{
    Loc Location
    Type string
}

//type LoadStatement struct{
    //AnnotatedNode
    //Lhs Identifier
    //Rhs Identifier

//}
//type StoreStatement struct{
    //AnnotatedNode
    //Lhs Identifier
    //Rhs Identifier

//}
//type BindStatement struct{
    //AnnotatedNode
    //Lhs Identifier
    //RhsExpr GenericExpression

//}
//type MapUpdateStatement struct{
    //AnnotatedNode
    //Name Identifier
    //Rhs Identifier //Optional
    //Keys []Identifier

//}
//type MapGetStatement struct{
    //AnnotatedNode
    //Name Identifier
    //Lhs Identifier
    //Keys []Identifier
    //IsValueRetrieve bool
//}

//type MatchStatement struct{
    //AnnotatedNode
    //Arg Identifier
    //Cases []MatchStatementCase

//}
//type ReadFromBCStatement struct{
    //AnnotatedNode
    //Lhs Identifier
    //RhsStr string

//}
//type AcceptPaymentStatement struct{
    //AnnotatedNode
//}

//type SendMsgsStatement struct{
    //AnnotatedNode
    //Arg Identifier

//}
//type CreateEvntStatement struct{
    //AnnotatedNode
    //Arg Identifier
//}

//type CallProcStatement struct{
    //AnnotatedNode
    //Arg Identifier
    //Messages []Identifier

//}
//type ThrowStatement struct{
    //AnnotatedNode
    //Arg Identifier // Optional
//}

type GenericLibEntry struct {
    VariableType string `json:"variable_type"`
    Expr GenericExpression `json:"expression"`
    CtrDefs []CtrDef `json:"ctr_defs"`
    NodeType string `json:"node_type"`
}

type LibraryVariable struct{
    VariableType string `json:"variable_type"` // Optional
    Expr GenericExpression `json:"expression"`
}

type LibraryType struct{
    CtrDefs []CtrDef `json:"ctr_defs"`
}

type Library struct{
    Name Identifier  `json:"library_name"`
    Entries []GenericLibEntry `json:"library_entries"`
}

type ExternalLibrary struct{
    Name Identifier `json:"name"`
    Alias Identifier `json:"alias"` // Optional

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
    Name Identifier `json:"field_name"`
    Type string `json:"field_type"`
    Expr GenericExpression `json:"expression"`
}

type Parameter struct{
    Name Identifier `json:"parameter_name"`
    Type string `json:"parameter_type"`
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
