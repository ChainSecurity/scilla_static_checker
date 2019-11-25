package ir

import (
	//"errors"
	"fmt"
	"os"
	"path"
	//"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	"strconv"
	//"strings"
)

type FactsDumper struct {
	visited        map[Node]bool
	idToPrefixID   map[uint64]string
	procFacts      []string
	unitFacts      []string
	planFacts      []string
	sendFacts      []string
	saveFacts      []string
	loadFacts      []string
	appDDFacts     []string
	appTDFacts     []string
	appTTFacts     []string
	argumentsFacts []string
	absDDFacts     []string
	msgFacts       []string
	dataFacts      []string
	strFacts       []string
	pickDataFacts  []string
	dataCaseFacts  []string
	natFacts       []string
	natTypeFacts   []string
	intFacts       []string
	intTypeFacts   []string
	rawFacts       []string
	rawTypeFacts   []string
	bindFacts      []string
	condFacts      []string
	condBindFacts  []string
	typeVarFacts   []string
	enumFacts      []string
	enumTypeFacts  []string
	Facts          []string
}

func (fd *FactsDumper) Visit(node Node, prev Node) Visitor {
	// DEBUG related check
	if node.ID() == 0 {
		panic(fmt.Sprintf("ID was not assigned %T", node))
	}

	visited, ok := fd.visited[node]
	if ok && visited {
		return nil
	}

	fd.unitFacts = append(fd.unitFacts, fmt.Sprintf("%d", node.ID()))

	switch n := node.(type) {
	case *Proc:
		//prefixID := fmt.Sprintf("Proc%d", n.ID())
		//fd.idToPrefixID[n.ID()] = prefixID
		fd.procFacts = append(fd.procFacts, strconv.FormatUint(n.ID(), 10))

		for i, u := range n.Plan {
			fact := fmt.Sprintf("plan_%d_%d\t%d\t%d", n.ID(), i, n.ID(), u.ID())
			fd.planFacts = append(fd.planFacts, fact)
		}

	case *Send:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Data.ID())
		fd.sendFacts = append(fd.sendFacts, fact)
	case *Save:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.saveFacts = append(fd.saveFacts, fact)
	case *Load:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.loadFacts = append(fd.loadFacts, fact)
	case *AppDD:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentsFacts = append(fd.argumentsFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appDDFacts = append(fd.appDDFacts, fact)
	case *AppTD:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentsFacts = append(fd.argumentsFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appTDFacts = append(fd.appTDFacts, fact)
	case *AppTT:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentsFacts = append(fd.argumentsFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appTTFacts = append(fd.appTTFacts, fact)
	case *AbsDD:
		for i, u := range n.Vars {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentsFacts = append(fd.argumentsFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Term.ID())
		fd.absDDFacts = append(fd.absDDFacts, fact)
	case *Msg:
		for k, v := range n.Data {
			fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), v.ID(), k)
			fd.dataFacts = append(fd.dataFacts, fact)
		}
		fact := fmt.Sprintf("%d", n.ID())
		fd.msgFacts = append(fd.msgFacts, fact)
	case *Str:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Data)
		fd.strFacts = append(fd.strFacts, fact)
	case *PickData:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.From)
		fd.pickDataFacts = append(fd.pickDataFacts, fact)
		for i, u := range n.With {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentsFacts = append(fd.argumentsFacts, argFact)
		}
	case *DataCase:

		body, ok := n.Body.(Node)
		if !ok {
			panic(fmt.Sprintf("Node is IDNode %T", body))
		}

		fact := fmt.Sprintf("%d\t%d\t%d\t%d", n.ID(), prev.ID(), n.Bind.ID(), body.ID())
		fd.dataCaseFacts = append(fd.dataCaseFacts, fact)
	case *Bind:
		fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), n.BindType.ID(), n.Cond.ID())
		fd.bindFacts = append(fd.bindFacts, fact)
	case *Cond:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Case)
		fd.condFacts = append(fd.condFacts, fact)

		for i, u := range n.Data {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.condBindFacts = append(fd.condBindFacts, fact)
		}
	case *Nat:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.NatType.ID(), n.Data)
		fd.natFacts = append(fd.natFacts, fact)
	case *NatType:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Size)
		fd.natTypeFacts = append(fd.natTypeFacts, fact)

	case *Int:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.IntType.ID(), n.Data)
		fd.intFacts = append(fd.intFacts, fact)
	case *IntType:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Size)
		fd.intTypeFacts = append(fd.intTypeFacts, fact)

	case *Raw:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.RawType.ID(), n.Data)
		fd.rawFacts = append(fd.rawFacts, fact)
	case *RawType:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Size)
		fd.rawTypeFacts = append(fd.rawTypeFacts, fact)

	default:
		fmt.Printf("+ %T %d\n", node, node.ID())
	}

	fd.visited[node] = true

	return fd
}

