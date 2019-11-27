package main

import (
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ir"
	"io/ioutil"
	"os"
	"path"
)

func makeCleanFolder(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	err := os.Mkdir(path, 0700)
	if err != nil {
		return err
	}
	return nil
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

	analysisFolder := "./souffle_analysis"
	factsInFolder := path.Join(analysisFolder, "facts_int")

	err = makeCleanFolder(factsInFolder)
	if err != nil {
		panic(err)
	}

	b := ir.BuildCFG(cm)
	ir.DumpFacts(b, factsInFolder)

}
