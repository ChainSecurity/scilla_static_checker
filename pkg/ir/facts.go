package ir

import (
	//"errors"
	"fmt"
	"os"
	"path"
	//"github.com/ChainSecurity/scilla_static_checker/pkg/ast"
	//"strings"
)

type FactsDumper struct {
	visited          map[Node]bool
	idToPrefixID     map[uint64]string
	absDDFacts       []string
	absTDFacts       []string
	absTTFacts       []string
	acceptFacts      []string
	appDDFacts       []string
	appTDFacts       []string
	appTTFacts       []string
	argumentFacts    []string
	bindFacts        []string
	bnrFacts         []string
	bnrTypeFacts     []string
	builtinFacts     []string
	builtinVarFacts  []string
	callProcFacts    []string
	condBindFacts    []string
	condFacts        []string
	constructorFacts []string
	dataCaseFacts    []string
	dataVarFacts     []string
	enumFacts        []string
	enumTypeFacts    []string
	eventFacts       []string
	fieldFacts       []string
	intFacts         []string
	intTypeFacts     []string
	jumpFacts        []string
	keyArgumentFacts []string
	loadFacts        []string
	mapTypeFacts     []string
	msgFacts         []string
	msgTypeFacts     []string
	natFacts         []string
	natTypeFacts     []string
	pickDataFacts    []string
	pickProcFacts    []string
	planFacts        []string
	procCaseFacts    []string
	procFacts        []string
	procTypeFacts    []string
	procedureFacts   []string
	rawFacts         []string
	rawTypeFacts     []string
	saveFacts        []string
	sendFacts        []string
	setKindFacts     []string
	strFacts         []string
	strTypeFacts     []string
	transitionFacts  []string
	typeVarFacts     []string
	unitFacts        []string
}

