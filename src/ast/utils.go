package ast

import (
    "encoding/json"
)

func Parse_cmod(b []byte) *ContractModule{
    var c ContractModule
    if err := json.Unmarshal(b, &c); err != nil {
        panic(err)
    }
    return &c
}
