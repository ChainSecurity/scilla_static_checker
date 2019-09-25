package main

import (
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ir"
	"io/ioutil"
	"os"
)

func main() {
	jsonPath := os.Args[1]
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println("Parsing")
	cm, err := ast.Parse_mod(byteValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished parsing")
	fmt.Println(cm == nil)
	ir.BuildCFG(cm)
	//fmt.Println(res["node_type"])

	//fmt.Println(res.Fruits[0])
}
