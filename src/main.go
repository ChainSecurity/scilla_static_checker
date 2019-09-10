package main

import (
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
    "./ast"
)

func main() {
    jsonPath := os.Args[1]
    jsonFile, err := os.Open(jsonPath)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)
    var res  map[string]interface{}

    json.Unmarshal([]byte(byteValue), &res)
    ast.Visit(res)
    //fmt.Println(res["node_type"])

    //fmt.Println(res.Fruits[0])
}
