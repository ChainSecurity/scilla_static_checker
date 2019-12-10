package souffle

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

var SOUFFLE = "souffle"

func ReadOutput(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()
	return txtlines, nil
}

func RunSouffle(datalog string, factsIn string, factsOut string) {
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

func MakeCleanFolder(path string) error {
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
