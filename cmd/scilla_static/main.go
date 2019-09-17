package main

import (
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
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
	cm := ast.Parse_cmod(byteValue)
	ast.Inspect(cm, func(n ast.AstNode) bool {
		var s string
		switch x := n.(type) {
		case *ast.Identifier:
			s = x.Id
		default:
			s = fmt.Sprintf("inspecting %T", x)
		}
		if s != "" {
			fmt.Println(s)
		}
		return true
	})
	//fmt.Println(res["node_type"])

	//fmt.Println(res.Fruits[0])
}
