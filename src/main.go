package main

import (
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
    ast.Visit(byteValue)
    //fmt.Println(res["node_type"])

    //fmt.Println(res.Fruits[0])
}
