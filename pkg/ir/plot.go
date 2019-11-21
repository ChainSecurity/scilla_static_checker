package ir

import (
	"errors"
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/multi"
	"strings"
)

type dotNode struct {
	id         int64
	name       string
	ports      []string
	portGroups map[string][]string
}

func (n dotNode) ID() int64 { return n.id }
func (n dotNode) Attributes() []encoding.Attribute {
	groupLabels := make([]string, (len(n.portGroups) + len(n.ports)))
	i := 0
	for groupName, ports := range n.portGroups {
		var portLabels = make([]string, len(ports))
		for j, p := range ports {
			portLabels[j] = fmt.Sprintf("<%s> %s", p, p)
		}
		groupLabel := fmt.Sprintf("{%s | { %s }}", groupName, strings.Join(portLabels, " | "))
		groupLabels[i] = groupLabel
		i = i + 1
	}

	for _, p := range n.ports {
		groupLabels[i] = fmt.Sprintf("<%s> %s", p, p)
		i = i + 1
	}
	label := fmt.Sprintf("{%s|{%s}}", n.name, strings.Join(groupLabels, " | "))
	attrs := []encoding.Attribute{
		{Key: "shape", Value: "record"},
		{Key: "label", Value: label},
	}
	return attrs
}

type dotPortedEdge struct {
	id          int64
	from, to    graph.Node
	fromPort    string
	fromCompass string
	toPort      string
	toCompass   string
}

func (e dotPortedEdge) From() graph.Node { return e.from }
func (e dotPortedEdge) To() graph.Node   { return e.to }
func (e dotPortedEdge) ID() int64        { return e.id }
func (e dotPortedEdge) ReversedEdge() graph.Edge {
	e.from, e.to = e.to, e.from
	e.fromPort, e.toPort = e.toPort, e.fromPort
	e.fromCompass, e.toCompass = e.toCompass, e.fromCompass
	return e
}
func (e dotPortedEdge) ReversedLine() graph.Line {
	e.from, e.to = e.to, e.from
	e.fromPort, e.toPort = e.toPort, e.fromPort
	e.fromCompass, e.toCompass = e.toCompass, e.fromCompass
	return e
}

func (e dotPortedEdge) Weight() float64 { return 0 }

func (e dotPortedEdge) FromPort() (port, compass string) {
	return e.fromPort, e.fromCompass
}
func (e dotPortedEdge) ToPort() (port, compass string) {
	return e.toPort, e.toCompass
}

type dotBuilder struct {
	nodeCounter int64
	edgeCounter int64
	nodes       []graph.Node
	edges       []*dotPortedEdge
	typeCache   map[Type]graph.Node
	dataCache   map[Data]graph.Node
	kindCache   map[Kind]graph.Node
	unitCache   map[Unit]graph.Node
}

