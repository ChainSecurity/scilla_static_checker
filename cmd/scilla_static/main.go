package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ir"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var SOUFFLE = "souffle"

func runSouffle(datalog, factsIn, factsOut string) {
	cmd := exec.Command(SOUFFLE, "--fact-dir", factsIn, "--output-dir", factsOut, datalog)

	var out bytes.Buffer
	cmd.Stdout = &out

	var err = cmd.Run()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf(out.String())
	}

}

func readOutput(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()
	return txtlines

}

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
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	cm, err := ast.Parse_mod(byteValue)
	if err != nil {
		panic(err)
	}

	b := ir.BuildCFG(cm)

	analysisFolder := "./souffle_analysis"
	factsOutFolder := path.Join(analysisFolder, "facts_out")
	factsInFolder := path.Join(analysisFolder, "facts_in")

	err = makeCleanFolder(factsOutFolder)
	if err != nil {
		panic(err)
	}

	err = makeCleanFolder(factsInFolder)
	if err != nil {
		panic(err)
	}

	ir.DumpFacts(b, factsInFolder)

	runSouffle("souffle_analysis/analysis.dl", "souffle_analysis/facts_in", "souffle_analysis/facts_out")
	result := readOutput("souffle_analysis/facts_out/vulnerability.csv")
	fmt.Println("======RESULTS======")
	fmt.Println(result)
}
