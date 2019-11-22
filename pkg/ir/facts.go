package ir

import (
	//"errors"
	"fmt"
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
	bindFacts      []string
	condFacts      []string
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
			factString := fmt.Sprintf("plan_%d_%d\t%d\t%d", n.ID(), i, n.ID(), u.ID())
			fd.planFacts = append(fd.planFacts, factString)
		}

	case *Send:
		factString := fmt.Sprintf("%d\t%d", n.ID(), n.Data.ID())
		fd.sendFacts = append(fd.sendFacts, factString)
	case *Load:
		factString := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.loadFacts = append(fd.loadFacts, factString)
	case *Save:
		factString := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.saveFacts = append(fd.saveFacts, factString)
	default:
		fmt.Printf("+ %T %d\n", node, node.ID())
	}

	fd.visited[node] = true

	return fd
}

func DumpFacts(builder *CFGBuilder) {
	v := FactsDumper{
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
		bindFacts:      []string{},
		condFacts:      []string{},
		typeVarFacts:   []string{},
		enumFacts:      []string{},
		enumTypeFacts:  []string{},
		Facts:          []string{},
	}
	Walk(&v, builder.Transitions["test"], nil)

	//Printing
	fmt.Println()
	fmt.Println("Unit")
	for i, l := range v.unitFacts {
		fmt.Println(i, l)
	}

	fmt.Println()
	fmt.Println("Proc")
	for i, l := range v.procFacts {
		fmt.Println(i, l)
	}
	fmt.Println()
	fmt.Println("Plan")
	for i, l := range v.planFacts {
		fmt.Println(i, l)
	}
}