func dotWalkJump(b *dotBuilder, jump Jump) graph.Node {
	switch j := jump.(type) {
	case *CallProc:
		m := dotNode{
			b.getNodeId(),
			"CallProc",
			[]string{"To"},
			map[string][]string{},
		}

		for i, _ := range j.Args {
			v := j.Args[i]
			portName := fmt.Sprintf("%s_%d", "Arg", i)
			m.portGroups["Args"] = append(m.portGroups["Args"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, j.To),
			fromPort: "To",
		}

		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, &m)

		return &m
	case *PickProc:
		m := dotNode{
			b.getNodeId(),
			"PickProc",
			[]string{"From"},
			map[string][]string{},
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, j.From),
			fromPort: "From",
		}

		for i, _ := range j.With {
			v := j.With[i]
			portName := fmt.Sprintf("%s_%d", "With", i)
			m.portGroups["With"] = append(m.portGroups["With"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkProcCase(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, &m)
		return &m
	default:
		panic(errors.New(fmt.Sprintf("Wrong Jump type: %T", j)))
	}
}

func dotWalkProcCase(b *dotBuilder, pc ProcCase) graph.Node {
	m := dotNode{
		b.getNodeId(),
		"ProcCase",
		[]string{"Bind", "Body"},
		map[string][]string{},
	}

	e := dotPortedEdge{
		id:       b.getEdgeId(),
		from:     m,
		to:       dotWalkData(b, &pc.Bind),
		fromPort: "Bind",
	}

	b.edges = append(b.edges, &e)

	f := dotPortedEdge{
		id:       b.getEdgeId(),
		from:     m,
		to:       dotWalkData(b, &pc.Body),
		fromPort: "Body",
	}

	b.edges = append(b.edges, &f)

	b.nodes = append(b.nodes, &m)
	return &m
}

func dotWalkUnit(b *dotBuilder, u Unit) graph.Node {
	n, ok := b.unitCache[u]
	if ok {
		return n
	}
	switch x := u.(type) {
	//case *EnumType:
	case *Load:
		slot_field := fmt.Sprintf("Slot: %s", x.Slot)
		m := dotNode{
			b.getNodeId(),
			"Load",
			[]string{slot_field},
			map[string][]string{},
		}
		for i, _ := range x.Path {
			p := x.Path[i]
			portName := fmt.Sprintf("%s_%d", "Path", i)
			m.portGroups["Path"] = append(m.portGroups["Path"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, p),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}
		n = &m
		b.nodes = append(b.nodes, n)
	case *Save:
		slot_field := fmt.Sprintf("Slot: %s", x.Slot)
		m := dotNode{
			b.getNodeId(),
			"Save",
			[]string{slot_field, "Data"},
			map[string][]string{},
		}
		for i, _ := range x.Path {
			p := x.Path[i]
			portName := fmt.Sprintf("%s_%d", "Path", i)
			m.portGroups["Path"] = append(m.portGroups["Path"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, p),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}
		e := &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.Data),
			fromPort: "Data"}
		n = &m
		b.edges = append(b.edges, e)
		b.nodes = append(b.nodes, n)
	case *Event:
		m := dotNode{
			b.getNodeId(),
			"Event",
			[]string{"Data"},
			map[string][]string{},
		}
		e := &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.Data),
			fromPort: "Data"}
		n = &m
		b.edges = append(b.edges, e)
		b.nodes = append(b.nodes, n)
	case *Send:
		m := dotNode{
			b.getNodeId(),
			"Send",
			[]string{"Data"},
			map[string][]string{},
		}
		e := &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.Data),
			fromPort: "Data"}
		n = &m
		b.edges = append(b.edges, e)
		b.nodes = append(b.nodes, n)
	case *Accept:
		n = &dotNode{
			b.getNodeId(),
			"Accept",
			[]string{},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
	case *AppDD:
		n = dotWalkData(b, x)
	case *AbsTD:
		n = dotWalkData(b, x)
	case *AbsDD:
		n = dotWalkData(b, x)
	case *AppTD:
		n = dotWalkData(b, x)
	case *Int:
		n = dotWalkData(b, x)
	case *Nat:
		n = dotWalkData(b, x)
	case *Raw:
		n = dotWalkData(b, x)
	case *Str:
		n = dotWalkData(b, x)
	case *Bnr:
		n = dotWalkData(b, x)
	case *Exc:
		n = dotWalkData(b, x)
	case *Msg:
		n = dotWalkData(b, x)
	//case *Map:
	default:
		panic(errors.New(fmt.Sprintf("unhandeled Unit type: %T", x)))
	}
	b.unitCache[u] = n
	return n
	//return nil
}

