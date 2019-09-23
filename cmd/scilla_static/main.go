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
	fmt.Println("Parsing")
	cm, err := ast.Parse_mod(byteValue)
	fmt.Println("Finished parsing")
	ast.Inspect(cm, func(n ast.AstNode) bool {
		if e, ok := n.(*ast.LibEntry); ok {
			d := *e
			if lt, ok := d.(*ast.LibraryType); ok {
				fmt.Printf("%T %s \n", lt, lt.Name.Id)
				for _, c := range lt.CtrDefs {
					fmt.Printf("\t %s %s\n", c.CDName.Id, c.CArgTypes)
				}
			}
		}
		return true
	})
	//fmt.Println(res["node_type"])

	//fmt.Println(res.Fruits[0])
}
