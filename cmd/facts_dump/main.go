package main

import (
	"fmt"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ast"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ir"
	"github.com/ChainSecurity/scilla_static_checker/pkg/souffle"
	"io/ioutil"
	"os"
	"path"
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

	analysisFolder := "./souffle_analysis"
	factsInFolder := path.Join(analysisFolder, "facts_in")

	err = souffle.MakeCleanFolder(factsInFolder)
	if err != nil {
		panic(err)
	}

	b := ir.BuildCFG(cm)
	ir.DumpFacts(b, factsInFolder)

}