func (fd *FactsDumper) Visit(node Node, prev Node) Visitor {

	visited, ok := fd.visited[node]
	if ok && visited {
		return nil
	}

	// DEBUG related check
	if node.ID() == 0 {
		panic(fmt.Sprintf("ID was not assigned %T", node))
	}

	fd.unitFacts = append(fd.unitFacts, fmt.Sprintf("%d", node.ID()))

	switch n := node.(type) {
	case *ProcType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.procTypeFacts = append(fd.procTypeFacts, fact)
	case *Proc:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.ProcName)
		fd.procFacts = append(fd.procFacts, fact)

		for i, v := range n.Vars {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), v.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, fact)
		}
		for i, u := range n.Plan {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.planFacts = append(fd.planFacts, fact)
		}

		if n.Jump != nil {
			fact = fmt.Sprintf("%d\t%d", n.ID(), n.Jump.ID())
			fd.jumpFacts = append(fd.jumpFacts, fact)
		}

	case *SetKind:
		fact := fmt.Sprintf("%d", n.ID())
		fd.setKindFacts = append(fd.setKindFacts, fact)
	case *TypeVar:
		fact := fmt.Sprintf("%d", n.ID())
		fd.typeVarFacts = append(fd.typeVarFacts, fact)
	case *DataVar:
		fact := fmt.Sprintf("%d", n.ID())
		fd.dataVarFacts = append(fd.dataVarFacts, fact)
	case *Event:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Data.ID())
		fd.eventFacts = append(fd.eventFacts, fact)
	case *Send:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Data.ID())
		fd.sendFacts = append(fd.sendFacts, fact)
	case *Accept:
		fact := fmt.Sprintf("%d", n.ID())
		fd.acceptFacts = append(fd.acceptFacts, fact)
	case *Save:
		fact := fmt.Sprintf("%d\t%s\t%d", n.ID(), n.Slot, n.Data.ID())
		fd.saveFacts = append(fd.saveFacts, fact)
		for i, p := range n.Path {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), p.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
	case *Load:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.loadFacts = append(fd.loadFacts, fact)
		for i, p := range n.Path {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), p.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}

	case *AbsDD:
		for i, u := range n.Vars {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Term.ID())
		fd.absDDFacts = append(fd.absDDFacts, fact)
	case *AbsTD:
		for i, u := range n.Vars {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Term.ID())
		fd.absTDFacts = append(fd.absTDFacts, fact)
	case *AbsTT:
		for i, u := range n.Vars {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Term.ID())
		fd.absTTFacts = append(fd.absTTFacts, fact)
	case *AppDD:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appDDFacts = append(fd.appDDFacts, fact)
	case *AppTD:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appTDFacts = append(fd.appTDFacts, fact)
	case *AppTT:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appTTFacts = append(fd.appTTFacts, fact)
	case *Msg:
		for k, v := range n.Data {
			fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), v.ID(), k)
			fd.keyArgumentFacts = append(fd.keyArgumentFacts, fact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.MsgType.ID())
		fd.msgFacts = append(fd.msgFacts, fact)
	case *MsgType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.msgTypeFacts = append(fd.msgTypeFacts, fact)
	case *Str:
		strData := "-1"
		if n.Data != "" {
			strData = n.Data
		}
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.StrType, strData)
		fd.strFacts = append(fd.strFacts, fact)
	case *StrType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.strTypeFacts = append(fd.strTypeFacts, fact)
	case *PickData:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.From)
		fd.pickDataFacts = append(fd.pickDataFacts, fact)
		for i, u := range n.With {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}

	case *DataCase:
		fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), n.Bind.ID(), n.Body.ID())
		fd.dataCaseFacts = append(fd.dataCaseFacts, fact)
	case *Bind:
		var condID int64
		condID = -1
		if n.Cond != nil {
			condID = n.Cond.ID()
		}
		fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), n.BindType.ID(), condID)
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
	case *Enum:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.EnumType.ID(), n.Case)
		fd.enumFacts = append(fd.enumFacts, fact)
		for i, u := range n.Data {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, fact)
		}
	case *CallProc:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.callProcFacts = append(fd.callProcFacts, fact)
		for i, u := range n.Args {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, fact)
		}
	case *PickProc:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.From.ID())
		fd.pickProcFacts = append(fd.pickProcFacts, fact)
		for i, u := range n.With {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, fact)
		}
	case *ProcCase:
		fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), n.Bind.ID(), n.Body.ID())
		fd.procCaseFacts = append(fd.procCaseFacts, fact)
	case *Builtin:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.BuiltinType.ID(), n.Label)
		fd.builtinFacts = append(fd.builtinFacts, fact)
	case *BuiltinVar:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.BuiltinVarType.ID(), n.Label)
		fd.builtinVarFacts = append(fd.builtinVarFacts, fact)
	case *Bnr:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.BnrType.ID(), n.Data)
		fd.bnrFacts = append(fd.bnrFacts, fact)
	case *BnrType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.bnrTypeFacts = append(fd.bnrTypeFacts, fact)
	case *EnumType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.enumTypeFacts = append(fd.enumTypeFacts, fact)
		for k, types := range n.Constructors {
			typeListID := fmt.Sprintf("EnumTypeConstructor_%d_%s", n.ID(), k)
			fact := fmt.Sprintf("%d\t%s\t%s", n.ID(), typeListID, k)
			fd.keyArgumentFacts = append(fd.keyArgumentFacts, fact)
			for j, t := range types {
				listFact := fmt.Sprintf("%s\t%d\t%d", typeListID, t.ID(), j)
				fd.argumentFacts = append(fd.argumentFacts, listFact)
			}

		}
	case *MapType:
		fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), n.KeyType.ID(), n.ValType.ID())
		fd.mapTypeFacts = append(fd.mapTypeFacts, fact)
	default:
		fmt.Printf("+ %T %d\n", node, node.ID())
	}

	fd.visited[node] = true

	return fd
}

