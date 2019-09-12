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

func LibEntryUnmarshal(rawMsg *json.RawMessage) (LibEntry, error) {
        ntype, err := getNodeType(rawMsg)
        if err != nil {
            return nil, err
        }
        var n map[string]string
        err = json.Unmarshal(*rawMsg, &n)
        fmt.Println(n)
        switch ntype{
        case "LibraryVariable":
            var m LibraryVariable
            err := json.Unmarshal(*rawMsg, &m)
            fmt.Println(m)
            return &m, err
        case "LibraryType":
            var m LibraryType
            err = json.Unmarshal(*rawMsg, &m)
            return &m, err
        default:
            return nil, errors.New("Unsupported type found!")
        }
}

func (l *LibraryVariable) UnmarshalJSON(b []byte) error {
    var objMap map[string]*json.RawMessage
    err := json.Unmarshal(b, &objMap)
    if err != nil {
        return err
    }

    //var rawMsg *json.RawMessage
    //err = json.Unmarshal(*objMap["expression"], &rawMsg)
    //if err != nil {
        //return err
    //}

    //e, err := ExpressionUnmarshal(rawMsg)
    //if err != nil {
        //return err
    //}
    //l.Expr = e
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
    return nil
}

func Parse_cmod(b []byte) {
    b = []byte(`{ "library_entries": [ { "node_type": "LibraryVariable", "variable_type": "asd"}  ] }`)
    var c Library
    if err := json.Unmarshal(b, &c); err != nil {
        panic(err)
    }
}
