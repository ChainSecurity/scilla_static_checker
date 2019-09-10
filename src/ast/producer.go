package ast

import (
    "fmt"
    "log"
    "encoding/json"
)

func Visit(m map[string]interface{}) (AstNode) {
    switch m["node_type"]{
    case "ContractModule":
        b,_ := json.Marshal(m)
        var c ContractModule
        if err := json.Unmarshal(b, &c); err != nil {
            panic(err)
        }
        fmt.Println(c.C.Components[0].Body[0].RhsStr)
    default:
        err := fmt.Errorf("Unknown node_type found in JSON: %s", m["node_type"])
        log.Fatal(err)
    }
    return AnnotatedNode{}
}