func DumpFacts(builder *CFGBuilder, factsInFolder string) {
	fd := FactsDumper{
		visited:          map[Node]bool{},
		idToPrefixID:     map[uint64]string{},
		absDDFacts:       []string{},
		absTDFacts:       []string{},
		absTTFacts:       []string{},
		acceptFacts:      []string{},
		appDDFacts:       []string{},
		appTDFacts:       []string{},
		appTTFacts:       []string{},
		argumentFacts:    []string{},
		bindFacts:        []string{},
		bnrFacts:         []string{},
		bnrTypeFacts:     []string{},
		builtinFacts:     []string{},
		builtinVarFacts:  []string{},
		callProcFacts:    []string{},
		condBindFacts:    []string{},
		condFacts:        []string{},
		constructorFacts: []string{},
		dataCaseFacts:    []string{},
		dataVarFacts:     []string{},
		enumFacts:        []string{},
		enumTypeFacts:    []string{},
		eventFacts:       []string{},
		fieldFacts:       []string{},
		intFacts:         []string{},
		intTypeFacts:     []string{},
		keyArgumentFacts: []string{},
		loadFacts:        []string{},
		msgFacts:         []string{},
		natFacts:         []string{},
		natTypeFacts:     []string{},
		pickDataFacts:    []string{},
		pickProcFacts:    []string{},
		planFacts:        []string{},
		procCaseFacts:    []string{},
		procFacts:        []string{},
		procTypeFacts:    []string{},
		rawFacts:         []string{},
		rawTypeFacts:     []string{},
		saveFacts:        []string{},
		sendFacts:        []string{},
		setKindFacts:     []string{},
		strFacts:         []string{},
		strTypeFacts:     []string{},
		typeVarFacts:     []string{},
		unitFacts:        []string{},
	}

	fd.constructorFacts = append(fd.constructorFacts, fmt.Sprintf("%d", builder.Constructor.ID()))
	Walk(&fd, builder.Constructor, nil)

	for fname, _ := range builder.fieldTypeMap {
		fmt.Println("field", fname)
		fact := fmt.Sprintf("%s", fname)
		fd.fieldFacts = append(fd.fieldFacts, fact)
	}
	for tname, t := range builder.Transitions {
		fmt.Println("transition", tname)
		fact := fmt.Sprintf("%d", t.ID())
		fd.transitionFacts = append(fd.transitionFacts, fact)
		Walk(&fd, t, nil)
	}
	for pName, p := range builder.Procedures {
		fmt.Println("Procedure", pName)
		fact := fmt.Sprintf("%d", p.ID())
		fd.procedureFacts = append(fd.procedureFacts, fact)
		Walk(&fd, p, nil)
	}

	fileToFacts := map[string][]string{
		"absDD":       fd.absDDFacts,
		"absTD":       fd.absTDFacts,
		"absTT":       fd.absTTFacts,
		"accept":      fd.acceptFacts,
		"appDD":       fd.appDDFacts,
		"appTD":       fd.appTDFacts,
		"appTT":       fd.appTTFacts,
		"argument":    fd.argumentFacts,
		"bind":        fd.bindFacts,
		"bnr":         fd.bnrFacts,
		"bnrType":     fd.bnrTypeFacts,
		"builtin":     fd.builtinFacts,
		"builtinVar":  fd.builtinVarFacts,
		"callProc":    fd.callProcFacts,
		"cond":        fd.condFacts,
		"constructor": fd.constructorFacts,
		"condBind":    fd.condBindFacts,
		"dataCase":    fd.dataCaseFacts,
		"dataVar":     fd.dataVarFacts,
		"enum":        fd.enumFacts,
		"enumType":    fd.enumTypeFacts,
		"event":       fd.eventFacts,
		"field":       fd.fieldFacts,
		"int":         fd.intFacts,
		"intType":     fd.intTypeFacts,
		"jump":        fd.jumpFacts,
		"keyArgument": fd.keyArgumentFacts,
		"load":        fd.loadFacts,
		"mapType":     fd.mapTypeFacts,
		"msg":         fd.msgFacts,
		"nat":         fd.natFacts,
		"natType":     fd.natTypeFacts,
		"pickData":    fd.pickDataFacts,
		"pickProc":    fd.pickProcFacts,
		"procCase":    fd.procCaseFacts,
		"plan":        fd.planFacts,
		"proc":        fd.procFacts,
		"procType":    fd.procTypeFacts,
		"procedure":   fd.procedureFacts,
		"raw":         fd.rawFacts,
		"rawType":     fd.rawTypeFacts,
		"save":        fd.saveFacts,
		"send":        fd.sendFacts,
		"setKind":     fd.setKindFacts,
		"str":         fd.strFacts,
		"strType":     fd.strTypeFacts,
		"transition":  fd.transitionFacts,
		"typeVar":     fd.typeVarFacts,
		"unit":        fd.unitFacts,
	}

	for fileName, lines := range fileToFacts {
		filePath := path.Join(factsInFolder, fmt.Sprintf("%s.facts", fileName))
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
