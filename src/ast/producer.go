package ast

import (
    "fmt"
    "encoding/json"
)

func Visit(b []byte) {
    var c ContractModule
    if err := json.Unmarshal(b, &c); err != nil {
        panic(err)
    }
    if b, err := json.Marshal(c); err != nil {
        panic(err)
    } else {
        fmt.Println(string(b))
    }
}