func dotWalkType(b *dotBuilder, t Type) graph.Node {

	n, ok := b.typeCache[t]
	if ok {
		return n
	}
	switch x := t.(type) {
	case *EnumType:
		keys := make([]string, 0, len(*x))
		for k := range *x {
			keys = append(keys, k)
		}
		n = &dotNode{
			b.getNodeId(),
			"EnumType",
			keys,
			map[string][]string{},
		}
		b.typeCache[t] = n
		b.nodes = append(b.nodes, n)
		for _, k := range keys {
			for _, inner_t := range (*x)[k] {
				v := dotWalkType(b, inner_t)
				e := &dotPortedEdge{
					id:       b.getEdgeId(),
					from:     n,
					to:       v,
					fromPort: k}
				b.edges = append(b.edges, e)
			}
		}
	case *IntType:
		n = &dotNode{
			b.getNodeId(),
			"IntType",
			[]string{fmt.Sprintf("Size: %d", x.Size)},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
	case *RawType:
		n = &dotNode{
			b.getNodeId(),
			"RawType",
			[]string{fmt.Sprintf("Size: %d", x.Size)},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
	case *NatType:
		n = &dotNode{
			b.getNodeId(),
			"NatType",
			[]string{fmt.Sprintf("Size: %d", x.Size)},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
	case *MapType:
		n = &dotNode{
			b.getNodeId(),
			"MapType",
			[]string{"KeyType", "ValType"},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		kNode := dotWalkType(b, x.KeyType)
		ke := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       kNode,
			fromPort: "KeyType"}
		b.edges = append(b.edges, &ke)
		vNode := dotWalkType(b, x.ValType)
		ve := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       vNode,
			fromPort: "ValType"}
		b.edges = append(b.edges, &ve)
	case *AbsTT:

		m := dotNode{
			b.getNodeId(),
			"AbsTT",
			[]string{"Term"},
			map[string][]string{},
		}

		for i, _ := range x.Vars {
			v := x.Vars[i]
			portName := fmt.Sprintf("%s_%d", "Var", i)
			m.portGroups["Vars"] = append(m.portGroups["Vars"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkType(b, &v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkType(b, x.Term),
			fromPort: "Term"}
		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, m)
		n = &m
	case *AppTT:

		m := dotNode{
			b.getNodeId(),
			fmt.Sprintf("AppTT %p", x),
			[]string{"To"},
			map[string][]string{},
		}

		for i, _ := range x.Args {
			v := x.Args[i]
			portName := fmt.Sprintf("%s_%d", "Arg", i)
			m.portGroups["Arg"] = append(m.portGroups["Arg"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkType(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkType(b, x.To),
			fromPort: "To"}
		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, m)
		n = &m
	case *TypeVar:
		n = &dotNode{
			b.getNodeId(),
			"TypeVar",
			[]string{"Kind"},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		tNode := dotWalkKind(b, x.Kind)
		e := &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       tNode,
			fromPort: "Kind"}
		b.edges = append(b.edges, e)
	case *StrType:
		n = &dotNode{
			b.getNodeId(),
			"StrType",
			[]string{},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
	case *MsgType:
		n = &dotNode{
			b.getNodeId(),
			"MsgType",
			[]string{},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
	default:
		n = &dotNode{
			b.getNodeId(),
			"NIL",
			[]string{},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		fmt.Printf("UNHANDLED TYPE %T\n", t)
		//return nil
		//panic(errors.New(fmt.Sprintf("unhandeled type: %T", x)))
	}

	b.typeCache[t] = n
	return n
}

func dotWalkKind(b *dotBuilder, k Kind) graph.Node {
	n, ok := b.kindCache[k]
	if ok {
		return n
	}
	n = &dotNode{
		b.getNodeId(),
		"Kind",
		[]string{},
		map[string][]string{},
	}
	b.nodes = append(b.nodes, n)
	b.kindCache[k] = n
	return n
}

func dotWalkCond(b *dotBuilder, w *Cond) graph.Node {
	n := dotNode{
		b.getNodeId(),
		"Cond",
		[]string{"Data", fmt.Sprintf("Case: %s", w.Case)},
		map[string][]string{},
	}
	b.nodes = append(b.nodes, &n)
	for i, _ := range w.Data {

		vNode := dotWalkData(b, &w.Data[i])
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       vNode,
			fromPort: "Data"}
		b.edges = append(b.edges, &e)

	}
	return &n
}

func dotWalkDataCase(b *dotBuilder, d *DataCase) graph.Node {
	n := dotNode{
		b.getNodeId(),
		"DataCase",
		[]string{"Bind", "Body"},
		map[string][]string{},
	}
	var e *dotPortedEdge
	e = &dotPortedEdge{
		id:       b.getEdgeId(),
		from:     n,
		to:       dotWalkData(b, d.Body),
		fromPort: "Body"}
	b.edges = append(b.edges, e)
	e = &dotPortedEdge{
		id:       b.getEdgeId(),
		from:     n,
		to:       dotWalkData(b, &d.Bind),
		fromPort: "Bind"}
	b.edges = append(b.edges, e)
	return &n
}

func dotWalkData(b *dotBuilder, d Data) graph.Node {
	n, ok := b.dataCache[d]
	if ok {
		return n
	}
	switch x := d.(type) {
	case *DataVar:
		n = &dotNode{
			b.getNodeId(),
			"DataVar",
			[]string{"DataType"},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		tNode := dotWalkType(b, x.DataType)
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       tNode,
			fromPort: "DataType"}
		b.edges = append(b.edges, &e)
	case *PickData:
		n = &dotNode{
			b.getNodeId(),
			"PickData",
			[]string{"From", "With"},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		for i, _ := range x.With {
			dc := x.With[i]
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     n,
				to:       dotWalkDataCase(b, &dc),
				fromPort: "With"}
			b.edges = append(b.edges, &e)

		}
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       dotWalkData(b, x.From),
			fromPort: "From"}
		b.edges = append(b.edges, &e)
	case *Int:
		n = &dotNode{
			b.getNodeId(),
			"Int",
			[]string{"IntType", fmt.Sprintf("Data: %s", x.Data)},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		tNode := dotWalkType(b, x.IntType)
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       tNode,
			fromPort: "IntType"}
		b.edges = append(b.edges, &e)
	case *Nat:
		n = &dotNode{
			b.getNodeId(),
			"Nat",
			[]string{"NatType", fmt.Sprintf("Data: %s", x.Data)},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		tNode := dotWalkType(b, x.NatType)
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       tNode,
			fromPort: "NatType"}
		b.edges = append(b.edges, &e)
	case *Map:
		n = &dotNode{
			b.getNodeId(),
			"Map",
			[]string{"MapType", fmt.Sprintf("Data: %s", x.Data)},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		tNode := dotWalkType(b, x.MapType)
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       tNode,
			fromPort: "MapType"}
		b.edges = append(b.edges, &e)
	case *AbsTD:

		m := dotNode{
			b.getNodeId(),
			"AbsTD",
			[]string{"Term"},
			map[string][]string{},
		}

		for i, _ := range x.Vars {
			v := x.Vars[i]
			portName := fmt.Sprintf("%s_%d", "Var", i)
			m.portGroups["Vars"] = append(m.portGroups["Vars"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkType(b, &v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.Term),
			fromPort: "Term"}
		b.edges = append(b.edges, &e)

		b.nodes = append(b.nodes, m)
		n = &m

	case *AbsDD:
		m := dotNode{
			b.getNodeId(),
			"AbsDD",
			[]string{"Term"},
			map[string][]string{},
		}

		for i, _ := range x.Vars {
			v := &x.Vars[i]
			portName := fmt.Sprintf("%s_%d", "Var", i)
			m.portGroups["Var"] = append(m.portGroups["Var"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.Term),
			fromPort: "Term"}
		b.edges = append(b.edges, &e)
		n = &m
		b.nodes = append(b.nodes, n)
	case *AppTD:
		m := dotNode{
			b.getNodeId(),
			"AppTD",
			[]string{"To"},
			map[string][]string{},
		}

		for i, _ := range x.Args {
			v := x.Args[i]
			portName := fmt.Sprintf("%s_%d", "Arg", i)
			m.portGroups["Arg"] = append(m.portGroups["Arg"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkType(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.To),
			fromPort: "To"}
		b.edges = append(b.edges, &e)

		b.nodes = append(b.nodes, m)
		n = &m
	case *AppDD:
		m := dotNode{
			b.getNodeId(),
			"AppDD",
			[]string{"To"},
			map[string][]string{},
		}

		for i, _ := range x.Args {
			v := x.Args[i]
			portName := fmt.Sprintf("%s_%d", "Arg", i)
			m.portGroups["Args"] = append(m.portGroups["Args"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkData(b, x.To),
			fromPort: "To"}
		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, m)
		n = &m
	case *Builtin:
		n = &dotNode{
			b.getNodeId(),
			"Builtin",
			[]string{"BuiltinType"},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		tNode := dotWalkType(b, x.BuiltinType)
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       tNode,
			fromPort: "BuiltinType"}
		b.edges = append(b.edges, &e)

	case *Proc:

		m := dotNode{
			b.getNodeId(),
			"Proc",
			[]string{"Jump"},
			map[string][]string{},
		}

		for i, _ := range x.Vars {
			v := &x.Vars[i]
			portName := fmt.Sprintf("%s_%d", "Var", i)
			vNode := dotWalkData(b, v)
			m.portGroups["Var"] = append(m.portGroups["Var"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       vNode,
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		for i, _ := range x.Plan {
			p := x.Plan[i]
			portName := fmt.Sprintf("%s_%d", "Plan", i)
			uNode := dotWalkUnit(b, p)
			m.portGroups["Plan"] = append(m.portGroups["Plan"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       uNode,
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		if x.Jump != nil {
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkJump(b, x.Jump),
				fromPort: "Jump"}

			b.edges = append(b.edges, &e)
		}

		n = &m
		b.nodes = append(b.nodes, n)
	case *Load:
		n = dotWalkUnit(b, x)
	case *Msg:
		m := dotNode{
			b.getNodeId(),
			"Msg",
			[]string{"MsgType"},
			map[string][]string{},
		}

		i := 0
		for k, d := range x.Data {
			portName := fmt.Sprintf("%s_%s", "Data", k)
			m.portGroups["Data"] = append(m.portGroups["Data"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, d),
				fromPort: portName}

			b.edges = append(b.edges, &e)
			i++
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkType(b, x.MsgType),
			fromPort: "MsgType"}

		fmt.Println("Msg", e)

		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, m)
		n = &m
	case *Enum:
		m := dotNode{
			b.getNodeId(),
			"Enum",
			[]string{"EnumType", fmt.Sprintf("Case: %s", x.Case)},
			map[string][]string{},
		}

		for i, _ := range x.Data {
			v := x.Data[i]
			portName := fmt.Sprintf("%s_%d", "Data", i)
			m.portGroups["Data"] = append(m.portGroups["Data"], portName)
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     m,
				to:       dotWalkData(b, v),
				fromPort: portName}
			b.edges = append(b.edges, &e)
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkType(b, x.EnumType),
			fromPort: "EnumType"}

		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, m)
		n = &m
	case *Str:
		m := dotNode{
			b.getNodeId(),
			"Str",
			[]string{"StrType", fmt.Sprintf("Data: %s", x.Data)},
			map[string][]string{},
		}

		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     m,
			to:       dotWalkType(b, x.StrType),
			fromPort: "StrType"}

		b.edges = append(b.edges, &e)
		b.nodes = append(b.nodes, m)
		n = &m

	case *Bind:

		n = dotNode{
			b.getNodeId(),
			"Bind",
			[]string{"BindType", "Cond"},
			map[string][]string{},
		}
		m := dotWalkType(b, x.BindType)
		e := &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       m,
			fromPort: "BindType"}
		b.edges = append(b.edges, e)
		if x.Cond != nil {
			m = dotWalkCond(b, x.Cond)
			e = &dotPortedEdge{
				id:       b.getEdgeId(),
				from:     n,
				to:       m,
				fromPort: "Cond"}
			b.edges = append(b.edges, e)
		}

	default:
		panic(errors.New(fmt.Sprintf("unhandeled type: %T", x)))
	}

	if n == nil {
		panic(errors.New(fmt.Sprintf("dotWalkData returned nil: %T", d)))
	}

	b.dataCache[d] = n
	return n
}

func (b *dotBuilder) getNodeId() int64 {
	id := b.nodeCounter
	b.nodeCounter++
	return id
}

func (b *dotBuilder) getEdgeId() int64 {
	id := b.edgeCounter
	b.edgeCounter++
	return id
}

func directedPortedAttrGraphFrom(b *dotBuilder) graph.Multigraph {
	dg := multi.NewDirectedGraph()
	for _, e := range b.edges {
		dg.SetLine(e)
	}
	return dg
}

func GetDot(b *CFGBuilder) string {
	//keys := make([]string, 0, len(b.GlobalVarMap))
	//for key := range b.GlobalVarMap {
	//keys = append(keys, key)
	//}
	d := dotBuilder{0, 0, []graph.Node{}, []*dotPortedEdge{}, map[Type]graph.Node{}, map[Data]graph.Node{}, map[Kind]graph.Node{}, map[Unit]graph.Node{}}
	//if !ok {
	//panic(errors.New("var not found"))
	//}
	v := b.Transitions["get"]
	//v := b.Transitions["inc"]

	dotWalkData(&d, v)
	g := directedPortedAttrGraphFrom(&d)
	got, err := dot.MarshalMulti(g, "asd", "", "\t")
	_ = got
	if err != nil {
		panic(err)
	}
	return string(got)

}
