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
		log.Fatal(err.Error())
		fmt.Printf(out.String())
	} else {
		fmt.Printf(out.String())
	}

}

func readOutput(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
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

func main() {
	runSouffle("analysis.dl", "facts_in", "facts_out")
	result := readOutput("facts_out/b.csv")
	fmt.Println(result)
}
