package ast

import (
    "encoding/json"
    "errors"
    "fmt"
)


func getNodeType(rawMsg *json.RawMessage) (string, error){
    type ntype struct {
        NodeType string `json:"node_type"`
    }
    var n ntype
    err := json.Unmarshal(*rawMsg, &n)
    if err != nil {
        return "", err
    }
    return n.NodeType, nil
}

func ExpressionUnmarshal(rawMsg *json.RawMessage) (Expression, error) {
        ntype, err := getNodeType(rawMsg)
        if err != nil {
            return nil, err
        }

        fmt.Println(ntype)
        switch ntype{
        case "LiteralExpression":
            var m LiteralExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "LetExpression":
            var m LetExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "MessageExpression":
            var m MessageExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "FunExpression":
            var m FunExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "AppExpression":
            var m AppExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "ConstrExpression":
            var m ConstrExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "MatchExpression":
            var m MatchExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "BuiltinExpression":
            var m BuiltinExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "TFunExpression":
            var m TFunExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "TAppExpression":
            var m TAppExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "FixpointExpression":
            var m FixpointExpression
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        default:
            return nil, errors.New("Unsupported type found!")
        }
}

func LiteralUnmarshal(rawMsg *json.RawMessage) (Literal, error) {
        ntype, err := getNodeType(rawMsg)
        if err != nil {
            return nil, err
        }
        var n map[string]string
        err = json.Unmarshal(*rawMsg, &n)
        fmt.Println(ntype)
        switch ntype{
        case "StringLiteral":
            var m StringLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "BNumLiteral":
            var m BNumLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "ByStrLiteral":
            var m ByStrLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "ByStrXLiteral":
            var m ByStrXLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "IntLiteral":
            var m IntLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "UintLiteral":
            var m UintLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "MapLiteral":
            var m MapLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "ADTValLiteral":
            var m ADTValLiteral
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        default:
            return nil, errors.New("Unsupported type found!")
        }
}

func LibEntryUnmarshal(rawMsg *json.RawMessage) (LibEntry, error) {
        ntype, err := getNodeType(rawMsg)
        if err != nil {
            return nil, err
        }
        var n map[string]string
        err = json.Unmarshal(*rawMsg, &n)
        fmt.Println(ntype)
        switch ntype{
        case "LibraryVariable":
            var m LibraryVariable
            err := json.Unmarshal(*rawMsg, &m)
            return &m, err
        case "LibraryType":
            var m LibraryType
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        default:
            return nil, errors.New("Unsupported type found!")
        }
}

func (le *LetExpression) UnmarshalJSON(b []byte) error {
    var objMap map[string]*json.RawMessage
    err := json.Unmarshal(b, &objMap)
    if err != nil {
        return err
    }

    var rawMsg *json.RawMessage
    rawMsg = objMap["expression"]
    e, err := ExpressionUnmarshal(rawMsg)
    if err != nil {
        return err
    }

    rawMsg = objMap["body"]
    bd, err := ExpressionUnmarshal(rawMsg)
    if err != nil {
        return err
    }

    type core struct{
        AnnotatedNode
        Var *Identifier `json:"variable"`
        VarType string `json:"variable_type"` //Optional 
    }

    var c core
    err = json.Unmarshal(b, &c)
    if err != nil {
        return err
    }

    le.Expr = &e
    le.Body = &bd
    le.AnnotatedNode = c.AnnotatedNode
    le.Var = c.Var
    le.VarType = c.VarType
    return nil
}

func (le *LiteralExpression) UnmarshalJSON(b []byte) error {
    var objMap map[string]*json.RawMessage
    err := json.Unmarshal(b, &objMap)
    if err != nil {
        return err
    }

    var rawMsg *json.RawMessage
    rawMsg = objMap["value"]
    v, err := LiteralUnmarshal(rawMsg)
    if err != nil {
        return err
    }

    type core struct{
        AnnotatedNode
    }

    var c core
    err = json.Unmarshal(b, &c)
    if err != nil {
        return err
    }

    le.AnnotatedNode = c.AnnotatedNode
    le.Val = &v
    return nil
}

func (fe *FunExpression) UnmarshalJSON(b []byte) error {
    var objMap map[string]*json.RawMessage
    err := json.Unmarshal(b, &objMap)
    if err != nil {
        return err
    }

    var rawMsg *json.RawMessage
    rawMsg = objMap["rhs_expr"]
    e, err := ExpressionUnmarshal(rawMsg)
    if err != nil {
        return err
    }


    type core struct{
        AnnotatedNode
        FunType string `json:"fun_type"`
        Lhs *Identifier  `json:"lhs"`
    }

    var c core
    err = json.Unmarshal(b, &c)
    if err != nil {
        return err
    }

    fe.RhsExpr = &e
    fe.AnnotatedNode = c.AnnotatedNode
    fe.Lhs = c.Lhs
    fe.FunType = c.FunType
    return nil
}

func (l *LibraryVariable) UnmarshalJSON(b []byte) error {
    var objMap map[string]*json.RawMessage
    err := json.Unmarshal(b, &objMap)
    if err != nil {
        return err
    }

    var rawMsg *json.RawMessage
    rawMsg = objMap["expression"]
    e, err := ExpressionUnmarshal(rawMsg)
    if err != nil {
        return err
    }

    l.Expr = &e

    type core struct{
        VarType string `json:"variable_type"`
    }

    var c core
    err = json.Unmarshal(b, &c)
    if err != nil {
        return err
    }

    l.VarType = c.VarType
    return nil
}

func (l *Library) UnmarshalJSON(b []byte) error {
    var objMap map[string]*json.RawMessage
    err := json.Unmarshal(b, &objMap)
    if err != nil {
        return err
    }

    var rawMsgs []*json.RawMessage
    err = json.Unmarshal(*objMap["library_entries"], &rawMsgs)
    if err != nil {
        return err
    }

    l.Entries = make([]LibEntry, len(rawMsgs))
    for index, rawMsg := range rawMsgs {
        e, err := LibEntryUnmarshal(rawMsg)
        if err != nil {
            return err
        }
        l.Entries[index] = e
    }

    type core struct{
        Name *Identifier  `json:"library_name"`
    }

    var c core
    err = json.Unmarshal(b, &c)
    if err != nil {
        return err
    }

    l.Name = c.Name
    return nil
}

func Parse_cmod(b []byte) {
    var c ContractModule
    if err := json.Unmarshal(b, &c); err != nil {
        fmt.Println(err)
        panic(err)
    }
    fmt.Println(c)
}
