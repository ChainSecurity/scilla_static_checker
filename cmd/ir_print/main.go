package main

import (
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ir"
	"io/ioutil"
	"os"
)

type TestVisitor struct {
	visited map[ir.IRNode]bool
}

func (t TestVisitor) Visit(node ir.IRNode) ir.Visitor {
	//fmt.Printf("+ %T\n", node)
	visited, ok := t.visited[node]
	if ok && visited {
		return nil
	}

	t.visited[node] = true
	return t
}

func main() {
	jsonPath := os.Args[1]
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	cm, err := ast.Parse_mod(byteValue)
	if err != nil {
		panic(err)
	}

	b := ir.BuildCFG(cm)

	v := TestVisitor{map[ir.IRNode]bool{}}
	ir.Walk(v, b.Transitions["test"], nil)

}
