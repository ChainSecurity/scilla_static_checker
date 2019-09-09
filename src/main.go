package main

import (
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
)
type response1 struct {
    Page   int
    Fruits []string
}

type response2 struct {
    Page   int      `json:"page"`
    Fruits []string `json:"fruits"`
}

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
    fmt.Println(res["data"])

    //fmt.Println(res.Fruits[0])
}
