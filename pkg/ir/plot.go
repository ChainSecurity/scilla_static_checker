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
}

func dotWalkUnit(b *dotBuilder, u Unit) graph.Node {
	switch x := u.(type) {
	//case *EnumType:
	//case *Load:
	case *Save:
		slot_field := fmt.Sprintf("Slot: %s", x.Slot)
		n := &dotNode{
			b.getNodeId(),
			"Save",
			[]string{slot_field, "Path", "Data"},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)
		for _, a := range x.Path {
			e := dotPortedEdge{
				id:       b.getEdgeId(),
				from:     n,
				to:       dotWalkData(b, a),
				fromPort: "Path"}
			b.edges = append(b.edges, &e)

		}
		e := &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       dotWalkData(b, x.Data),
			fromPort: "Data"}
		b.edges = append(b.edges, e)
		return n
		//Path []Data
		//Data Data
	//case *Emit:
	//case *Send:
	//case *Have:
	//case *AbsTD:
	//case *AbsDD:
	//case *AppDD:
	//case *AppTD:
	//case *Int:
	//case *Nat:
	//case *Raw:
	//case *Str:
	//case *Bnr:
	//case *Exc:
	//case *Msg:
	//case *Map:
	default:
		panic(errors.New(fmt.Sprintf("unhandeled type: %T", x)))
	}
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
		n := &dotNode{
			b.getNodeId(),
			"EnumType",
			keys,
			map[string][]string{},
		}
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
			"AppTT",
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
	default:
		//return nil
		panic(errors.New(fmt.Sprintf("unhandeled type: %T", x)))
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
	for _, v := range w.Data {
		vNode := dotWalkBind(b, &v)
		e := dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       vNode,
			fromPort: "Data"}
		b.edges = append(b.edges, &e)

	}
	return &n
}

func dotWalkBind(b *dotBuilder, d *Bind) graph.Node {
	n := dotNode{
		b.getNodeId(),
		"Bind",
		[]string{"BindType", "Cond"},
		map[string][]string{},
	}
	var m graph.Node
	m = dotWalkType(b, d.BindType)
	var e *dotPortedEdge
	e = &dotPortedEdge{
		id:       b.getEdgeId(),
		from:     n,
		to:       m,
		fromPort: "BindType"}
	b.edges = append(b.edges, e)
	if d.Cond != nil {
		m = dotWalkCond(b, d.Cond)
		e = &dotPortedEdge{
			id:       b.getEdgeId(),
			from:     n,
			to:       m,
			fromPort: "Cond"}
		b.edges = append(b.edges, e)
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
		to:       dotWalkBind(b, &d.Bind),
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
		for _, dc := range x.With {
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
		n = &m
		b.nodes = append(b.nodes, n)
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
			[]string{},
			map[string][]string{},
		}
		b.nodes = append(b.nodes, n)

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
		n = &m
	default:
		panic(errors.New(fmt.Sprintf("unhandeled type: %T", x)))
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
	d := dotBuilder{0, 0, []graph.Node{}, []*dotPortedEdge{}, map[Type]graph.Node{}, map[Data]graph.Node{}, map[Kind]graph.Node{}}
	//v, ok := stackMapPeek(b.fieldStack, "a")
	//if !ok {
	//panic(errors.New("var not found"))
	//}
	v := b.constructor
	dotWalkData(&d, v)
	g := directedPortedAttrGraphFrom(&d)
	got, err := dot.MarshalMulti(g, "asd", "", "\t")
	_ = got
	if err != nil {
		panic(err)
	}
	return string(got)

}
