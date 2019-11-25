package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var SOUFFLE = "souffle"

func runSouffle(datalog, factsIn, factsOut string) {
	cmd := exec.Command(SOUFFLE, "--fact-dir", factsIn, "--output-dir", factsOut, datalog)

	var out bytes.Buffer
	cmd.Stdout = &out

	var err = cmd.Run()
	if err != nil {
		fmt.Println("Error happened")
		log.Fatal(err.Error())
		fmt.Printf(out.String())
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

// We run the analysis in the root directory
func main() {
	runSouffle("souffle_analysis/analysis.dl", "souffle_analysis/facts_in", "souffle_analysis/facts_out")
	result := readOutput("souffle_analysis/facts_out/vulnerability.csv")
	fmt.Println(result)
}