func DumpFacts(builder *CFGBuilder) {
	fd := FactsDumper{
		visited:        map[Node]bool{},
		idToPrefixID:   map[uint64]string{},
		procFacts:      []string{},
		unitFacts:      []string{},
		planFacts:      []string{},
		sendFacts:      []string{},
		saveFacts:      []string{},
		loadFacts:      []string{},
		appDDFacts:     []string{},
		appTDFacts:     []string{},
		appTTFacts:     []string{},
		argumentsFacts: []string{},
		absDDFacts:     []string{},
		msgFacts:       []string{},
		dataFacts:      []string{},
		strFacts:       []string{},
		pickDataFacts:  []string{},
		dataCaseFacts:  []string{},
		natFacts:       []string{},
		natTypeFacts:   []string{},
		intFacts:       []string{},
		intTypeFacts:   []string{},
		rawFacts:       []string{},
		rawTypeFacts:   []string{},
		bindFacts:      []string{},
		condFacts:      []string{},
		condBindFacts:  []string{},
		typeVarFacts:   []string{},
		enumFacts:      []string{},
		enumTypeFacts:  []string{},
		Facts:          []string{},
	}
	Walk(&fd, builder.Transitions["test"], nil)

	fileToFacts := map[string][]string{
		"proc":      fd.procFacts,
		"unit":      fd.unitFacts,
		"plan":      fd.planFacts,
		"send":      fd.sendFacts,
		"save":      fd.saveFacts,
		"load":      fd.loadFacts,
		"appDD":     fd.appDDFacts,
		"appTD":     fd.appTDFacts,
		"appTT":     fd.appTTFacts,
		"arguments": fd.argumentsFacts,
		"absDD":     fd.absDDFacts,
		"msg":       fd.msgFacts,
		"data":      fd.dataFacts,
		"str":       fd.strFacts,
		"pick_data": fd.pickDataFacts,
		"data_case": fd.dataCaseFacts,
		"nat":       fd.natFacts,
		"nat_type":  fd.natTypeFacts,
		"int":       fd.intFacts,
		"int_type":  fd.intTypeFacts,
		"raw":       fd.rawFacts,
		"raw_type":  fd.rawTypeFacts,
		"bind":      fd.bindFacts,
		"cond":      fd.condFacts,
		"cond_bind": fd.condBindFacts,
		"type":      fd.typeVarFacts,
		"enum":      fd.enumFacts,
		"enum_type": fd.enumTypeFacts,
	}

	outFolder := "out"
	err := os.Mkdir(outFolder, 0700)
	if err != nil {
		panic(err)
	}

	for fileName, lines := range fileToFacts {
		filePath := path.Join(outFolder, fmt.Sprintf("%s.tsv", fileName))
		f, err := os.Create(filePath)
		if err != nil {
			f.Close()
			panic(err)
		}
		for _, line := range lines {
			fmt.Fprintln(f, line)
			if err != nil {
				panic(err)
			}
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
