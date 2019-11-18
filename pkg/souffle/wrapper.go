package main

import (
	"bytes"
	"fmt"
	"log"
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

func main() {
	runSouffle("analysis.dl", "facts_in", "facts_out")
}
