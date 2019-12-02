package main

import (
	"bufio"
	"fmt"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ast"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ir"
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
	cm, err := ast.Parse_mod(byteValue)
	if err != nil {
		panic(err)
	}

	dotPath := os.Args[2]
	b := ir.BuildCFG(cm)

	dot := ir.GetDot(b)
	f, err := os.Create(dotPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprint(w, dot)
	w.Flush()
}
